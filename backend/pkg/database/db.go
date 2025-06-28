package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"stockpick-backend/pkg/models"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &DB{db},
		nil
}

// InsertStock inserts a new stock into the database
func (d *DB) InsertStock(stock *models.Stock) error {
	query := `INSERT INTO stocks (stock_id, symbol, company_name, exchange, sector, industry, currency, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (symbol) DO UPDATE SET
		company_name = EXCLUDED.company_name, exchange = EXCLUDED.exchange, sector = EXCLUDED.sector,
		industry = EXCLUDED.industry, currency = EXCLUDED.currency, is_active = EXCLUDED.is_active,
		updated_at = NOW() RETURNING stock_id`

	stock.StockID = uuid.New()
	stock.CreatedAt = time.Now()
	stock.UpdatedAt = time.Now()

	err := d.QueryRow(query,
		stock.StockID, stock.Symbol, stock.CompanyName, stock.Exchange, stock.Sector,
		stock.Industry, stock.Currency, stock.IsActive, stock.CreatedAt, stock.UpdatedAt).Scan(&stock.StockID)

	if err != nil {
		return fmt.Errorf("failed to insert stock: %w", err)
	}
	return nil
}

// GetStockBySymbol retrieves a stock by its symbol
func (d *DB) GetStockBySymbol(symbol string) (*models.Stock, error) {
	query := `SELECT stock_id, symbol, company_name, exchange, sector, industry, currency, is_active, created_at, updated_at
		FROM stocks WHERE symbol = $1`
	
	stock := &models.Stock{}
	err := d.QueryRow(query, symbol).Scan(
		&stock.StockID, &stock.Symbol, &stock.CompanyName, &stock.Exchange, &stock.Sector,
		&stock.Industry, &stock.Currency, &stock.IsActive, &stock.CreatedAt, &stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Stock not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get stock by symbol: %w", err)
	}
	return stock, nil
}

// GetAllStocks retrieves all stocks from the database
func (d *DB) GetAllStocks() ([]models.Stock, error) {
	query := `SELECT stock_id, symbol, company_name, exchange, sector, industry, currency, is_active, created_at, updated_at FROM stocks ORDER BY symbol ASC`

	rows, err := d.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all stocks: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(
			&stock.StockID, &stock.Symbol, &stock.CompanyName, &stock.Exchange, &stock.Sector,
			&stock.Industry, &stock.Currency, &stock.IsActive, &stock.CreatedAt, &stock.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning stock row: %v", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock rows: %w", err)
	}

	return stocks, nil
}

// InsertHistoricalPrice inserts a new historical price record
func (d *DB) InsertHistoricalPrice(price *models.HistoricalPrice) error {
	query := `INSERT INTO historical_prices (time, stock_id, open_price, high_price, low_price, close_price, volume, vwap, price_change, pct_change)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (time, stock_id) DO UPDATE SET
		open_price = EXCLUDED.open_price, high_price = EXCLUDED.high_price, low_price = EXCLUDED.low_price,
		close_price = EXCLUDED.close_price, volume = EXCLUDED.volume, vwap = EXCLUDED.vwap,
		price_change = EXCLUDED.price_change, pct_change = EXCLUDED.pct_change`

	_, err := d.Exec(query, price.Time, price.StockID, price.OpenPrice, price.HighPrice, price.LowPrice,
		price.ClosePrice, price.Volume, price.VWAP, price.PriceChange, price.PctChange)
	if err != nil {
		return fmt.Errorf("failed to insert historical price: %w", err)
	}
	return nil
}

// GetHistoricalPrices retrieves historical prices for a stock within a time range
func (d *DB) GetHistoricalPrices(stockID uuid.UUID, from, to time.Time) ([]models.HistoricalPrice, error) {
	query := `SELECT time, stock_id, open_price, high_price, low_price, close_price, volume, vwap, price_change, pct_change
		FROM historical_prices WHERE stock_id = $1 AND time BETWEEN $2 AND $3 ORDER BY time ASC`

	rows, err := d.Query(query, stockID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query historical prices: %w", err)
	}
	defer rows.Close()

	var prices []models.HistoricalPrice
	for rows.Next() {
		var price models.HistoricalPrice
		err := rows.Scan(
			&price.Time, &price.StockID, &price.OpenPrice, &price.HighPrice, &price.LowPrice,
			&price.ClosePrice, &price.Volume, &price.VWAP, &price.PriceChange, &price.PctChange,
		)
		if err != nil {
			log.Printf("Error scanning historical price row: %v", err)
			continue
		}
		prices = append(prices, price)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating historical prices rows: %w", err)
	}

	return prices, nil
}

