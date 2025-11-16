'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, LoginRequest, RegisterRequest, TwoFAVerifyRequest } from '@/types';
import { authService } from '@/services/auth.service';
import { handleApiError } from '@/lib/api';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (data: LoginRequest) => Promise<{ requires2FA?: boolean }>;
  verify2FA: (data: TwoFAVerifyRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
  setToken: (token: string) => void;
  setUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is logged in on mount
    const initAuth = async () => {
      const token = localStorage.getItem('token');
      const storedUser = localStorage.getItem('user');

      if (token && storedUser) {
        try {
          const parsedUser = JSON.parse(storedUser);
          setUser(parsedUser);
          // Don't refresh on mount to avoid logout loop
          // refreshUser();
        } catch (error) {
          console.error('Failed to parse stored user:', error);
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
      }
      setLoading(false);
    };

    initAuth();
  }, []);

  const login = async (data: LoginRequest): Promise<{ requires2FA?: boolean; success?: boolean }> => {
    try {
      const response = await authService.login(data);

      if (response.requires_2fa) {
        return { requires2FA: true };
      }

      if (response.token && response.user) {
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        setUser(response.user);
        return { success: true };
      } else {
        throw new Error('Invalid response from server');
      }
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  };

  const verify2FA = async (data: TwoFAVerifyRequest): Promise<void> => {
    try {
      const response = await authService.verify2FA(data);

      if (response.token && response.user) {
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        setUser(response.user);
      }
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  };

  const register = async (data: RegisterRequest): Promise<void> => {
    try {
      await authService.register(data);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  };

  const refreshUser = async () => {
    try {
      const userData = await authService.getMe();
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } catch (error) {
      console.error('Failed to refresh user:', error);
      logout();
    }
  };

  const handleSetToken = (token: string) => {
    localStorage.setItem('token', token);
  };

  const handleSetUser = (userData: User) => {
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
  };

  return (
    <AuthContext.Provider value={{
      user,
      loading,
      login,
      verify2FA,
      register,
      logout,
      refreshUser,
      setToken: handleSetToken,
      setUser: handleSetUser
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
