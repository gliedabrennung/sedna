import { useState, useCallback } from 'react';
import { api } from '@/api';
import { useAuthStore } from '@/store/authStore';
import { useNavigate } from 'react-router-dom';

export function useAuth() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const setAuth = useAuthStore((state) => state.setAuth);
  const navigate = useNavigate();

  const login = useCallback(
    async (username: string, password: string) => {
      setIsLoading(true);
      setError(null);
      try {
        const res = await api.post('/auth/login', { username, password });
        setAuth(res.data.user);
        navigate('/');
      } catch (err: any) {
        setError(err.response?.data?.message || 'Login failed');
      } finally {
        setIsLoading(false);
      }
    },
    [setAuth, navigate]
  );

  const register = useCallback(
    async (username: string, password: string) => {
      setIsLoading(true);
      setError(null);
      try {
        await api.post('/auth/register', { username, password });
        navigate('/login');
      } catch (err: any) {
        setError(err.response?.data?.message || 'Registration failed');
      } finally {
        setIsLoading(false);
      }
    },
    [navigate]
  );

  return { login, register, isLoading, error };
}
