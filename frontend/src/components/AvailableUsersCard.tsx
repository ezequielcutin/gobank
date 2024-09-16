import React, { useState, useEffect } from 'react';
import { getUsers } from '../services/api';

interface User {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  number: number;
}

const AvailableUsersCard: React.FC<{ currentUserId: number }> = ({ currentUserId }) => {
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        console.log('Fetching users...');
        const response = await getUsers();
        console.log('Users fetched:', response.data);
        const filteredUsers = response.data.filter((user: User) => user.id !== currentUserId);
        console.log('Filtered users:', filteredUsers);
        setUsers(filteredUsers);
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    };

    fetchUsers();
  }, [currentUserId]);

  console.log('Rendering AvailableUsersCard with users:', users);

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-4">
        <h2 className="text-2xl font-bold text-white">Available Users</h2>
      </div>
      <div className="p-6">
        {users.length === 0 ? (
          <p className="text-gray-600 italic">No other users available.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left bg-gray-100">
                  <th className="py-2 px-4 font-semibold text-gray-600">User ID</th>
                  <th className="py-2 px-4 font-semibold text-gray-600">Name</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user, index) => (
                  <tr 
                    key={user.id} 
                    className={`border-t ${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition-colors duration-150 ease-in-out`}
                  >
                    <td className="py-3 px-4">{user.id}</td>
                    <td className="py-3 px-4">{`${user.firstName} ${user.lastName}`}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default AvailableUsersCard;