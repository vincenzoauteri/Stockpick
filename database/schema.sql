-- Create the stocks table
CREATE TABLE stocks (
    stock_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for the stock
    symbol TEXT UNIQUE NOT NULL,                          -- Stock ticker symbol (e.g., AAPL)
    company_name TEXT NOT NULL,                           -- Full company name
    exchange TEXT,                                        -- Exchange where the stock is traded (e.g., NASDAQ)
    sector TEXT,                                          -- Industry sector (e.g., Technology)
    industry TEXT,                                        -- Specific industry (e.g., Consumer Electronics)
    currency TEXT,                                        -- Trading currency (e.g., USD)
    is_active BOOLEAN DEFAULT TRUE,                       -- Indicates if the stock is actively trading
    created_at TIMESTAMPTZ DEFAULT NOW(),                 -- Timestamp of record creation
    updated_at TIMESTAMPTZ DEFAULT NOW()                  -- Timestamp of last record update
);

-- Create the historical_prices table (Hypertable)
CREATE TABLE historical_prices (
    time TIMESTAMPTZ NOT NULL,                            -- Timestamp of the price data (TimescaleDB time dimension)
    stock_id UUID NOT NULL REFERENCES stocks(stock_id),   -- Foreign key to stocks table (TimescaleDB space dimension)
    open_price DOUBLE PRECISION,                          -- Opening price
    high_price DOUBLE PRECISION,                          -- Highest price
    low_price DOUBLE PRECISION,                           -- Lowest price
    close_price DOUBLE PRECISION,                         -- Closing price
    volume BIGINT,                                        -- Trading volume
    vwap DOUBLE PRECISION,                                -- Volume-Weighted Average Price
    price_change DOUBLE PRECISION,                        -- Absolute price change
    pct_change DOUBLE PRECISION,                          -- Percentage price change
    PRIMARY KEY (time, stock_id)
);

-- Convert to TimescaleDB hypertable, partitioned by time and symbol for performance
SELECT create_hypertable('historical_prices', 'time', 'stock_id', number_partitions => 4);

-- Create the financial_statements table
CREATE TABLE financial_statements (
    statement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for the statement record
    stock_id UUID NOT NULL REFERENCES stocks(stock_id),      -- Foreign key to stocks table
    date DATE NOT NULL,                                      -- Date of the financial statement (e.g., end of quarter/year)
    period TEXT NOT NULL,                                    -- 'annual' or 'quarterly'
    revenue NUMERIC(20, 2),                                  -- Total revenue
    net_income NUMERIC(20, 2),                               -- Net income
    eps NUMERIC(10, 4),                                      -- Earnings Per Share
    total_assets NUMERIC(20, 2),                             -- Total assets
    total_liabilities NUMERIC(20, 2),                        -- Total liabilities
    total_equity NUMERIC(20, 2),                             -- Total equity
    free_cash_flow NUMERIC(20, 2),                           -- Free Cash Flow
    debt_to_equity_ratio NUMERIC(10, 4),                     -- Debt-to-Equity Ratio
    p_e_ratio NUMERIC(10, 4),                                -- Price-to-Earnings Ratio
    p_b_ratio NUMERIC(10, 4),                                -- Price-to-Book Ratio
    roic NUMERIC(10, 4),                                     -- Return on Invested Capital
    created_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of record creation
    updated_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of last record update
    UNIQUE (stock_id, date, period)                          -- Ensure unique statement per stock per period
);

-- Create the analyst_targets table
CREATE TABLE analyst_targets (
    target_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),    -- Unique identifier for the analyst target record
    stock_id UUID NOT NULL REFERENCES stocks(stock_id),      -- Foreign key to stocks table
    date DATE NOT NULL,                                      -- Date of the analyst consensus
    consensus_price_target NUMERIC(10, 2),                   -- Average price target across analysts
    high_price_target NUMERIC(10, 2),                        -- Highest price target
    low_price_target NUMERIC(10, 2),                         -- Lowest price target
    consensus_rating TEXT,                                   -- Overall rating category (e.g., 'strong_buy', 'hold')
    consensus_rating_value NUMERIC(5, 2),                    -- Numerical average of consensus weights (1=Strong Sell, 5=Strong Buy)
    buy_ratings_count INTEGER,                               -- Count of 'Buy' ratings
    hold_ratings_count INTEGER,                              -- Count of 'Hold' ratings
    sell_ratings_count INTEGER,                              -- Count of 'Sell' ratings
    total_analysts_contributing INTEGER,                     -- Total unique analysts contributing
    created_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of record creation
    updated_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of last record update
    UNIQUE (stock_id, date)                                  -- Ensure unique consensus per stock per date
);

-- Create the sentiment_scores table
CREATE TABLE sentiment_scores (
    sentiment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for the sentiment record
    stock_id UUID NOT NULL REFERENCES stocks(stock_id),      -- Foreign key to stocks table
    timestamp TIMESTAMPTZ NOT NULL,                          -- Timestamp of the sentiment data
    absolute_index DOUBLE PRECISION,                         -- How much people are talking about the stock
    relative_index DOUBLE PRECISION,                         -- Discussion volume relative to previous day
    sentiment_score DOUBLE PRECISION,                        -- Overall percentage of positive activity (e.g., 0 to 100 or -1 to 1)
    general_perception TEXT,                                 -- Overall perception (e.g., 'positive', 'negative', 'neutral')
    source TEXT,                                             -- Source of sentiment (e.g., 'Reddit', 'Twitter', 'Overall')
    created_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of record creation
    updated_at TIMESTAMPTZ DEFAULT NOW(),                    -- Timestamp of last record update
    UNIQUE (stock_id, timestamp, source)                     -- Ensure unique sentiment per stock per timestamp per source
);
