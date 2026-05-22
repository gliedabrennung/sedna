import type { FC } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { AuthForm } from '@/components/AuthForm';

export const Register: FC = () => {
  const { register, isLoading, error } = useAuth();

  return (
    <AuthForm
      title="Create Account"
      submitLabel="Sign Up"
      altText="Already have an account?"
      altLinkText="Sign in"
      altLinkTo="/login"
      isLoading={isLoading}
      error={error}
      onSubmit={register}
    />
  );
};
