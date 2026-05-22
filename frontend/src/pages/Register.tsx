import { useState, type FormEvent, type FC } from 'react';
import { useAuth } from '../hooks/useAuth';
import { Link } from 'react-router-dom';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';

export const Register: FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { register, isLoading, error } = useAuth();

  const handleRegister = (e: FormEvent) => {
    e.preventDefault();
    if (username.trim() && password) {
      register(username.trim(), password);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-950 px-4">
      <div className="w-full max-w-sm p-8 bg-zinc-900 border border-zinc-800 rounded-2xl shadow-xl">
        <h2 className="text-2xl font-semibold text-zinc-100 text-center mb-6">Create Account</h2>
        {error && (
          <div className="p-3 mb-4 text-sm text-red-400 bg-red-400/10 rounded-lg">
            {error}
          </div>
        )}

        <form onSubmit={handleRegister} className="space-y-4">
          <Input
            label="Username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Input
            label="Password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button type="submit" isLoading={isLoading}>
            Sign Up
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-zinc-500">
          Already have an account?{' '}
          <Link to="/login" className="text-indigo-400 hover:text-indigo-300 transition-colors">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
};
