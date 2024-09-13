import React from 'react';
import { useNavigate } from 'react-router-dom';
import CurrentBalanceCard from './CurrentBalanceCard';
import TransferFundsCard from './TransferFundsCard';
import TransactionHistoryCard from './TransactionHistoryCard';


const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    // Clear the token and user ID from local storage
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    // Redirect to the login page
    navigate('/login');
  };

  return (
    <div>
      <h1>Dashboard</h1>
      <button
        onClick={handleLogout}
        className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
      >
        Logout
      </button>
      <div className="grid grid-cols-3 gap-4">
        <CurrentBalanceCard />
        <TransferFundsCard />
        <TransactionHistoryCard />
      </div>
    </div>
  );
};

export default Dashboard;