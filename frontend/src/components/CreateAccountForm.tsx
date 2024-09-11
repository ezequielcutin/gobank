import React, { useState } from 'react';
import { createAccount } from '../services/api';
import { User, Mail } from 'lucide-react';

const CreateAccountForm: React.FC<{ onSuccess: () => void }> = ({ onSuccess }) => {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');

  const handleCreateAccount = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      await createAccount(firstName, lastName);
      onSuccess();
      setFirstName('');
      setLastName('');
    } catch (error) {
      console.error('Account creation failed:', error);
    }
  };

  return (
    <form onSubmit={handleCreateAccount} className="space-y-4">
      <div className="relative">
        <label htmlFor="firstName" className="block text-sm font-medium text-gray-700 mb-1">First Name</label>
        <div className="relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <User className="h-5 w-5 text-gray-400" />
          </div>
          <input
            id="firstName"
            type="text"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            required
            className="block w-full pl-10 sm:text-sm border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            placeholder="John"
          />
        </div>
      </div>
      <div className="relative">
        <label htmlFor="lastName" className="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
        <div className="relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Mail className="h-5 w-5 text-gray-400" />
          </div>
          <input
            id="lastName"
            type="text"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            required
            className="block w-full pl-10 sm:text-sm border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            placeholder="Doe"
          />
        </div>
      </div>
      <button
        type="submit"
        className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-150 ease-in-out"
      >
        Create Account
      </button>
    </form>
  );
};

export default CreateAccountForm;