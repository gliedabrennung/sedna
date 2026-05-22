import { useEffect, useState } from 'react';
import { Sidebar } from './Sidebar';
import { ChatWindow } from './ChatWindow';
import { useAuthStore } from '../store/authStore';
import { Navigate } from 'react-router-dom';
import { api } from '../api';
import type { User } from '../types';

export function Layout() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const user = useAuthStore((s) => s.user);
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) {
      setIsChecking(false);
      return;
    }

    if (!user?.id) {
      useAuthStore.getState().logout();
      setIsChecking(false);
      return;
    }

    let cancelled = false;

    api
      .get<User[]>(`/users/bulk?ids=${user.id}`)
      .then((res) => {
        if (cancelled) return;
        if (res.data?.length > 0) {
          useAuthStore.getState().setAuth(res.data[0]);
        } else {
          useAuthStore.getState().logout();
        }
      })
      .catch(() => {
        if (!cancelled) useAuthStore.getState().logout();
      })
      .finally(() => {
        if (!cancelled) setIsChecking(false);
      });

    return () => {
      cancelled = true;
    };
  }, [isAuthenticated]);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (isChecking) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-zinc-950 text-zinc-100">
        <div className="flex flex-col items-center gap-3">
          <svg
            className="animate-spin h-8 w-8 text-indigo-500"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <span className="text-sm font-medium text-zinc-400">Verifying session...</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="flex h-screen overflow-hidden bg-zinc-950 text-zinc-100">
      <Sidebar />
      <ChatWindow />
    </div>
  );
}
