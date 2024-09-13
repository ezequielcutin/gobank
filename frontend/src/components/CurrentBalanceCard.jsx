// frontend/src/components/CurrentBalanceCard.jsx

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { getUserIdByEmail } from '../services/api';


const CurrentBalanceCard = () => {
  const [balance, setBalance] = useState(null);

  useEffect(() => {
    const fetchBalance = async () => {
      const userEmail = localStorage.getItem('userEmail');
      console.log("User Email Retrieved for getBalance: ", userEmail);
      if (userEmail) {
        try {
          const userId = await getUserIdByEmail(userEmail);
          console.log("User ID Retrieved for getBalance: ", userId);
          const response = await axios.get(`http://localhost:3000/balance/${userId}`);
          setBalance(response.data.balance);
        } catch (error) {
          console.error('Error fetching balance:', error);
          setBalance(0);
        }
      }
    };

    fetchBalance();
  }, []);

  return (
    <div className="card">
      <h2>Current Balance</h2>
      {balance === null ? <p>Loading balance...</p> : <p>${balance.toFixed(2)}</p>}
    </div>
  );
};

export default CurrentBalanceCard;