// InsertFinancialStatement inserts a new financial statement record
func (d *DB) InsertFinancialStatement(statement *models.FinancialStatement) error {
	query := `INSERT INTO financial_statements (statement_id, stock_id, date, period, revenue, net_income, eps, total_assets, total_liabilities, total_equity, free_cash_flow, debt_to_equity_ratio, p_e_ratio, p_b_ratio, roic, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (stock_id, date, period) DO UPDATE SET
		revenue = EXCLUDED.revenue, net_income = EXCLUDED.net_income, eps = EXCLUDED.eps,
		total_assets = EXCLUDED.total_assets, total_liabilities = EXCLUDED.total_liabilities,
		total_equity = EXCLUDED.total_equity, free_cash_flow = EXCLUDED.free_cash_flow,
		debt_to_equity_ratio = EXCLUDED.debt_to_equity_ratio, p_e_ratio = EXCLUDED.p_e_ratio,
		p_b_ratio = EXCLUDED.p_b_ratio, roic = EXCLUDED.roic, updated_at = NOW()`

	statement.StatementID = uuid.New()
	statement.CreatedAt = time.Now()
	statement.UpdatedAt = time.Now()

	_, err := d.Exec(query, statement.StatementID, statement.StockID, statement.Date, statement.Period,
		statement.Revenue, statement.NetIncome, statement.EPS, statement.TotalAssets, statement.TotalLiabilities,
		statement.TotalEquity, statement.FreeCashFlow, statement.DebtToEquityRatio, statement.PERatio,
		statement.PBRatio, statement.ROIC, statement.CreatedAt, statement.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert financial statement: %w", err)
	}
	return nil
}

