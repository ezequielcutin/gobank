import React, { useEffect, useState } from 'react';
import { getUsers, deleteUser } from '../services/api';
import { Trash2 } from 'lucide-react';

interface Account {
  id: number;
  firstName: string;
  lastName: string;
  balance: number;
  number: number; // Add account number to the interface
}

interface AccountListProps {
  onDelete: () => void;
}

const AccountList: React.FC<AccountListProps> = ({ onDelete }) => {
  const [accounts, setAccounts] = useState<Account[]>([]);

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const response = await getUsers();
        setAccounts(response.data);
      } catch (error) {
        console.error('Error fetching accounts:', error);
      }
    };

    fetchAccounts();
  }, []);

  const handleDelete = async (id: number) => {
    try {
      await deleteUser(id);
      setAccounts(accounts.filter(account => account.id !== id));
      onDelete();
    } catch (error) {
      console.error('Error deleting account:', error);
    }
  };

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Account Number</th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Balance</th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {accounts.map((account, index) => (
            <tr key={account.id} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{account.id}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{account.number}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{account.firstName} {account.lastName}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${account.balance.toFixed(2)}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button
                  onClick={() => handleDelete(account.id)}
                  className="text-red-600 hover:text-red-900 focus:outline-none focus:underline inline-flex items-center"
                >
                  <Trash2 className="h-4 w-4 mr-1" />
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AccountList;