import React, {useState, useEffect} from 'react';
import { useNavigate } from 'react-router-dom';
import CurrentBalanceCard from './CurrentBalanceCard.jsx';
import TransferFundsCard from './TransferFundsCard.jsx';
import TransactionHistoryCard from './TransactionHistoryCard.jsx';
import AvailableUsersCard from './AvailableUsersCard';
import { getUserDetails } from '../services/api';

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
    <div className="min-h-screen bg-gray-100 py-6 flex flex-col justify-center sm:py-12">
      <div className="relative py-3 sm:max-w-xl md:max-w-2xl lg:max-w-4xl xl:max-w-6xl mx-auto">
        <div className="relative px-4 py-10 bg-white shadow-lg sm:rounded-3xl sm:p-20">
          <div className="flex justify-between items-center mb-10">
          <h1 className="text-4xl font-bold text-gray-800">{userName ? `${userName}'s Finances` : 'Your Finances'}</h1>
          <button
              onClick={handleLogout}
              className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition duration-300"
            >
              Logout
            </button>
          </div>
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
  );
};

export default Dashboard;