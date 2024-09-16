import React from 'react';
import { useNavigate } from 'react-router-dom';

const Home: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-blue-100 to-purple-100">
      <div className="bg-white p-8 rounded-2xl shadow-xl max-w-md w-full">
        <h1 className="text-4xl font-bold mb-8 text-center bg-gradient-to-r from-blue-500 to-purple-600 bg-clip-text text-transparent">
          Welcome to GoBank
        </h1>
        <div className="space-y-4">
          <button
            onClick={() => navigate('/login')}
            className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 text-white py-3 px-4 rounded-lg hover:from-blue-600 hover:to-indigo-700 transition duration-300 shadow-md"
          >
            Login
          </button>
          <button
            onClick={() => navigate('/register')}
            className="w-full bg-gradient-to-r from-purple-500 to-pink-600 text-white py-3 px-4 rounded-lg hover:from-purple-600 hover:to-pink-700 transition duration-300 shadow-md"
          >
            Register
          </button>
          <button
            onClick={() => navigate('/admin')}
            className="w-full bg-gradient-to-r from-gray-700 to-gray-900 text-white py-3 px-4 rounded-lg hover:from-gray-800 hover:to-black transition duration-300 shadow-md"
          >
            Admin View (Monopoly Bank)
          </button>
        </div>
      </div>
    </div>
  );
};

export default Home;