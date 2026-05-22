import type { ButtonHTMLAttributes, FC } from 'react';
import clsx from 'clsx';
import { Spinner } from './Spinner';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading?: boolean;
  variant?: 'primary' | 'ghost';
}

export const Button: FC<ButtonProps> = ({
  children,
  isLoading,
  disabled,
  variant = 'primary',
  className,
  ...props
}) => {
  return (
    <button
      disabled={disabled || isLoading}
      className={clsx(
        'w-full py-2.5 px-4 rounded-[var(--radius-md)] font-medium transition-all duration-200',
        'disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2',
        variant === 'primary' && 'gradient-accent text-white hover:opacity-90 shadow-sm hover:shadow-md',
        variant === 'ghost' && 'bg-transparent text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] hover:bg-[var(--color-surface-tertiary)]',
        className
      )}
      {...props}
    >
      {isLoading && <Spinner size="sm" className="text-white" />}
      {children}
    </button>
  );
};
