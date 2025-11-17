# CS2 P2P Skins Trading - Frontend

Frontend application for the CS2 P2P Skins trading platform, built with Next.js 14 and TypeScript.

## Features

- **Authentication System**
  - User registration and login
  - JWT token-based authentication
  - 2FA (TOTP) support
  - Steam OAuth integration (backend ready)

- **Skins Management**
  - View all available skins
  - Add skins to personal inventory
  - Remove skins from inventory
  - Search and filter skins
  - Rarity-based color coding

- **Trading System**
  - Browse all open trade offers
  - Create custom trade offers
  - Accept offers from other users
  - Cancel your own offers
  - View offer status

- **Trade History**
  - View completed trades
  - See trade details and participants
  - Track received and given items

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Forms**: React Hook Form
- **Icons**: Lucide React
- **State Management**: React Context API

## Getting Started

### Prerequisites

- Node.js 18+ installed
- Backend API running on **HTTPS** (see backend README)

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure environment variables:
   - The `.env.local` file is already configured
   - Default API URL: `https://localhost:8443/api` (**HTTPS**)

4. **IMPORTANTE: Aceite o certificado SSL auto-assinado**

   Antes de rodar o frontend, você precisa aceitar o certificado SSL do backend no navegador:

   a. Abra no navegador: **https://localhost:8443/healthz**

   b. Você verá um aviso de segurança. Clique em:
      - Chrome/Edge: "Advanced" → "Proceed to localhost (unsafe)"
      - Firefox: "Advanced" → "Accept the Risk and Continue"

   c. Você deve ver: `{"status":"healthy",...}`

   📖 **Detalhes completos:** Veja `SSL_SETUP.md` para mais informações

5. Run the development server:
   ```bash
   npm run dev
   ```

6. Open [http://localhost:3000](http://localhost:3000) in your browser

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── (dashboard)/       # Protected dashboard routes
│   │   │   ├── dashboard/     # Dashboard page
│   │   │   ├── skins/         # Skins management
│   │   │   ├── offers/        # Trade offers
│   │   │   └── trades/        # Trade history
│   │   ├── login/             # Login page
│   │   ├── register/          # Registration page
│   │   ├── layout.tsx         # Root layout
│   │   └── page.tsx           # Home page (redirects)
│   ├── components/            # Reusable components
│   │   ├── ui/               # UI components (Button, Input, Card)
│   │   ├── Navbar.tsx        # Navigation bar
│   │   └── SkinCard.tsx      # Skin display card
│   ├── contexts/             # React contexts
│   │   └── AuthContext.tsx   # Authentication state
│   ├── services/             # API service layers
│   │   ├── auth.service.ts
│   │   ├── skins.service.ts
│   │   ├── offers.service.ts
│   │   └── trades.service.ts
│   ├── lib/                  # Utilities
│   │   └── api.ts           # Axios configuration
│   └── types/               # TypeScript types
│       └── index.ts         # Shared types
├── public/                  # Static assets
├── .env.local              # Environment variables
├── next.config.mjs         # Next.js configuration
├── tailwind.config.ts      # Tailwind CSS configuration
└── package.json            # Dependencies
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## Features in Detail

### Authentication Flow

1. Users register with email and password
2. Login with credentials
3. If 2FA is enabled, enter TOTP code
4. JWT token is stored in localStorage
5. Token is automatically added to API requests

### Skins Management

- Browse all available CS2 skins
- Add skins to your inventory (simulated inventory)
- Remove skins from your inventory
- View skin details (name, weapon type, rarity, float value)
- Skins are color-coded by rarity

### Trading Flow

1. User creates an offer by selecting:
   - A skin they want to offer (from their inventory)
   - A skin they want to receive (from available skins)
2. Other users can browse all open offers
3. Users can accept offers if they have the requested skin
4. Once accepted, a trade is executed:
   - Skins are swapped between users
   - Trade is recorded in history
   - Offer status changes to "accepted"

## API Integration

The frontend communicates with the backend API at `https://localhost:8443/api` via **HTTPS with TLS encryption**. All authenticated requests include a JWT token in the Authorization header.

### Security Features

- ✅ **HTTPS/TLS** - All data encrypted in transit
- ✅ **JWT Authentication** - Secure token-based auth
- ✅ **Self-signed certificates** in development (production requires CA certificates)
- ✅ **Automatic SSL handling** - Next.js configured to accept dev certificates

### API Endpoints Used

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/2fa/verify` - 2FA verification
- `GET /auth/me` - Get current user
- `GET /skins` - Get all skins
- `GET /skins/my` - Get user's skins
- `POST /skins/my` - Add skin to inventory
- `DELETE /skins/my/:id` - Remove skin from inventory
- `GET /offers` - Get all offers
- `GET /offers/my` - Get user's offers
- `POST /offers` - Create new offer
- `POST /offers/:id/accept` - Accept offer
- `POST /offers/:id/cancel` - Cancel offer
- `GET /trades/my` - Get user's trades

## Styling

The application uses Tailwind CSS with custom configurations:

- Dark mode support (system preference)
- Custom rarity color classes for CS2 skins
- Responsive design for all screen sizes
- Consistent component styling

## Future Enhancements

- Steam OAuth login integration
- Real-time notifications for trade updates
- Advanced filtering and sorting
- User profiles and ratings
- Trade chat system
- Mobile app version
- Skin price estimates
- Trade history analytics

## Contributing

1. Follow the existing code style
2. Use TypeScript for all new code
3. Add proper error handling
4. Test all features before committing

## License

This project is part of the CS2 P2P Skins trading platform.