// GetFinancialStatements retrieves financial statements for a stock by period
func (d *DB) GetFinancialStatements(stockID uuid.UUID, period string) ([]models.FinancialStatement, error) {
	query := `SELECT statement_id, stock_id, date, period, revenue, net_income, eps, total_assets, total_liabilities, total_equity, free_cash_flow, debt_to_equity_ratio, p_e_ratio, p_b_ratio, roic, created_at, updated_at
		FROM financial_statements WHERE stock_id = $1 AND period = $2 ORDER BY date DESC`

	rows, err := d.Query(query, stockID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to query financial statements: %w", err)
	}
	defer rows.Close()

	var statements []models.FinancialStatement
	for rows.Next() {
		var statement models.FinancialStatement
		err := rows.Scan(
			&statement.StatementID, &statement.StockID, &statement.Date, &statement.Period,
			&statement.Revenue, &statement.NetIncome, &statement.EPS, &statement.TotalAssets, &statement.TotalLiabilities,
			&statement.TotalEquity, &statement.FreeCashFlow, &statement.DebtToEquityRatio, &statement.PERatio,
			&statement.PBRatio, &statement.ROIC, &statement.CreatedAt, &statement.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning financial statement row: %v", err)
			continue
		}
		statements = append(statements, statement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating financial statements rows: %w", err)
	}

	return statements, nil
}

// InsertAnalystTarget inserts a new analyst target record
func (d *DB) InsertAnalystTarget(target *models.AnalystTarget) error {
	query := `INSERT INTO analyst_targets (target_id, stock_id, date, consensus_price_target, high_price_target, low_price_target, consensus_rating, consensus_rating_value, buy_ratings_count, hold_ratings_count, sell_ratings_count, total_analysts_contributing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (stock_id, date) DO UPDATE SET
		consensus_price_target = EXCLUDED.consensus_price_target, high_price_target = EXCLUDED.high_price_target,
		low_price_target = EXCLUDED.low_price_target, consensus_rating = EXCLUDED.consensus_rating,
		consensus_rating_value = EXCLUDED.consensus_rating_value, buy_ratings_count = EXCLUDED.buy_ratings_count,
		hold_ratings_count = EXCLUDED.hold_ratings_count, sell_ratings_count = EXCLUDED.sell_ratings_count,
		total_analysts_contributing = EXCLUDED.total_analysts_contributing, updated_at = NOW()`

	target.TargetID = uuid.New()
	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()

	_, err := d.Exec(query, target.TargetID, target.StockID, target.Date, target.ConsensusPriceTarget,
		target.HighPriceTarget, target.LowPriceTarget, target.ConsensusRating, target.ConsensusRatingValue,
		target.BuyRatingsCount, target.HoldRatingsCount, target.SellRatingsCount, target.TotalAnalystsContributing,
		target.CreatedAt, target.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert analyst target: %w", err)
	}
	return nil
}

// GetAnalystTargets retrieves analyst targets for a stock
func (d *DB) GetAnalystTargets(stockID uuid.UUID) ([]models.AnalystTarget, error) {
	query := `SELECT target_id, stock_id, date, consensus_price_target, high_price_target, low_price_target, consensus_rating, consensus_rating_value, buy_ratings_count, hold_ratings_count, sell_ratings_count, total_analysts_contributing, created_at, updated_at
		FROM analyst_targets WHERE stock_id = $1 ORDER BY date DESC`

	rows, err := d.Query(query, stockID)
	if err != nil {
		return nil, fmt.Errorf("failed to query analyst targets: %w", err)
	}
	defer rows.Close()

	var targets []models.AnalystTarget
	for rows.Next() {
		var target models.AnalystTarget
		err := rows.Scan(
			&target.TargetID, &target.StockID, &target.Date, &target.ConsensusPriceTarget,
			&target.HighPriceTarget, &target.LowPriceTarget, &target.ConsensusRating, &target.ConsensusRatingValue,
			&target.BuyRatingsCount, &target.HoldRatingsCount, &target.SellRatingsCount, &target.TotalAnalystsContributing,
			&target.CreatedAt, &target.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning analyst target row: %v", err)
			continue
		}
		targets = append(targets, target)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating analyst targets rows: %w", err)
	}

	return targets, nil
}

// InsertSentimentScore inserts a new sentiment score record
func (d *DB) InsertSentimentScore(sentiment *models.SentimentScore) error {
	query := `INSERT INTO sentiment_scores (sentiment_id, stock_id, timestamp, absolute_index, relative_index, sentiment_score, general_perception, source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (stock_id, timestamp, source) DO UPDATE SET
		absolute_index = EXCLUDED.absolute_index, relative_index = EXCLUDED.relative_index,
		sentiment_score = EXCLUDED.sentiment_score, general_perception = EXCLUDED.general_perception,
		updated_at = NOW()`

	sentiment.SentimentID = uuid.New()
	sentiment.CreatedAt = time.Now()
	sentiment.UpdatedAt = time.Now()

	_, err := d.Exec(query, sentiment.SentimentID, sentiment.StockID, sentiment.Timestamp,
		sentiment.AbsoluteIndex, sentiment.RelativeIndex, sentiment.SentimentScore, sentiment.GeneralPerception,
		sentiment.Source, sentiment.CreatedAt, sentiment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert sentiment score: %w", err)
	}
	return nil
}

// GetSentimentScores retrieves sentiment scores for a stock within a time range and source
func (d *DB) GetSentimentScores(stockID uuid.UUID, from, to time.Time, source string) ([]models.SentimentScore, error) {
	query := `SELECT sentiment_id, stock_id, timestamp, absolute_index, relative_index, sentiment_score, general_perception, source, created_at, updated_at
		FROM sentiment_scores WHERE stock_id = $1 AND timestamp BETWEEN $2 AND $3 AND source = $4 ORDER BY timestamp ASC`

	rows, err := d.Query(query, stockID, from, to, source)
	if err != nil {
		return nil, fmt.Errorf("failed to query sentiment scores: %w", err)
	}
	defer rows.Close()

	var sentiments []models.SentimentScore
	for rows.Next() {
		var sentiment models.SentimentScore
		err := rows.Scan(
			&sentiment.SentimentID, &sentiment.StockID, &sentiment.Timestamp, &sentiment.AbsoluteIndex,
			&sentiment.RelativeIndex, &sentiment.SentimentScore, &sentiment.GeneralPerception, &sentiment.Source,
			&sentiment.CreatedAt, &sentiment.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning sentiment score row: %v", err)
			continue
		}
		sentiments = append(sentiments, sentiment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sentiment scores rows: %w", err)
	}

	return sentiments, nil
}