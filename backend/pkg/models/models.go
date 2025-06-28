package models

import (
	"time"

	"github.com/google/uuid"
)

// Stock represents a stock in the database
type Stock struct {
	StockID     uuid.UUID `json:"stock_id" db:"stock_id"`
	Symbol      string    `json:"symbol" db:"symbol"`
	CompanyName string    `json:"company_name" db:"company_name"`
	Exchange    string    `json:"exchange" db:"exchange"`
	Sector      string    `json:"sector" db:"sector"`
	Industry    string    `json:"industry" db:"industry"`
	Currency    string    `json:"currency" db:"currency"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// HistoricalPrice represents historical OHLCV data for a stock
type HistoricalPrice struct {
	Time       time.Time `json:"time" db:"time"`
	StockID    uuid.UUID `json:"stock_id" db:"stock_id"`
	OpenPrice  float64   `json:"open_price" db:"open_price"`
	HighPrice  float64   `json:"high_price" db:"high_price"`
	LowPrice   float64   `json:"low_price" db:"low_price"`
	ClosePrice float64   `json:"close_price" db:"close_price"`
	Volume     int64     `json:"volume" db:"volume"`
	VWAP       float64   `json:"vwap" db:"vwap"`
	PriceChange float64   `json:"price_change" db:"price_change"`
	PctChange  float64   `json:"pct_change" db:"pct_change"`
}

// FinancialStatement represents annual or quarterly financial data
type FinancialStatement struct {
	StatementID      uuid.UUID `json:"statement_id" db:"statement_id"`
	StockID          uuid.UUID `json:"stock_id" db:"stock_id"`
	Date             time.Time `json:"date" db:"date"`
	Period           string    `json:"period" db:"period"`
	Revenue          float64   `json:"revenue" db:"revenue"`
	NetIncome        float64   `json:"net_income" db:"net_income"`
	EPS              float64   `json:"eps" db:"eps"`
	TotalAssets      float64   `json:"total_assets" db:"total_assets"`
	TotalLiabilities float64   `json:"total_liabilities" db:"total_liabilities"`
	TotalEquity      float64   `json:"total_equity" db:"total_equity"`
	FreeCashFlow     float64   `json:"free_cash_flow" db:"free_cash_flow"`
	DebtToEquityRatio float64   `json:"debt_to_equity_ratio" db:"debt_to_equity_ratio"`
	PERatio          float64   `json:"p_e_ratio" db:"p_e_ratio"`
	PBRatio          float64   `json:"p_b_ratio" db:"p_b_ratio"`
	ROIC             float64   `json:"roic" db:"roic"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// AnalystTarget represents aggregated analyst consensus data
type AnalystTarget struct {
	TargetID              uuid.UUID `json:"target_id" db:"target_id"`
	StockID               uuid.UUID `json:"stock_id" db:"stock_id"`
	Date                  time.Time `json:"date" db:"date"`
	ConsensusPriceTarget  float64   `json:"consensus_price_target" db:"consensus_price_target"`
	HighPriceTarget       float64   `json:"high_price_target" db:"high_price_target"`
	LowPriceTarget        float64   `json:"low_price_target" db:"low_price_target"`
	ConsensusRating       string    `json:"consensus_rating" db:"consensus_rating"`
	ConsensusRatingValue  float64   `json:"consensus_rating_value" db:"consensus_rating_value"`
	BuyRatingsCount       int       `json:"buy_ratings_count" db:"buy_ratings_count"`
	HoldRatingsCount      int       `json:"hold_ratings_count" db:"hold_ratings_count"`
	SellRatingsCount      int       `json:"sell_ratings_count" db:"sell_ratings_count"`
	TotalAnalystsContributing int       `json:"total_analysts_contributing" db:"total_analysts_contributing"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// SentimentScore represents social media sentiment data
type SentimentScore struct {
	SentimentID      uuid.UUID `json:"sentiment_id" db:"sentiment_id"`
	StockID          uuid.UUID `json:"stock_id" db:"stock_id"`
	Timestamp        time.Time `json:"timestamp" db:"timestamp"`
	AbsoluteIndex    float64   `json:"absolute_index" db:"absolute_index"`
	RelativeIndex    float64   `json:"relative_index" db:"relative_index"`
	SentimentScore   float64   `json:"sentiment_score" db:"sentiment_score"`
	GeneralPerception string    `json:"general_perception" db:"general_perception"`
	Source           string    `json:"source" db:"source"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
