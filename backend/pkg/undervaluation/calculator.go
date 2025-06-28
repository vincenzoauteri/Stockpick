package undervaluation

import (
	"fmt"
	"math"
	"stockpick-backend/pkg/models"
)

// UndervaluationScore represents the calculated undervaluation score for a stock
type UndervaluationScore struct {
	StockID          string  `json:"stock_id"`
	Symbol           string  `json:"symbol"`
	Score            float64 `json:"score"`
	FundamentalScore float64 `json:"fundamental_score"`
	AnalystScore     float64 `json:"analyst_score"`
	SentimentScore   float64 `json:"sentiment_score"`
	// Add more detailed breakdown if needed
}

// CalculateUndervaluation calculates a composite undervaluation score for a stock.
// This is a simplified model based on the provided specs.
func CalculateUndervaluation(
	stock *models.Stock,
	latestPrice float64,
	financialStatements []models.FinancialStatement,
	analystTargets []models.AnalystTarget,
	sentimentScores []models.SentimentScore,
) (*UndervaluationScore, error) {

	if stock == nil || latestPrice == 0 {
		return nil, fmt.Errorf("invalid input: stock or latest price is missing")
	}

	// --- Fundamental Analysis Score (Simplified) ---
	fundamentalScore := 0.0
	if len(financialStatements) > 0 {
		latestFS := financialStatements[0] // Assuming latest is first due to DESC order in query

		// P/E Ratio (lower is better, relative to some benchmark)
		if latestFS.PERatio > 0 {
			// Example: If P/E is very low (e.g., < 10), give higher score
			if latestFS.PERatio < 10 {
				fundamentalScore += 0.3
			} else if latestFS.PERatio < 20 {
				fundamentalScore += 0.15
			}
		}

		// EPS Growth (higher is better)
		// This requires historical EPS, which we don't have directly in latestFS. For simplicity,
		// we'll assume a positive EPS is good for now. In a real scenario, compare current EPS to previous.
		if latestFS.EPS > 0 {
			fundamentalScore += 0.2
		}

		// ROIC (higher is better, e.g., > 15%)
		if latestFS.ROIC > 0.15 { // 15%
			fundamentalScore += 0.2
		}

		// Free Cash Flow (positive and growing is good)
		if latestFS.FreeCashFlow > 0 {
			fundamentalScore += 0.15
		}
	}

	// --- Analyst Consensus Score (Simplified) ---
	analystScore := 0.0
	if len(analystTargets) > 0 {
		latestAT := analystTargets[0] // Assuming latest is first

		// Price Target Upside
		if latestAT.ConsensusPriceTarget > 0 && latestPrice > 0 {
			upside := (latestAT.ConsensusPriceTarget - latestPrice) / latestPrice
			if upside > 0.20 { // > 20% upside
				analystScore += 0.4
			} else if upside > 0.10 { // > 10% upside
				analystScore += 0.2
			}
		}

		// Consensus Rating Value (1=Strong Sell, 5=Strong Buy)
		if latestAT.ConsensusRatingValue >= 4.0 { // Buy or Strong Buy
			analystScore += 0.3
		} else if latestAT.ConsensusRatingValue >= 3.0 { // Hold
			analystScore += 0.15
		}
	}

	// --- Sentiment Analysis Score (Simplified) ---
	sentimentScore := 0.0
	if len(sentimentScores) > 0 {
		latestSS := sentimentScores[0] // Assuming latest is first

		// Sentiment Score (e.g., 0 to 100, higher is better)
		// Normalize to 0-1 scale if it's 0-100
	normalizedSentiment := latestSS.SentimentScore
	if normalizedSentiment > 1 {
		normalizedSentiment /= 100.0
	}

		if normalizedSentiment > 0.7 { // High positive sentiment
			sentimentScore += 0.2
		} else if normalizedSentiment > 0.5 { // Neutral to slightly positive
			sentimentScore += 0.1
		}

		// Absolute Index (discussion volume) - higher indicates more interest
		if latestSS.AbsoluteIndex > 100000 { // Arbitrary high volume threshold
			sentimentScore += 0.1
		}
	}

	// --- Combine Scores with Weights (Example Weights) ---
	// Weights are illustrative and can be fine-tuned
	const ( 
		FundamentalWeight = 0.50
		AnalystWeight     = 0.30
		SentimentWeight   = 0.20
	)

	compositeScore := (
		fundamentalScore*FundamentalWeight + 
		analystScore*AnalystWeight + 
		sentimentScore*SentimentWeight
	) * 100 // Scale to 0-100 for easier interpretation

	// Cap the score at 100
	compositeScore = math.Min(compositeScore, 100.0)

	return &UndervaluationScore{
		StockID:          stock.StockID.String(),
		Symbol:           stock.Symbol,
		Score:            compositeScore,
		FundamentalScore: fundamentalScore * 100, // Scaled for output
		AnalystScore:     analystScore * 100,
		SentimentScore:   sentimentScore * 100,
	}, nil
}
