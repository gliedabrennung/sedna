import type { FC } from 'react';
import clsx from 'clsx';

interface AvatarProps {
  username?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const sizeClasses = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm font-semibold',
  lg: 'w-12 h-12 text-base font-semibold',
};

const gradients = [
  'from-indigo-500 to-purple-500',
  'from-emerald-500 to-teal-500',
  'from-rose-500 to-pink-500',
  'from-amber-500 to-orange-500',
  'from-cyan-500 to-blue-500',
  'from-fuchsia-500 to-violet-500',
];

function getGradient(name: string): string {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return gradients[Math.abs(hash) % gradients.length];
}

export const Avatar: FC<AvatarProps> = ({ username = '?', size = 'md', className }) => {
  const initial = username.charAt(0).toUpperCase();
  const gradient = getGradient(username);

  return (
    <div
      className={clsx(
        'rounded-full bg-gradient-to-br text-white flex items-center justify-center flex-shrink-0 select-none shadow-sm',
        gradient,
        sizeClasses[size],
        className
      )}
    >
      {initial}
    </div>
  );
};
