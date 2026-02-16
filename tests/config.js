export const BASE_URL = __ENV.BASE_URL || "http://localhost:3154";
export const API = `${BASE_URL}/v1`;

export function authHeader(token) {
  return { headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` } };
}

export const jsonHeaders = { headers: { "Content-Type": "application/json" } };
