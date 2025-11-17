import axios, { AxiosError } from 'axios';
import https from 'https';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'https://localhost:8443/api';

// In development, allow self-signed certificates
// IMPORTANT: This should ONLY be used in development, never in production
const httpsAgent = process.env.NODE_ENV !== 'production'
  ? new https.Agent({ rejectUnauthorized: false })
  : undefined;

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  // Only applied on server-side (Node.js), not in browser
  httpsAgent,
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors and unwrap data
api.interceptors.response.use(
  (response) => {
    // Unwrap the backend's standard response format { success, message, data }
    // Extract the actual data from response.data.data if success is true
    if (response.data && typeof response.data === 'object' && response.data.success && 'data' in response.data) {
      response.data = response.data.data;
    }
    return response;
  },
  (error: AxiosError<any>) => {
    if (error.response?.status === 401) {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const handleApiError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const responseData = error.response?.data;

    // Try to extract error message from various possible formats
    if (typeof responseData === 'string') {
      return responseData;
    }

    if (responseData && typeof responseData === 'object') {
      return responseData.message || responseData.error || responseData.errors || error.message;
    }

    return error.message || 'An unexpected error occurred';
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'An unexpected error occurred';
};
