import { api } from '@/lib/api';
import { Skin, UserSkinWithDetails } from '@/types';

export const skinsService = {
  async getAllSkins(): Promise<Skin[]> {
    const response = await api.get<Skin[]>('/skins');
    return response.data;
  },

  async getMySkins(): Promise<UserSkinWithDetails[]> {
    const response = await api.get<UserSkinWithDetails[]>('/skins/my');
    return response.data;
  },

  async addSkinToInventory(skinId: number, quantity: number = 1): Promise<void> {
    await api.post('/skins/my', { skin_id: skinId, quantity });
  },

  async removeSkinFromInventory(skinId: number, quantity: number = 1): Promise<void> {
    await api.delete(`/skins/my/${skinId}`, { data: { quantity } });
  },
};
