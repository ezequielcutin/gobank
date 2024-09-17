import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import CurrentBalanceCard from './CurrentBalanceCard';
import TransferFundsCard from './TransferFundsCard';
import TransactionHistoryCard from './TransactionHistoryCard';
import AvailableUsersCard from './AvailableUsersCard';
import { getUserDetails } from '../services/api';
import { LogOutIcon } from 'lucide-react';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [userName, setUserName] = useState('');
  const [currentUserId, setCurrentUserId] = useState<number | null>(null);

  useEffect(() => {
    const fetchUserDetails = async () => {
      const userEmail = localStorage.getItem('userEmail');
      if (userEmail) {
        try {
          const userDetails = await getUserDetails(userEmail);
          setUserName(`${userDetails.firstName} ${userDetails.lastName}`);
          setCurrentUserId(userDetails.id);
        } catch (error) {
          console.error('Error fetching user details:', error);
        }
      }
    };

    fetchUserDetails();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userEmail');
    navigate('/');
  };

  const handleTransferSuccess = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-100 p-4 sm:p-6 lg:p-8">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white shadow-lg rounded-2xl overflow-hidden">
          <div className="p-6 sm:p-10">
            <header className="flex justify-between items-center mb-8">
              <h1 className="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800">
                {userName ? `${userName}'s Finances` : 'Your Finances'}
              </h1>
              <button
                onClick={handleLogout}
                className="flex items-center px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors duration-300"
              >
                <LogOutIcon className="w-4 h-4 mr-2" />
                Logout
              </button>
            </header>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div className="md:col-span-2 lg:col-span-1">
                <CurrentBalanceCard refreshTrigger={refreshTrigger} />
              </div>
              <div className="md:col-span-2 lg:col-span-1">
                <TransferFundsCard onTransferSuccess={handleTransferSuccess} />
              </div>
              <div className="md:col-span-2 lg:col-span-1">
                {currentUserId !== null && <AvailableUsersCard currentUserId={currentUserId} />}
              </div>
              <div className="md:col-span-2 lg:col-span-3">
                <TransactionHistoryCard refreshTrigger={refreshTrigger} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;