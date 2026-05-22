import { type ForwardRefRenderFunction, type InputHTMLAttributes, forwardRef } from 'react';
import clsx from 'clsx';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

const InputBase: ForwardRefRenderFunction<HTMLInputElement, InputProps> = (
  { label, error, className, id, ...props },
  ref
) => {
  const inputId = id || (label ? label.toLowerCase().replace(/\s+/g, '-') : undefined);

  return (
    <div className="w-full">
      {label && (
        <label
          htmlFor={inputId}
          className="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5"
        >
          {label}
        </label>
      )}
      <input
        ref={ref}
        id={inputId}
        className={clsx(
          'w-full px-4 py-2.5 bg-[var(--color-surface-primary)] border rounded-[var(--radius-md)] text-[var(--color-text-primary)] placeholder-[var(--color-text-muted)]',
          'focus:outline-none focus:ring-2 focus:ring-[var(--color-accent-start)]/30 focus:border-[var(--color-accent-start)] transition-all duration-200',
          error ? 'border-[var(--color-danger)]/50 focus:ring-[var(--color-danger)]/30' : 'border-[var(--color-border-primary)]',
          className
        )}
        {...props}
      />
      {error && (
        <span className="text-xs text-[var(--color-danger)] mt-1.5 block">{error}</span>
      )}
    </div>
  );
};

export const Input = forwardRef(InputBase);
