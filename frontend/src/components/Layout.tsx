import { useEffect, useState, useRef, type FC } from 'react';
import { Sidebar } from '@/components/Sidebar';
import { ChatWindow } from '@/components/ChatWindow';
import { useAuthStore } from '@/store/authStore';
import { Navigate } from 'react-router-dom';
import { api } from '@/api';
import { Spinner } from '@/components/ui/Spinner';
import type { User } from '@/types';

export const Layout: FC = () => {
  const user = useAuthStore((s) => s.user);
  const [isChecking, setIsChecking] = useState(true);
  const hasVerified = useRef(false);

  useEffect(() => {
    if (hasVerified.current) return;

    if (!user?.id) {
      useAuthStore.getState().logout();
      setIsChecking(false);
      return;
    }

    hasVerified.current = true;
    let cancelled = false;

    api
      .get<User>(`/users/me`)
      .then((res) => {
        if (cancelled) return;
        if (res.data && res.data.id) {
          useAuthStore.getState().setAuth(res.data);
        }
      })
      .catch(() => {
        /* session check failed — keep current user from localStorage, don't force logout */
      })
      .finally(() => {
        if (!cancelled) setIsChecking(false);
      });

    return () => {
      cancelled = true;
    };
  }, [user?.id]);

  if (isChecking) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-[var(--color-surface-primary)]">
        <div className="flex flex-col items-center gap-3 animate-fade-in">
          <Spinner size="lg" />
          <span className="text-sm font-medium text-[var(--color-text-muted)]">Verifying session...</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="flex h-screen overflow-hidden bg-[var(--color-surface-primary)] text-[var(--color-text-primary)]">
      <Sidebar />
      <ChatWindow />
    </div>
  );
};
