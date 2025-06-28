package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver

	"stockpick-backend/pkg/database"
	"stockpick-backend/pkg/fmp"
	"stockpick-backend/pkg/models"
	"stockpick-backend/pkg/undervaluation"
)

type App struct {
	Router *mux.Router
	DB     *database.DB
	FMP    *fmp.Client
}

func (a *App) Initialize(dbHost, dbPort, dbUser, dbPassword, dbName, fmpAPIKey string) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	a.DB = db

	a.FMP = fmp.NewClient(fmpAPIKey)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Printf("Server starting on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/health", a.healthCheckHandler).Methods("GET")
	a.Router.HandleFunc("/api/ingest/historical-prices/{symbol}", a.ingestHistoricalPricesHandler).Methods("POST")
	a.Router.HandleFunc("/api/stocks/{symbol}/history", a.getHistoricalPricesHandler).Methods("GET")
	a.Router.HandleFunc("/api/stocks", a.getStocksHandler).Methods("GET")
	a.Router.HandleFunc("/api/undervalued", a.getUndervaluedStocksHandler).Methods("GET")
}

func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is healthy!")
}

func (a *App) ingestHistoricalPricesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	log.Printf("Ingesting historical prices for %s", symbol)

	// For simplicity, fetch data for the last year. In a real app, this would be more dynamic.
	to := time.Now()
	from := to.AddDate(-1, 0, 0) // Last 1 year

	fmpPrices, err := a.FMP.GetHistoricalPrices(symbol, from, to)
	if err != nil {
		log.Printf("Error fetching historical prices from FMP for %s: %v", symbol, err)
		http.Error(w, "Failed to fetch historical prices", http.StatusInternalServerError)
		return
	}

	// First, get or create the stock in our DB
	stock, err := a.DB.GetStockBySymbol(symbol)
	if err != nil {
		log.Printf("Error getting stock by symbol %s: %v", symbol, err)
		http.Error(w, "Failed to process stock", http.StatusInternalServerError)
		return
	}

	if stock == nil {
		// Attempt to get company profile to populate stock details
		profiles, err := a.FMP.GetCompanyProfile(symbol)
		if err != nil || len(profiles) == 0 {
			log.Printf("Could not get company profile for new stock %s: %v", symbol, err)
			http.Error(w, "Failed to get company profile for new stock", http.StatusInternalServerError)
			return
		}
		profile := profiles[0]

		stock = &models.Stock{
			Symbol:      symbol,
			CompanyName: profile.CompanyName,
			Exchange:    profile.Exchange,
			Sector:      profile.Sector,
			Industry:    profile.Industry,
			Currency:    "USD", // FMP usually provides USD for US stocks
			IsActive:    true,
		}
		err = a.DB.InsertStock(stock)
		if err != nil {
			log.Printf("Error inserting new stock %s: %v", symbol, err)
			http.Error(w, "Failed to insert new stock", http.StatusInternalServerError)
			return
		}
		log.Printf("Inserted new stock: %s (%s)", stock.CompanyName, stock.Symbol)
	}

	// Insert historical prices into DB
	for _, p := range fmpPrices {
		priceTime, err := time.Parse("2006-01-02", p.Date)
		if err != nil {
			log.Printf("Error parsing date %s: %v", p.Date, err)
			continue
		}
		hp := &models.HistoricalPrice{
			Time:       priceTime,
			StockID:    stock.StockID,
			OpenPrice:  p.Open,
			HighPrice:  p.High,
			LowPrice:   p.Low,
			ClosePrice: p.Close,
			Volume:     p.Volume,
			VWAP:       p.VWAP,
			PriceChange: p.Change,
			PctChange:  p.PctChange,
		}
		err = a.DB.InsertHistoricalPrice(hp)
		if err != nil {
			log.Printf("Error inserting historical price for %s on %s: %v", symbol, p.Date, err)
			// Continue to next price, don't stop the whole ingestion
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully ingested historical prices for %s", symbol)
}

func (a *App) getHistoricalPricesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	stock, err := a.DB.GetStockBySymbol(symbol)
	if err != nil {
		log.Printf("Error getting stock by symbol %s: %v", symbol, err)
		http.Error(w, "Failed to retrieve stock data", http.StatusInternalServerError)
		return
	}

	if stock == nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	// For simplicity, retrieve data for the last year
	to := time.Now()
	from := to.AddDate(-1, 0, 0)

	prices, err := a.DB.GetHistoricalPrices(stock.StockID, from, to)
	if err != nil {
		log.Printf("Error retrieving historical prices for %s: %v", symbol, err)
		http.Error(w, "Failed to retrieve historical prices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

func (a *App) getStocksHandler(w http.ResponseWriter, r *http.Request) {
	stocks, err := a.DB.GetAllStocks()
	if err != nil {
		log.Printf("Error retrieving all stocks: %v", err)
		http.Error(w, "Failed to retrieve stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func (a *App) getUndervaluedStocksHandler(w http.ResponseWriter, r *http.Request) {
	allStocks, err := a.DB.GetAllStocks()
	if err != nil {
		log.Printf("Error retrieving all stocks for undervaluation: %v", err)
		http.Error(w, "Failed to retrieve stocks for undervaluation", http.StatusInternalServerError)
		return
	}

	var undervaluedStocks []undervaluation.UndervaluationScore
	for _, stock := range allStocks {
		// Fetch latest price (from historical prices)
		// For simplicity, get the very last closing price
		prices, err := a.DB.GetHistoricalPrices(stock.StockID, time.Now().AddDate(-1, 0, 0), time.Now())
		if err != nil || len(prices) == 0 {
			log.Printf("Could not get latest price for %s: %v", stock.Symbol, err)
			continue // Skip if no price data
		}
		latestPrice := prices[len(prices)-1].ClosePrice

		// Fetch latest financial statements (annual)
		financialStatements, err := a.DB.GetFinancialStatements(stock.StockID, "annual")
		if err != nil {
			log.Printf("Could not get financial statements for %s: %v", stock.Symbol, err)
			financialStatements = []models.FinancialStatement{} // Ensure it's not nil for calculator
		}

		// Fetch latest analyst targets
		analystTargets, err := a.DB.GetAnalystTargets(stock.StockID)
		if err != nil {
			log.Printf("Could not get analyst targets for %s: %v", stock.Symbol, err)
			analystTargets = []models.AnalystTarget{} // Ensure it's not nil
		}

		// Fetch latest sentiment scores (e.g., overall source)
		sentimentScores, err := a.DB.GetSentimentScores(stock.StockID, time.Now().AddDate(0, 0, -7), time.Now(), "Overall")
		if err != nil {
			log.Printf("Could not get sentiment scores for %s: %v", stock.Symbol, err)
			sentimentScores = []models.SentimentScore{} // Ensure it's not nil
		}

		score, err := undervaluation.CalculateUndervaluation(&stock, latestPrice, financialStatements, analystTargets, sentimentScores)
		if err != nil {
			log.Printf("Error calculating undervaluation for %s: %v", stock.Symbol, err)
			continue
		}

		// Define a threshold for "undervalued"
		if score.Score >= 50 { // Example threshold
			undervaluedStocks = append(undervaluedStocks, *score)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(undervaluedStocks)
}


func main() {
	app := App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("FMP_API_KEY"),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	app.Run(":" + port)
}
