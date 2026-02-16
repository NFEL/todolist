import http from "k6/http";
import { check, sleep } from "k6";
import { API, authHeader, jsonHeaders } from "./config.js";

export const options = {
  stages: [
    { duration: "10s", target: 10 },
    { duration: "20s", target: 10 },
    { duration: "5s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    checks: ["rate>0.95"],
  },
};

export function setup() {
  const unique = `load_${Date.now()}`;
  const creds = { username: unique, password: "loadtest123", email: `${unique}@test.com` };

  http.post(`${API}/auth/register`, JSON.stringify(creds), jsonHeaders);
  const login = http.post(
    `${API}/auth/login`,
    JSON.stringify({ username: creds.username, password: creds.password }),
    jsonHeaders,
  );
  return { token: login.json().data.access };
}

export default function (data) {
  const opts = authHeader(data.token);

  // Create a task
  const create = http.post(
    `${API}/tasks`,
    JSON.stringify({ name: `load-${__VU}-${__ITER}`, description: "load test" }),
    opts,
  );
  check(create, { "create 201": (r) => r.status === 201 });

  // List tasks
  const list = http.get(`${API}/tasks?limit=10`, opts);
  check(list, { "list 200": (r) => r.status === 200 });

  // Get profile
  const profile = http.get(`${API}/user/profile`, opts);
  check(profile, { "profile 200": (r) => r.status === 200 });

  // sleep(0.5);
}
