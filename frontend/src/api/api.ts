const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

export const fetchStocks = async () => {
  const response = await fetch(`${API_BASE_URL}/stocks`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
};

export const fetchHistoricalPrices = async (symbol: string) => {
  const response = await fetch(`${API_BASE_URL}/stocks/${symbol}/history`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
};

export const fetchUndervaluedStocks = async () => {
  const response = await fetch(`${API_BASE_URL}/undervalued`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
};

export const ingestHistoricalPrices = async (symbol: string) => {
  const response = await fetch(`${API_BASE_URL}/ingest/historical-prices/${symbol}`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.text(); // Assuming backend returns a simple message
};
