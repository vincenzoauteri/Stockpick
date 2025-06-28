package fmp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ( 
	BaseURL = "https://financialmodelingprep.com/api/v3"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) get(path string, queryParams map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", BaseURL, path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	q.Add("apikey", c.APIKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("FMP API returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}

// GetHistoricalPrices fetches historical OHLCV data for a given symbol.
func (c *Client) GetHistoricalPrices(symbol string, from, to time.Time) ([]HistoricalPriceFMP, error) {
	path := fmt.Sprintf("/historical-price-full/%s", symbol)
	queryParams := map[string]string{
		"from": from.Format("2006-01-02"),
		"to":   to.Format("2006-01-02"),
	}

	body, err := c.get(path, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Symbol string               `json:"symbol"`
		Historical []HistoricalPriceFMP `json:"historical"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal historical prices: %w", err)
	}

	return response.Historical, nil
}

// GetCompanyProfile fetches company profile data for a given symbol.
func (c *Client) GetCompanyProfile(symbol string) ([]CompanyProfileFMP, error) {
	path := fmt.Sprintf("/profile/%s", symbol)
	body, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}

	var profiles []CompanyProfileFMP
	if err := json.Unmarshal(body, &profiles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company profile: %w", err)
	}
	return profiles, nil
}

// GetFinancialStatements fetches financial statements (income, balance, cash flow) for a given symbol and period.
func (c *Client) GetFinancialStatements(symbol, statementType, period string) ([]FinancialStatementFMP, error) {
	path := fmt.Sprintf("/%s-statement/%s", statementType, symbol)
	queryParams := map[string]string{
		"period": period,
	}
	body, err := c.get(path, queryParams)
	if err != nil {
		return nil, err
	}

	var statements []FinancialStatementFMP
	if err := json.Unmarshal(body, &statements); err != nil {
		return nil, fmt.Errorf("failed to unmarshal financial statements: %w", err)
	}
	return statements, nil
}

// GetAnalystEstimates fetches analyst estimates for a given symbol.
func (c *Client) GetAnalystEstimates(symbol string) ([]AnalystEstimateFMP, error) {
	path := fmt.Sprintf("/analyst-estimates/%s", symbol)
	body, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}

	var estimates []AnalystEstimateFMP
	if err := json.Unmarshal(body, &estimates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal analyst estimates: %w", err)
	}
	return estimates, nil
}

// GetPriceTargetConsensus fetches price target consensus for a given symbol.
func (c *Client) GetPriceTargetConsensus(symbol string) ([]PriceTargetFMP, error) {
	path := fmt.Sprintf("/price-target-consensus/%s", symbol)
	body, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}

	var targets []PriceTargetFMP
	if err := json.Unmarshal(body, &targets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal price targets: %w", err)
	}
	return targets, nil
}

// GetSocialSentiment fetches social sentiment data for a given symbol.
func (c *Client) GetSocialSentiment(symbol string, from, to time.Time) ([]SocialSentimentFMP, error) {
	path := fmt.Sprintf("/historical/social-sentiment/%s", symbol)
	queryParams := map[string]string{
		"from": from.Format("2006-01-02"),
		"to":   to.Format("2006-01-02"),
	}
	body, err := c.get(path, queryParams)
	if err != nil {
		return nil, err
	}

	var sentiment []SocialSentimentFMP
	if err := json.Unmarshal(body, &sentiment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal social sentiment: %w", err)
	}
	return sentiment, nil
}
