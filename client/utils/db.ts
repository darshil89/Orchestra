// check for production or development environment

export const API_URL =
  process.env.NODE_ENV === "production"
    ? "https://api.orchestra.com"
    : "http://localhost:8080";
