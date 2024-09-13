// frontend/src/components/TransferFundsCard.jsx

import React, { useState } from 'react';
import axios from 'axios';

const TransferFundsCard = () => {
  const [amount, setAmount] = useState('');
  const [recipient, setRecipient] = useState('');

  const handleTransfer = async (e) => {
    e.preventDefault();
    
    try {
      // Make API call to transfer funds
      const response = await axios.post('/transfer', { amount, recipient });
      console.log('Transfer successful:', response.data);
      
      // Reset form fields after successful transfer
      setAmount('');
      setRecipient('');
    } catch (error) {
      console.error('Transfer failed:', error.message);
    }
  };

  return (
    <div className="card">
      <h2>Transfer Funds</h2>
      <form onSubmit={handleTransfer}>
        <input
          type="number"
          placeholder="Amount"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
        <input
          type="text"
          placeholder="Recipient"
          value={recipient}
          onChange={(e) => setRecipient(e.target.value)}
        />
        <button type="submit" className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
          Transfer
        </button>
      </form>
    </div>
  );
};

export default TransferFundsCard;