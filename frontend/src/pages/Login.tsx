import type { FC } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { AuthForm } from '@/components/AuthForm';

export const Login: FC = () => {
  const { login, isLoading, error } = useAuth();

  return (
    <AuthForm
      title="Welcome Back"
      submitLabel="Sign In"
      altText="Don't have an account?"
      altLinkText="Sign up"
      altLinkTo="/register"
      isLoading={isLoading}
      error={error}
      onSubmit={login}
    />
  );
};
