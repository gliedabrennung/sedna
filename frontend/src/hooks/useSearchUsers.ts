import { useState, useEffect } from 'react';
import { api } from '@/api';
import type { User } from '@/types';

export function useSearchUsers(query: string) {
  const [results, setResults] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (query.trim().length < 3) {
      setResults([]);
      setError(null);
      return;
    }

    const controller = new AbortController();

    const delayDebounceFn = setTimeout(async () => {
      setIsLoading(true);
      setError(null);
      try {
        const res = await api.get<User[]>(`/users/search?q=${encodeURIComponent(query)}`, {
          signal: controller.signal,
        });
        setResults(res.data || []);
      } catch (err: any) {
        if (err.name !== 'CanceledError') {
          setError(err.response?.data?.message || 'Search failed');
          setResults([]);
        }
      } finally {
        setIsLoading(false);
      }
    }, 400);

    return () => {
      clearTimeout(delayDebounceFn);
      controller.abort();
    };
  }, [query]);

  return { results, isLoading, error };
}
