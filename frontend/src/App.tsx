import React, { useState } from 'react';
import { DollarSign, PlusCircle, ArrowRightCircle } from 'lucide-react';
import AccountList from './components/AccountList';
import TransferForm from './components/TransferForm';
import CreateAccountForm from './components/CreateAccountForm';
import { CheckCircle, XCircle } from 'lucide-react';

type NotificationType = {
  type: 'success' | 'error';
  message: string;
} | null;

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [notification, setNotification] = useState<NotificationType>(null);

  const refreshAccounts = () => {
    setRefreshKey(prevKey => prevKey + 1);
  };

  const showNotification = (type: 'success' | 'error', message: string) => {
    setNotification({ type, message });
    setTimeout(() => setNotification(null), 5000);
  };

  return (
    <div className="min-h-screen bg-gray-50 font-sans">
      <header className="bg-gradient-to-r from-blue-600 to-blue-800 shadow-lg">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 flex items-center justify-between">
          <div className="flex items-center">
            <img src="/gobank-removebg.png" alt="GoBank Logo" className="h-12 w-12 mr-2" />
            <h1 className="text-3xl font-bold text-white">GoBank</h1>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {notification && (
          <div 
            className={`mb-4 p-4 rounded-md ${
              notification.type === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            } flex items-center`} 
            role="alert"
          >
            {notification.type === 'success' ? (
              <CheckCircle className="h-5 w-5 mr-2" />
            ) : (
              <XCircle className="h-5 w-5 mr-2" />
            )}
            <p className="font-medium">{notification.message}</p>
          </div>
        )}
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="bg-gradient-to-r from-blue-500 to-blue-600 p-4">
              <h2 className="text-xl font-semibold text-white flex items-center">
                <PlusCircle className="h-5 w-5 mr-2" />
                Create Account
              </h2>
            </div>
            <div className="p-6">
              <CreateAccountForm onSuccess={() => {
                refreshAccounts();
                showNotification('success', 'Account created successfully!');
              }} />
            </div>
          </div>
          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="bg-gradient-to-r from-green-500 to-green-600 p-4">
              <h2 className="text-xl font-semibold text-white flex items-center">
                <ArrowRightCircle className="h-5 w-5 mr-2" />
                Transfer Funds
              </h2>
            </div>
            <div className="p-6">
              <TransferForm onTransfer={() => {
                refreshAccounts();
                showNotification('success', 'Transfer successful!');
              }} />
            </div>
          </div>
        </div>
        <div className="mt-8 bg-white shadow-md rounded-lg overflow-hidden">
          <div className="bg-gradient-to-r from-purple-500 to-purple-600 p-4">
            <h2 className="text-xl font-semibold text-white flex items-center">
              <DollarSign className="h-5 w-5 mr-2" />
              Accounts
            </h2>
          </div>
          <div className="p-6">
            <AccountList 
              key={refreshKey} 
              onDelete={() => {
                refreshAccounts();
                showNotification('success', 'Account deleted successfully!');
              }} 
            />
          </div>
        </div>
      </main>
    </div>
  );
}