// src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:3000/'; // Ensure this matches your backend

export const getUsers = () => {
    return axios.get(`${API_URL}`);
};

export const createUser = async (firstName, lastName, email, password) => {
    console.log('createUser called with:', { firstName, lastName, email, password });
    try {
      const response = await axios.post(`${API_URL}register`, { firstName, lastName, email, password });
      localStorage.setItem('userEmail', response.data.user.email);
      console.log("Email stored in local storage: ", response.data.user.email);
      console.log('createUser response:', response);
      return response;
    } catch (error) {
      console.error('createUser error:', error);
      throw error;
    }
  };

export const transferFunds = (fromId, toId, amount) => {
    return axios.post(`${API_URL}transfer`, { fromId, toId, amount });
};

export const deleteUser = (id) => {
    return axios.delete(`${API_URL}account${id}`);
};

export const loginUser = async (email, password) => {
  try {
    const response = await axios.post(`${API_URL}login`, { email, password });
    if (response.data.user && response.data.user.id) {
      localStorage.setItem('userEmail', response.data.user.email);
      console.log("Email stored in local storage: ", response.data.user.email);
    }
    return response;
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
};



export const getUserIdByEmail = async (email) => {
  try {
    const response = await axios.get(`${API_URL}user-by-email/${email}`);
    return response.data.id;
  } catch (error) {
    console.error('Error fetching user ID:', error);
    throw error;
  }
};