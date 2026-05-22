import type { FC } from 'react';
import clsx from 'clsx';
import type { ConnectionStatus } from '@/types';

interface StatusDotProps {
  status: ConnectionStatus;
  className?: string;
}

const statusColors: Record<ConnectionStatus, string> = {
  connected: 'bg-[var(--color-success)]',
  connecting: 'bg-amber-400 animate-pulse-soft',
  disconnected: 'bg-zinc-500',
};

const statusLabels: Record<ConnectionStatus, string> = {
  connected: 'Online',
  connecting: 'Connecting...',
  disconnected: 'Offline',
};

export const StatusDot: FC<StatusDotProps> = ({ status, className }) => (
  <div className={clsx('flex items-center gap-1.5', className)}>
    <span className={clsx('w-2 h-2 rounded-full', statusColors[status])} />
    <span className="text-xs text-[var(--color-text-muted)]">{statusLabels[status]}</span>
  </div>
);
