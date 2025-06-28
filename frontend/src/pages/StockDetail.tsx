import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { fetchHistoricalPrices } from '../api/api';
import HistoricalPriceChart from '../components/HistoricalPriceChart';

interface HistoricalPrice {
  time: string;
  open_price: number;
  high_price: number;
  low_price: number;
  close_price: number;
  volume: number;
}

const StockDetail: React.FC = () => {
  const { symbol } = useParams<{ symbol: string }>();
  const [prices, setPrices] = useState<HistoricalPrice[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getHistoricalPrices = async () => {
      if (!symbol) return;
      try {
        const data = await fetchHistoricalPrices(symbol);
        setPrices(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    getHistoricalPrices();
  }, [symbol]);

  if (loading) return <div>Loading {symbol} historical data...</div>;
  if (error) return <div>Error: {error}</div>;
  if (prices.length === 0) return <div>No historical data available for {symbol}.</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">{symbol} Historical Prices</h1>
      <HistoricalPriceChart prices={prices} symbol={symbol!} />
    </div>
  );
};

export default StockDetail;
