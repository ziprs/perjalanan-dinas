/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/spd',
  assetPrefix: '/spd',
  images: {
    unoptimized: true,
  },
  env: {
    API_URL: process.env.API_URL || 'http://localhost:8080/api',
  },
}

module.exports = nextConfig
