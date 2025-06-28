import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import StockList from './components/StockList';
import StockDetail from './pages/StockDetail';
import UndervaluedStocks from './pages/UndervaluedStocks';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <nav className="bg-gray-800 p-4">
          <div className="container mx-auto flex justify-between items-center">
            <Link to="/" className="text-white text-xl font-bold">StockPick</Link>
            <div className="space-x-4">
              <Link to="/stocks" className="text-gray-300 hover:text-white">Stocks</Link>
              <Link to="/undervalued" className="text-gray-300 hover:text-white">Undervalued</Link>
            </div>
          </div>
        </nav>

        <main className="p-4">
          <Routes>
            <Route path="/" element={<StockList />} />
            <Route path="/stocks" element={<StockList />} />
            <Route path="/stocks/:symbol" element={<StockDetail />} />
            <Route path="/undervalued" element={<UndervaluedStocks />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;