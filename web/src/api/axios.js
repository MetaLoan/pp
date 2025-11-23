import axios from 'axios';

const apiBase =
  import.meta.env?.VITE_API_BASE ||
  (typeof process !== 'undefined' && process.env?.VITE_API_BASE) ||
  'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: apiBase,
  timeout: 5000,
});

// Avoid multiple concurrent demo logins
let authPromise = null;

// Auto register/login a demo account when running locally to keep
// the trading and wallet endpoints usable without a manual login step.
const ensureDemoToken = async () => {
  const existing = localStorage.getItem('token');
  if (existing) return existing;

  if (authPromise) return authPromise;

  authPromise = (async () => {
    const email =
      import.meta.env?.VITE_DEMO_EMAIL ||
      (typeof process !== 'undefined' && process.env?.VITE_DEMO_EMAIL) ||
      'demo@example.com';
    const password =
      import.meta.env?.VITE_DEMO_PASS ||
      (typeof process !== 'undefined' && process.env?.VITE_DEMO_PASS) ||
      'demo1234';
    const nickname =
      import.meta.env?.VITE_DEMO_NICK ||
      (typeof process !== 'undefined' && process.env?.VITE_DEMO_NICK) ||
      'Demo';

    const client = axios.create({ baseURL: apiBase, timeout: 5000 });
    try {
      await client.post('/auth/register', { email, password, nickname });
    } catch (e) {
      // ignore if already exists
    }

    const loginResp = await client.post('/auth/login', { email, password });
    const token = loginResp.data.token;
    const user = loginResp.data.user;
    if (token) {
      localStorage.setItem('token', token);
      if (user) localStorage.setItem('user', JSON.stringify(user));
    }
    return token;
  })();

  try {
    const token = await authPromise;
    return token;
  } finally {
    authPromise = null;
  }
};

// Request interceptor to add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to transparently recover from 401 by auto-login
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const { response, config } = error || {};
    if (
      response?.status === 401 &&
      !config?._retry &&
      config?.url &&
      !config.url.includes('/auth/')
    ) {
      config._retry = true;
      try {
        const token = await ensureDemoToken();
        if (token) {
          config.headers = config.headers || {};
          config.headers.Authorization = `Bearer ${token}`;
          return api(config);
        }
      } catch (e) {
        // fall through to original error
      }
    }
    return Promise.reject(error);
  }
);

export default api;
