/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    NEXY_PUBLIC_API_URL: process.env.NEXY_PUBLIC_API_URL,
  },
  images: {
    domains: ['localhost', 'workoscdn.com'],
  }
};

export default nextConfig;
