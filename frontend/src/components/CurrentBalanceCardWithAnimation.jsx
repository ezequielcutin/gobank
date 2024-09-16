import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { getUserIdByEmail } from '../services/api';

let gradientAngle = 0;

const CurrentBalanceCard = () => {
  const [balance, setBalance] = useState(null);

  useEffect(() => {
    //-------------------------- backend stuff --------------------------
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


    //-------------------------- front end stuff --------------------------
    const gradientVar = document.querySelector(".text-white.rounded-lg.shadow-lg.p-6");

    if (gradientVar) {
      const handleMouseMove = (event) => {
        const rect = gradientVar.getBoundingClientRect();  // Get element's size and position

        // Get normalized mouse position relative to the element (0 to 1 scale)
        const normalizedX = (((event.clientX - rect.left) / rect.width) - 0.5) * 2;
        const normalizedY = (((event.clientY - rect.top) / rect.height) + 0.5) * 2;

        // Log corner outputs based on mouse position
        console.log(`Mouse Position (normalized): (${normalizedX.toFixed(2)}, ${normalizedY.toFixed(2)})`);

        // Calculate theta (angle) using atan2 for proper angle calculation
        var theta = Math.atan2(normalizedY, normalizedX);
        console.log("Theta (angle in radians): ", theta);

        // Convert radians to degrees for the gradient
        gradientAngle = theta * (180 / Math.PI);
        console.log("Gradient Angle (in degrees):", gradientAngle);

        // Apply the gradient directly to the element's background
        gradientVar.style.background = `linear-gradient(${gradientAngle}deg, red, blue, yellow)`;
      };

      gradientVar.addEventListener('mousemove', handleMouseMove);

      // Clean up the event listener on component unmount
      return () => {
        gradientVar.removeEventListener('mousemove', handleMouseMove);
      };
    }

  }, []);


  return (
    <div
      className="text-white rounded-lg shadow-lg p-6"
    >
      <h2 className="text-2xl font-bold mb-4">Current Balance</h2>
      {balance === null ? (
        <div className="h-12 w-24 bg-white/20 animate-pulse rounded"></div>
      ) : (
        <p className="text-4xl font-bold">${balance.toFixed(2)}</p>
      )}
    </div>
  );
}
export default CurrentBalanceCard;
