import http from "k6/http";
import { check, group } from "k6";
import { API, authHeader, jsonHeaders } from "./config.js";

export const options = {
  scenarios: {
    e2e: { executor: "shared-iterations", vus: 1, iterations: 1 },
  },
  thresholds: {
    checks: ["rate==1.0"],
  },
};

export default function () {
  const unique = `user_${Date.now()}_${__VU}`;
  const creds = { username: unique, password: "testpass123", email: `${unique}@test.com` };
  let accessToken, refreshToken, taskId;

  // ── Auth: Register ──
  group("POST /auth/register", () => {
    const res = http.post(`${API}/auth/register`, JSON.stringify(creds), jsonHeaders);
    check(res, {
      "register 201": (r) => r.status === 201,
      "register success": (r) => r.json().success === true,
      "register has id": (r) => r.json().data.id > 0,
    });
  });

  // ── Auth: Register duplicate ──
  group("POST /auth/register duplicate", () => {
    const res = http.post(`${API}/auth/register`, JSON.stringify(creds), jsonHeaders);
    check(res, {
      "duplicate 400": (r) => r.status === 400,
    });
  });

  // ── Auth: Login ──
  group("POST /auth/login", () => {
    const res = http.post(
      `${API}/auth/login`,
      JSON.stringify({ username: creds.username, password: creds.password }),
      jsonHeaders,
    );
    check(res, {
      "login 200": (r) => r.status === 200,
      "login has access": (r) => !!r.json().data.access,
      "login has refresh": (r) => !!r.json().data.refresh,
    });
    accessToken = res.json().data.access;
    refreshToken = res.json().data.refresh;
  });

  // ── Auth: Login bad password ──
  group("POST /auth/login bad password", () => {
    const res = http.post(
      `${API}/auth/login`,
      JSON.stringify({ username: creds.username, password: "wrong" }),
      jsonHeaders,
    );
    check(res, { "bad login 401": (r) => r.status === 401 });
  });

  // ── Auth: Refresh ──
  group("POST /auth/refresh", () => {
    const res = http.post(
      `${API}/auth/refresh`,
      JSON.stringify({ refresh_token: refreshToken }),
      jsonHeaders,
    );
    check(res, {
      "refresh 200": (r) => r.status === 200,
      "refresh has access": (r) => !!r.json().data.access,
    });
    accessToken = res.json().data.access;
    refreshToken = res.json().data.refresh;
  });

  // ── User: Profile ──
  group("GET /user/profile", () => {
    const res = http.get(`${API}/user/profile`, authHeader(accessToken));
    check(res, {
      "profile 200": (r) => r.status === 200,
      "profile username": (r) => r.json().data.username === creds.username,
      "profile email": (r) => r.json().data.email === creds.email,
    });
  });

  // ── User: Profile unauthorized ──
  group("GET /user/profile no auth", () => {
    const res = http.get(`${API}/user/profile`, jsonHeaders);
    check(res, { "no auth 401": (r) => r.status === 401 });
  });

  // ── Task: Create ──
  group("POST /tasks", () => {
    const res = http.post(
      `${API}/tasks`,
      JSON.stringify({ name: "Test task", description: "k6 e2e test" }),
      authHeader(accessToken),
    );
    check(res, {
      "create task 201": (r) => r.status === 201,
      "create task name": (r) => r.json().data.name === "Test task",
      "create task status": (r) => r.json().data.status === "Created",
    });
    taskId = res.json().data.id;
  });

  // ── Task: Create validation ──
  group("POST /tasks invalid", () => {
    const res = http.post(`${API}/tasks`, JSON.stringify({}), authHeader(accessToken));
    check(res, { "invalid task 400": (r) => r.status === 400 });
  });

  // ── Task: List ──
  group("GET /tasks", () => {
    const res = http.get(`${API}/tasks`, authHeader(accessToken));
    check(res, {
      "list tasks 200": (r) => r.status === 200,
      "list has tasks": (r) => r.json().data.tasks.length > 0,
      "list has total": (r) => r.json().data.total > 0,
    });
  });

  // ── Task: List with query params ──
  group("GET /tasks?status=0&limit=5", () => {
    const res = http.get(`${API}/tasks?status=0&limit=5`, authHeader(accessToken));
    check(res, {
      "filtered list 200": (r) => r.status === 200,
    });
  });

  // ── Task: Get by ID ──
  group("GET /tasks/:id", () => {
    const res = http.get(`${API}/tasks/${taskId}`, authHeader(accessToken));
    check(res, {
      "get task 200": (r) => r.status === 200,
      "get task id match": (r) => r.json().data.id === taskId,
    });
  });

  // ── Task: Get not found ──
  group("GET /tasks/999999", () => {
    const res = http.get(`${API}/tasks/999999`, authHeader(accessToken));
    check(res, { "not found 404": (r) => r.status === 404 });
  });

  // ── Task: Update ──
  group("PUT /tasks/:id", () => {
    const res = http.put(
      `${API}/tasks/${taskId}`,
      JSON.stringify({ name: "Updated task", status: 1 }),
      authHeader(accessToken),
    );
    check(res, {
      "update task 200": (r) => r.status === 200,
      "update task name": (r) => r.json().data.name === "Updated task",
      "update task status": (r) => r.json().data.status === "Started",
    });
  });

  // ── Task: Archive ──
  group("PATCH /tasks/:id/archive", () => {
    const res = http.patch(`${API}/tasks/${taskId}/archive`, null, authHeader(accessToken));
    check(res, {
      "archive 200": (r) => r.status === 200,
      "archive status": (r) => r.json().data.status === "Canceled",
    });
  });

  // ── Task: Delete ──
  group("DELETE /tasks/:id", () => {
    const res = http.del(`${API}/tasks/${taskId}`, null, authHeader(accessToken));
    check(res, { "delete task 200": (r) => r.status === 200 });
  });

  // ── Task: Verify deleted ──
  group("GET /tasks/:id after delete", () => {
    const res = http.get(`${API}/tasks/${taskId}`, authHeader(accessToken));
    check(res, { "deleted 404": (r) => r.status === 404 });
  });

  // ── Auth: Logout ──
  group("POST /auth/logout", () => {
    const res = http.post(`${API}/auth/logout`, null, authHeader(accessToken));
    check(res, { "logout 200": (r) => r.status === 200 });
  });

  // ── Auth: Verify logged out ──
  group("GET /user/profile after logout", () => {
    const res = http.get(`${API}/user/profile`, authHeader(accessToken));
    check(res, { "revoked 401": (r) => r.status === 401 });
  });
}
