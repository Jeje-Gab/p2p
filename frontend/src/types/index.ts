// User types
export interface User {
  id: number;
  email: string;
  steam_id?: string;
  twofa_enabled: boolean;
  role: 'user' | 'admin';
  created_at: string;
  updated_at: string;
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token?: string;
  requires_2fa?: boolean;
  user?: User;
}

export interface TwoFASetupResponse {
  secret: string;
  qr_code: string;
}

export interface TwoFAVerifyRequest {
  email: string;
  password: string;
  code: string;
}

// Skin types
export const SkinRarity = {
  ConsumerGrade: 'Consumer Grade',
  IndustrialGrade: 'Industrial Grade',
  MilSpecGrade: 'Mil-Spec Grade',
  Restricted: 'Restricted',
  Classified: 'Classified',
  Covert: 'Covert',
  ExceedinglyRare: 'Exceedingly Rare',
} as const;

export type SkinRarityType = typeof SkinRarity[keyof typeof SkinRarity];

export interface Skin {
  id: number;
  name: string;
  weapon_type: string;
  rarity: string;
  float_value?: number;
  image_url: string;
  created_at: string;
  updated_at: string;
}

export interface UserSkin {
  id: number;
  user_id: number;
  skin_id: number;
  quantity: number;
  created_at: string;
  updated_at: string;
}

export interface UserSkinWithDetails extends UserSkin {
  skin: Skin;
}

// Offer types
export const OfferStatus = {
  Open: 'open',
  Accepted: 'accepted',
  Cancelled: 'cancelled',
} as const;

export type OfferStatusType = typeof OfferStatus[keyof typeof OfferStatus];

export interface Offer {
  id: number;
  owner_id: number;
  skin_offered_id: number;
  skin_requested_id: number;
  status: OfferStatusType;
  created_at: string;
  updated_at: string;
}

export interface OfferWithDetails extends Offer {
  owner_email: string;
  skin_offered: Skin;
  skin_requested: Skin;
}

export interface CreateOfferRequest {
  skin_offered_id: number;
  skin_requested_id: number;
}

// Trade types
export interface Trade {
  id: number;
  offer_id: number;
  from_user_id: number;
  to_user_id: number;
  executed_at: string;
}

export interface TradeWithDetails extends Trade {
  offer: OfferWithDetails;
  from_user_email: string;
  to_user_email: string;
}
