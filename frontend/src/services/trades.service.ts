import { api } from '@/lib/api';
import { TradeWithDetails } from '@/types';

export const tradesService = {
  async getMyTrades(): Promise<TradeWithDetails[]> {
    const response = await api.get<TradeWithDetails[]>('/trades/my');
    return response.data;
  },

  async getAllTrades(): Promise<TradeWithDetails[]> {
    const response = await api.get<TradeWithDetails[]>('/trades');
    return response.data;
  },
};
