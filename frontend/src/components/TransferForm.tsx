import React, { useState, useEffect } from 'react';
import { transferFunds } from '../services/api';
import { CreditCard, DollarSign } from 'lucide-react';

const TransferForm: React.FC<{ onTransfer: () => void }> = ({ onTransfer }) => {
  const [fromId, setFromId] = useState('');
  const [toId, setToId] = useState('');
  const [amount, setAmount] = useState('');

  useEffect(() => {
    const storedUserId = localStorage.getItem('userId');
    if (storedUserId) {
      setFromId(storedUserId);
    }
  }, []);

  const handleTransfer = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      await transferFunds(fromId, parseInt(toId), parseInt(amount));
      onTransfer();
      setToId('');
      setAmount('');
    } catch (error) {
      console.error('Transfer failed:', error);
    }
  };

  return (
    <form onSubmit={handleTransfer} className="space-y-4">
      <div className="relative">
        <label htmlFor="fromId" className="block text-sm font-medium text-gray-700 mb-1">From Account ID</label>
        <div className="relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <CreditCard className="h-5 w-5 text-gray-400" />
          </div>
          <input
            id="fromId"
            type="number"
            value={fromId}
            onChange={(e) => setFromId(e.target.value)}
            required
            className="block w-full pl-10 sm:text-sm border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
            placeholder="123"
          />
        </div>
      </div>
      <div className="relative">
        <label htmlFor="toId" className="block text-sm font-medium text-gray-700 mb-1">To Account ID</label>
        <div className="relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <CreditCard className="h-5 w-5 text-gray-400" />
          </div>
          <input
            id="toId"
            type="number"
            value={toId}
            onChange={(e) => setToId(e.target.value)}
            required
            className="block w-full pl-10 sm:text-sm border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
            placeholder="456"
          />
        </div>
      </div>
      <div className="relative">
        <label htmlFor="amount" className="block text-sm font-medium text-gray-700 mb-1">Amount</label>
        <div className="relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <DollarSign className="h-5 w-5 text-gray-400" />
          </div>
          <input
            id="amount"
            type="number"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            required
            className="block w-full pl-10 sm:text-sm border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
            placeholder="100"
          />
        </div>
      </div>
      <button
        type="submit"
        className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition duration-150 ease-in-out"
      >
        Transfer
      </button>
    </form>
  );
};

export default TransferForm;