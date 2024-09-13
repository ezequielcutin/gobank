// frontend/src/components/TransactionHistoryCard.jsx

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { getUserIdByEmail } from '../services/api';


const TransactionHistoryCard = () => {
  const [transactions, setTransactions] = useState(null);

  useEffect(() => {
    const fetchTransactions = async () => {
      const userEmail = localStorage.getItem('userEmail');
      console.log("User Email Retrieved for getTransfers: ", userEmail);
      if (userEmail) {
        try {
          const userId = await getUserIdByEmail(userEmail);
          console.log("User ID Retrieved for getTransfers: ", userId);
          const response = await axios.get(`http://localhost:3000/transactions/${userId}`);
          setTransactions(response.data);
        } catch (error) {
          console.error('Error fetching transactions:', error);
          setTransactions([]);
        }
      }
    };

    fetchTransactions();
  }, []);

  if (transactions === null) {
    return <div className="card">Loading transactions...</div>;
  }

  return (
    <div className="card">
      <h2>Transaction History</h2>
      {transactions.length === 0 ? (
        <p>No transactions found.</p>
      ) : (
        <ul>
          {transactions.map((transaction) => (
            <li key={transaction.id}>
              <p>Amount: {transaction.amount}</p>
              <p>Type: {transaction.type}</p>
              <p>Date: {new Date(transaction.created_at).toLocaleString()}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TransactionHistoryCard;