import { type ForwardRefRenderFunction, type InputHTMLAttributes, forwardRef } from 'react';
import clsx from 'clsx';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

const InputBase: ForwardRefRenderFunction<HTMLInputElement, InputProps> = (
  { label, error, className, ...props },
  ref
) => {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-zinc-400 mb-1">
          {label}
        </label>
      )}
      <input
        ref={ref}
        className={clsx(
          'w-full px-4 py-2 bg-zinc-950 border rounded-lg text-zinc-100 placeholder-zinc-600',
          'focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 transition-colors',
          error ? 'border-red-500/50 focus:ring-red-500' : 'border-zinc-800',
          className
        )}
        {...props}
      />
      {error && (
        <span className="text-xs text-red-400 mt-1 block">
          {error}
        </span>
      )}
    </div>
  );
};

export const Input = forwardRef(InputBase);
