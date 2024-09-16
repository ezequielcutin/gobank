import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { loginUser } from '../services/api';

const AdminView: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const adminEmail = 'admin@gmail.com'; // Monopoly Bank email
    const adminPassword = 'megapassword'; // Monopoly Bank password

    const performLogin = async () => {
      try {
        const response = await loginUser(adminEmail, adminPassword);
        // Assuming loginUser stores token and user info in localStorage
        navigate('/dashboard');
      } catch (error) {
        console.error('Admin login failed:', error);
      }
    };

    performLogin();
  }, [navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <p>Logging in as Admin...</p>
    </div>
  );
};

export default AdminView;