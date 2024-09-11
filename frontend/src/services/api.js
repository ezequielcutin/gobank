// src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:3000/account'; // Ensure this matches your backend

export const getAccounts = () => {
    return axios.get(`${API_URL}`);
};

export const createAccount = (firstName, lastName) => {
    return axios.post(`${API_URL}`, { firstName, lastName });
};

export const transferFunds = (fromId, toId, amount) => {
    return axios.post('http://localhost:3000/transfer', { fromId, toId, amount });
};

export const deleteAccount = (id) => {
    return axios.delete(`http://localhost:3000/account/${id}`);
};
// Add more functions as needed