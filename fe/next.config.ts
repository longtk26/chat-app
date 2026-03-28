import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    rewrites: async () => [
        {
            source: "/api/:path*",
            destination: "http://localhost:8080/api/:path*",
        },
        {
            source: "/socket.io/:path*",
            destination: "http://localhost:8080/socket.io/:path*",
        },
    ],
};

export default nextConfig;
