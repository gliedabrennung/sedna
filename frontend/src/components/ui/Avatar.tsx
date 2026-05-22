import type { FC } from 'react';
import clsx from 'clsx';

interface AvatarProps {
  username?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const Avatar: FC<AvatarProps> = ({ username = '?', size = 'md', className }) => {
  const initial = username.charAt(0).toUpperCase();

  const sizes = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm font-semibold',
    lg: 'w-12 h-12 text-base font-semibold',
  };

  return (
    <div
      className={clsx(
        'rounded-full bg-indigo-500/20 text-indigo-400 flex items-center justify-center flex-shrink-0 select-none',
        sizes[size],
        className
      )}
    >
      {initial}
    </div>
  );
};
