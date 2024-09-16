import React, { useState, useEffect } from 'react';
import { transferFunds, getUserIdByEmail } from '../services/api';

const TransferFundsCard = ({ onTransferSuccess }) => {
  const [amount, setAmount] = useState('');
  const [toId, setToId] = useState('');
  const [fromId, setFromId] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    const fetchUserId = async () => {
      const userEmail = localStorage.getItem('userEmail');
      if (userEmail) {
        try {
          const userId = await getUserIdByEmail(userEmail);
          setFromId(userId);
        } catch (error) {
          console.error('Error fetching user ID:', error);
        }
      }
    };

    fetchUserId();
  }, []);

  const handleTransfer = async (e) => {
    e.preventDefault();
    
    try {
      await transferFunds(fromId, parseInt(toId), parseInt(amount));
      setMessage('Transfer successful');
      onTransferSuccess(); // Call this function to update parent components
      
      // Reset form fields after successful transfer
      setAmount('');
      setToId('');
    } catch (error) {
      setMessage('Transfer failed: ' + (error.response ? error.response.data.Error : error.message));
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-4">
        <h2 className="text-2xl font-bold text-white">Transfer Funds</h2>
      </div>
      <div className="p-6">
        {message && (
          <p className={`mb-4 p-2 rounded ${message.includes('successful') ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
            {message}
          </p>
        )}
        <form onSubmit={handleTransfer} className="space-y-4">
          <div>
            <label htmlFor="amount" className="block text-sm font-medium text-gray-700 mb-1">Amount</label>
            <input
              id="amount"
              type="number"
              placeholder="Enter amount"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
            />
          </div>
          <div>
            <label htmlFor="toId" className="block text-sm font-medium text-gray-700 mb-1">To User ID</label>
            <input
              id="toId"
              type="number"
              placeholder="Enter recipient's User ID"
              value={toId}
              onChange={(e) => setToId(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
            />
          </div>
          <button 
            type="submit" 
            className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white p-2 rounded hover:from-blue-600 hover:to-purple-700 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Transfer
          </button>
        </form>
      </div>
    </div>
  );
};

export default TransferFundsCard;