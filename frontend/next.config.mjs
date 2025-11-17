/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
      {
        protocol: 'https',
        hostname: 'community.cloudflare.steamstatic.com',
      },
      {
        protocol: 'https',
        hostname: 'steamcdn-a.akamaihd.net',
      },
      {
        protocol: 'https',
        hostname: 'steamcommunity-a.akamaihd.net',
      },
    ],
    // Disable image optimization for external images to avoid 404s
    unoptimized: process.env.NODE_ENV === 'development',
  },
};

export default nextConfig;
