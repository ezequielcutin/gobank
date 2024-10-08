// frontend/src/components/TransactionHistoryCard.jsx

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { getUserIdByEmail } from '../services/api';
import BalanceChart from './BalanceChart';

const TransactionHistoryCard = ({ refreshTrigger }) => {
  const [transactions, setTransactions] = useState(null);
  const [balanceHistory, setBalanceHistory] = useState([]);

  useEffect(() => {
    const fetchTransactions = async () => {
      const userEmail = localStorage.getItem('userEmail');
      if (userEmail) {
        try {
          const userId = await getUserIdByEmail(userEmail);
          const response = await axios.get(`http://localhost:3000/transactions/${userId}`);
          setTransactions(response.data);
          calculateBalanceHistory(response.data);
        } catch (error) {
          console.error('Error fetching transactions:', error);
          setTransactions([]);
        }
      }
    };

    fetchTransactions();
  }, [refreshTrigger]);

  const calculateBalanceHistory = (transactions) => {
    let balance = 999999999; // Starting balance
    const history = transactions.map(transaction => {
      balance += transaction.amount;
      return {
        date: new Date(transaction.createdAt),
        balance: balance
      };
    });
    setBalanceHistory(history);
  };

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-4">
        <h2 className="text-2xl font-bold text-white">Transaction History</h2>
      </div>
      <div className="p-6">
        {transactions === null ? (
          <p className="text-gray-600 italic">Loading transactions...</p>
        ) : transactions.length === 0 ? (
          <p className="text-gray-600 italic">No transactions found.</p>
        ) : (
          <>
            <div className="overflow-x-auto mb-6">
              <table className="w-full">
                <thead>
                  <tr className="bg-gray-100">
                    <th className="py-2 px-4 text-left font-semibold text-gray-600">Amount</th>
                    <th className="py-2 px-4 text-left font-semibold text-gray-600">Type</th>
                    <th className="py-2 px-4 text-left font-semibold text-gray-600">Date</th>
                  </tr>
                </thead>
                <tbody>
                  {transactions.map((transaction, index) => {
                    const date = new Date(transaction.createdAt.replace('Z', ''));
                    return (
                      <tr
                        key={transaction.id}
                        className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition-colors duration-150 ease-in-out`}
                      >
                        <td className="py-3 px-4">
                          <span className={transaction.amount < 0 ? 'text-red-600' : 'text-green-600'}>
                            ${Math.abs(transaction.amount).toFixed(2)}
                          </span>
                        </td>
                        <td className="py-3 px-4">
                          {transaction.type}
                        </td>
                        <td className="py-3 px-4">{date.toLocaleString()}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
            {balanceHistory.length > 0 && (
              <div className="mt-8">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-4 rounded-t-lg">
                  <h2 className="text-2xl font-bold text-white">Balance Over Time</h2>
                </div>
                <BalanceChart balanceHistory={balanceHistory} />
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default TransactionHistoryCard;