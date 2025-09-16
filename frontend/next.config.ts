import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    reactCompiler: true,
  },
  /* config options here */
  output: "standalone",
};

export default nextConfig;
