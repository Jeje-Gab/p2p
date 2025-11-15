import { api } from '@/lib/api';
import { OfferWithDetails, CreateOfferRequest } from '@/types';

export const offersService = {
  async getAllOffers(): Promise<OfferWithDetails[]> {
    const response = await api.get<OfferWithDetails[]>('/offers');
    return response.data;
  },

  async getMyOffers(): Promise<OfferWithDetails[]> {
    const response = await api.get<OfferWithDetails[]>('/offers/my');
    return response.data;
  },

  async getOfferById(id: number): Promise<OfferWithDetails> {
    const response = await api.get<OfferWithDetails>(`/offers/${id}`);
    return response.data;
  },

  async createOffer(data: CreateOfferRequest): Promise<OfferWithDetails> {
    const response = await api.post<OfferWithDetails>('/offers', data);
    return response.data;
  },

  async acceptOffer(id: number): Promise<void> {
    await api.post(`/offers/${id}/accept`);
  },

  async cancelOffer(id: number): Promise<void> {
    await api.post(`/offers/${id}/cancel`);
  },
};
