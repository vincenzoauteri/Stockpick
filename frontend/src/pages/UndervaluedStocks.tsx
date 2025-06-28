import React, { useEffect, useState } from 'react';
import { fetchUndervaluedStocks } from '../api/api';
import { Link } from 'react-router-dom';

interface UndervaluationScore {
  stock_id: string;
  symbol: string;
  score: number;
  fundamental_score: number;
  analyst_score: number;
  sentiment_score: number;
}

const UndervaluedStocks: React.FC = () => {
  const [undervalued, setUndervalued] = useState<UndervaluationScore[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getUndervaluedStocks = async () => {
      try {
        const data = await fetchUndervaluedStocks();
        setUndervalued(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    getUndervaluedStocks();
  }, []);

  if (loading) return <div>Loading undervalued stocks...</div>;
  if (error) return <div>Error: {error}</div>;
  if (undervalued.length === 0) return <div>No undervalued stocks found at this time.</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Undervalued Stocks</h1>
      <div className="overflow-x-auto">
        <table className="min-w-full bg-white shadow-md rounded-lg overflow-hidden">
          <thead className="bg-gray-800 text-white">
            <tr>
              <th className="py-2 px-4 text-left">Symbol</th>
              <th className="py-2 px-4 text-left">Undervaluation Score</th>
              <th className="py-2 px-4 text-left">Fundamental Score</th>
              <th className="py-2 px-4 text-left">Analyst Score</th>
              <th className="py-2 px-4 text-left">Sentiment Score</th>
            </tr>
          </thead>
          <tbody>
            {undervalued.map((stock) => (
              <tr key={stock.stock_id} className="border-b border-gray-200 hover:bg-gray-100">
                <td className="py-2 px-4"><Link to={`/stocks/${stock.symbol}`} className="text-blue-600 hover:underline">{stock.symbol}</Link></td>
                <td className="py-2 px-4">{stock.score.toFixed(2)}</td>
                <td className="py-2 px-4">{stock.fundamental_score.toFixed(2)}</td>
                <td className="py-2 px-4">{stock.analyst_score.toFixed(2)}</td>
                <td className="py-2 px-4">{stock.sentiment_score.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UndervaluedStocks;
