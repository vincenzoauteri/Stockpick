import React, { useEffect, useState } from 'react';
import { fetchStocks, ingestHistoricalPrices } from '../api/api';
import { Link } from 'react-router-dom';

interface Stock {
  stock_id: string;
  symbol: string;
  company_name: string;
  exchange: string;
}

const StockList: React.FC = () => {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [ingesting, setIngesting] = useState<string | null>(null);

  useEffect(() => {
    const getStocks = async () => {
      try {
        const data = await fetchStocks();
        setStocks(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    getStocks();
  }, []);

  const handleIngest = async (symbol: string) => {
    setIngesting(symbol);
    try {
      await ingestHistoricalPrices(symbol);
      alert(`Successfully ingested data for ${symbol}`);
    } catch (err: any) {
      alert(`Failed to ingest data for ${symbol}: ${err.message}`);
    } finally {
      setIngesting(null);
    }
  };

  if (loading) return <div>Loading stocks...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Available Stocks</h1>
      <div className="overflow-x-auto">
        <table className="min-w-full bg-white shadow-md rounded-lg overflow-hidden">
          <thead className="bg-gray-800 text-white">
            <tr>
              <th className="py-2 px-4 text-left">Symbol</th>
              <th className="py-2 px-4 text-left">Company Name</th>
              <th className="py-2 px-4 text-left">Exchange</th>
              <th className="py-2 px-4 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {stocks.map((stock) => (
              <tr key={stock.stock_id} className="border-b border-gray-200 hover:bg-gray-100">
                <td className="py-2 px-4"><Link to={`/stocks/${stock.symbol}`} className="text-blue-600 hover:underline">{stock.symbol}</Link></td>
                <td className="py-2 px-4">{stock.company_name}</td>
                <td className="py-2 px-4">{stock.exchange}</td>
                <td className="py-2 px-4">
                  <button
                    onClick={() => handleIngest(stock.symbol)}
                    disabled={ingesting === stock.symbol}
                    className={`bg-green-500 text-white px-3 py-1 rounded-md text-sm ${ingesting === stock.symbol ? 'opacity-50 cursor-not-allowed' : 'hover:bg-green-600'}`}
                  >
                    {ingesting === stock.symbol ? 'Ingesting...' : 'Ingest Data'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default StockList;
