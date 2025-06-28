package fmp

import "time"

// HistoricalPriceFMP represents a single historical price entry from FMP API
type HistoricalPriceFMP struct {
	Date      string  `json:"date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	VWAP      float64 `json:"vwap"`
	Change    float64 `json:"change"`
	PctChange float64 `json:"changePercent"`
}

// CompanyProfileFMP represents company profile data from FMP API
type CompanyProfileFMP struct {
	Symbol        string `json:"symbol"`
	Price         float64 `json:"price"`
	Beta          float64 `json:"beta"`
	VolAvg        int64 `json:"volAvg"`
	MktCap        int64 `json:"mktCap"`
	LastDiv       float64 `json:"lastDiv"`
	Range         string `json:"range"`
	Changes       float64 `json:"changes"`
	CompanyName   string `json:"companyName"`
	Exchange      string `json:"exchange"`
	Industry      string `json:"industry"`
	Website       string `json:"website"`
	Description   string `json:"description"`
	CEO           string `json:"ceo"`
	Sector        string `json:"sector"`
	Country       string `json:"country"`
	FullTimeEmployees int `json:"fullTimeEmployees"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
}

// FinancialStatementFMP represents a single financial statement entry from FMP API
type FinancialStatementFMP struct {
	Date                 string  `json:"date"`
	Symbol               string  `json:"symbol"`
	ReportedCurrency     string  `json:"reportedCurrency"`
	Cik                  string  `json:"cik"`
	FillingDate          string  `json:"fillingDate"`
	AcceptedDate         string  `json:"acceptedDate"`
	CalendarYear         string  `json:"calendarYear"`
	Period               string  `json:"period"`
	Revenue              float64 `json:"revenue"`
	CostOfRevenue        float64 `json:"costOfRevenue"`
	GrossProfit          float64 `json:"grossProfit"`
	OperatingExpenses    float64 `json:"operatingExpenses"`
	EBITDA               float64 `json:"ebitda"`
	NetIncome            float64 `json:"netIncome"`
	EPS                  float64 `json:"eps"`
	TotalAssets          float64 `json:"totalAssets"`
	TotalLiabilities     float64 `json:"totalLiabilities"`
	TotalEquity          float64 `json:"totalEquity"`
	FreeCashFlow         float64 `json:"freeCashFlow"`
	Debt                 float64 `json:"debt"`
	DebtToEquityRatio    float64 `json:"debtToEquityRatio"`
	// Add other relevant fields as needed based on FMP documentation
}

// AnalystEstimateFMP represents a single analyst estimate entry from FMP API
type AnalystEstimateFMP struct {
	Symbol           string  `json:"symbol"`
	Date             string  `json:"date"`
	EstimatedRevenue float64 `json:"estimatedRevenue"`
	EstimatedEPS     float64 `json:"estimatedEps"`
	// Add other relevant fields as needed
}

// PriceTargetFMP represents a single price target entry from FMP API
type PriceTargetFMP struct {
	Symbol               string  `json:"symbol"`
	PublishedDate        string  `json:"publishedDate"`
	AnalystCompany       string  `json:"analystCompany"`
	PriceTarget          float64 `json:"priceTarget"`
	Recommendation       string  `json:"recommendation"`
	RecommendationStrongBuy int `json:"recommendationStrongBuy"`
	RecommendationBuy    int `json:"recommendationBuy"`
	RecommendationHold   int `json:"recommendationHold"`
	RecommendationSell   int `json:"recommendationSell"`
	RecommendationStrongSell int `json:"recommendationStrongSell"`
}

// SocialSentimentFMP represents a single social sentiment entry from FMP API
type SocialSentimentFMP struct {
	Symbol        string    `json:"symbol"`
	Date          time.Time `json:"date"`
	AbsoluteIndex float64   `json:"absoluteIndex"`
	RelativeIndex float64   `json:"relativeIndex"`
	Sentiment     float64   `json:"sentiment"` // This is the "sentiment field" (overall percentage of positive activity)
	GeneralPerception string `json:"generalPerception"`
	Source        string    `json:"source"`
}
