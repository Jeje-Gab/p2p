import { api } from '@/lib/api';
import {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User,
  TwoFASetupResponse,
  TwoFAVerifyRequest,
} from '@/types';

export const authService = {
  async register(data: RegisterRequest): Promise<User> {
    const response = await api.post<User>('/auth/register', data);
    return response.data;
  },

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  async verify2FA(data: TwoFAVerifyRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/auth/2fa/verify', data);
    return response.data;
  },

  async getMe(): Promise<User> {
    const response = await api.get<User>('/auth/me');
    return response.data;
  },

  async setup2FA(): Promise<TwoFASetupResponse> {
    const response = await api.post<TwoFASetupResponse>('/auth/2fa/setup');
    return response.data;
  },

  async enable2FA(code: string): Promise<{ message: string }> {
    const response = await api.post('/auth/2fa/enable', { code });
    return response.data;
  },

  async disable2FA(code: string): Promise<{ message: string }> {
    const response = await api.post('/auth/2fa/disable', { code });
    return response.data;
  },

  async getSteamLoginUrl(): Promise<string> {
    const response = await api.get<{ auth_url: string }>('/auth/steam/login');
    return response.data.auth_url;
  },

  async handleSteamCallback(params: Record<string, string>): Promise<LoginResponse> {
    const queryString = new URLSearchParams(params).toString();
    const response = await api.get<LoginResponse>(`/auth/steam/callback?${queryString}`);
    return response.data;
  },
};
