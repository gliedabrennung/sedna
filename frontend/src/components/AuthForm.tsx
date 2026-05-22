import { useState, type SyntheticEvent, type FC } from 'react';
import { Link } from 'react-router-dom';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';

interface AuthFormProps {
  title: string;
  submitLabel: string;
  altText: string;
  altLinkText: string;
  altLinkTo: string;
  isLoading: boolean;
  error: string | null;
  onSubmit: (username: string, password: string) => void;
}

export const AuthForm: FC<AuthFormProps> = ({
  title,
  submitLabel,
  altText,
  altLinkText,
  altLinkTo,
  isLoading,
  error,
  onSubmit,
}) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    if (username.trim() && password) {
      onSubmit(username.trim(), password);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--color-surface-primary)] px-4">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-96 h-96 rounded-full bg-[var(--color-accent-start)]/5 blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-96 h-96 rounded-full bg-[var(--color-accent-end)]/5 blur-3xl" />
      </div>

      <div className="w-full max-w-sm p-8 bg-[var(--color-surface-secondary)] border border-[var(--color-border-primary)] rounded-[var(--radius-xl)] shadow-[var(--shadow-card)] animate-scale-in relative">
        <div className="flex justify-center mb-6">
          <div className="w-12 h-12 rounded-[var(--radius-md)] gradient-accent flex items-center justify-center shadow-lg">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
            </svg>
          </div>
        </div>

        <h1 className="text-2xl font-bold text-[var(--color-text-primary)] text-center mb-1">{title}</h1>
        <p className="text-sm text-[var(--color-text-muted)] text-center mb-6">
          {title === 'Welcome Back' ? 'Sign in to continue messaging' : 'Create your account to get started'}
        </p>

        {error && (
          <div className="p-3 mb-4 text-sm text-[var(--color-danger)] bg-[var(--color-danger-bg)] rounded-[var(--radius-md)] animate-slide-up">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            id="auth-username"
            label="Username"
            type="text"
            required
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Input
            id="auth-password"
            label="Password"
            type="password"
            required
            autoComplete={title === 'Welcome Back' ? 'current-password' : 'new-password'}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button id="auth-submit" type="submit" isLoading={isLoading}>
            {submitLabel}
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-[var(--color-text-muted)]">
          {altText}{' '}
          <Link to={altLinkTo} className="text-[var(--color-accent-start)] hover:text-[var(--color-accent-hover)] transition-colors font-medium">
            {altLinkText}
          </Link>
        </p>
      </div>
    </div>
  );
};
