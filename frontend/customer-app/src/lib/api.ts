import axios, { AxiosError } from 'axios';
import type { LoginRequest, RegisterRequest, AuthResponse, RefreshTokenResponse } from '@/types/auth';
import type { Book, BookFilters, BooksResponse, Category } from '@/types/book';
import type { User } from '@/types/user';
import type { WishlistItem, WishlistResponse } from '@/types/wishlist';

const api = axios.create({
  baseURL: 'http://localhost',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
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

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && originalRequest) {
      const refreshToken = localStorage.getItem('refresh_token');

      if (refreshToken) {
        try {
          const { data } = await axios.post<RefreshTokenResponse>(
            'http://localhost:8082/api/v1/auth/refresh',
            { refresh_token: refreshToken }
          );

          localStorage.setItem('token', data.token);
          localStorage.setItem('refresh_token', data.refresh_token);

          originalRequest.headers.Authorization = `Bearer ${data.token}`;
          return api(originalRequest);
        } catch (refreshError) {
          localStorage.removeItem('token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
        }
      }
    }

    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (data: LoginRequest) =>
    api.post<AuthResponse>('/api/v1/auth/login', data),

  register: (data: RegisterRequest) =>
    api.post<AuthResponse>('/api/v1/auth/register', data),

  logout: () =>
    api.post('/api/v1/auth/logout'),

  me: () =>
    api.get<{ data: User }>('/api/v1/users/me'),
};

// Books API
export const booksAPI = {
  list: (params?: BookFilters) =>
    api.get<BooksResponse>('/api/v1/books', { params }),

  get: (id: string) =>
    api.get<{ data: Book }>(`/api/v1/books/${id}`),

  create: (book: Partial<Book>) =>
    api.post<{ data: Book }>('/api/v1/books', book),

  update: (id: string, book: Partial<Book>) =>
    api.put<{ data: Book }>(`/api/v1/books/${id}`, book),

  delete: (id: string) =>
    api.delete(`/api/v1/books/${id}`),
};

// Categories API
export const categoriesAPI = {
  list: () =>
    api.get<{ data: Category[] }>('/api/v1/categories'),

  get: (id: string) =>
    api.get<{ data: Category }>(`/api/v1/categories/${id}`),
};

// Wishlist API
export const wishlistAPI = {
  list: () =>
    api.get<WishlistResponse>('/api/v1/users/me/wishlist'),

  add: (book_id: string) =>
    api.post<{ data: WishlistItem }>('/api/v1/users/me/wishlist', { book_id }),

  remove: (book_id: string) =>
    api.delete(`/api/v1/users/me/wishlist/${book_id}`),
};

export default api;
