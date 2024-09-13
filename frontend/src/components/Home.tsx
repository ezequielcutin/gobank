import React from 'react';
import { useNavigate } from 'react-router-dom';

const Home: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <h1 className="text-4xl font-bold mb-8">Welcome to GoBank</h1>
      <div className="space-y-4">
        <button
          onClick={() => navigate('/login')}
          className="w-64 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition duration-300"
        >
          Login
        </button>
        <button
          onClick={() => navigate('/register')}
          className="w-64 bg-green-500 text-white py-2 px-4 rounded hover:bg-green-600 transition duration-300"
        >
          Register
        </button>
        <button
          onClick={() => navigate('/admin')}
          className="w-64 bg-green-500 text-white py-2 px-4 rounded hover:bg-green-600 transition duration-300"
        >
          Admin View (Monopoly Bank)
        </button>
      </div>
    </div>
  );
};

export default Home;