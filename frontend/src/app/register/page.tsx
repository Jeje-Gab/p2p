'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { useAuth } from '@/contexts/AuthContext';
import { RegisterRequest } from '@/types';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { authService } from '@/services/auth.service';

interface RegisterFormData extends RegisterRequest {
  confirmPassword: string;
}

export default function RegisterPage() {
  const router = useRouter();
  const { register: registerUser } = useAuth();
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterFormData>();

  const password = watch('password');

  const handleSteamLogin = async () => {
    try {
      setLoading(true);
      const authUrl = await authService.getSteamLoginUrl();
      window.location.href = authUrl;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to initiate Steam login');
      setLoading(false);
    }
  };

  const onSubmit = async (data: RegisterFormData) => {
    setError('');
    setLoading(true);

    try {
      await registerUser({ email: data.email, password: data.password });
      setSuccess(true);
      setTimeout(() => {
        router.push('/login');
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900 px-4">
        <Card className="w-full max-w-md">
          <CardBody className="text-center py-8">
            <div className="text-green-600 text-6xl mb-4">✓</div>
            <h2 className="text-2xl font-bold mb-2">Registration Successful!</h2>
            <p className="text-gray-600 dark:text-gray-400">
              Redirecting to login page...
            </p>
          </CardBody>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900 px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <h1 className="text-2xl font-bold text-center">Create Account</h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 text-center mt-2">
            Join CS2 P2P Skins Trading
          </p>
        </CardHeader>
        <CardBody>
          <Button
            type="button"
            variant="secondary"
            className="w-full bg-[#171a21] hover:bg-[#1b2838] text-white border-none flex items-center justify-center gap-2"
            onClick={handleSteamLogin}
            disabled={loading}
          >
            <svg width="20" height="20" viewBox="0 0 256 256" fill="currentColor">
              <path d="M127.997 0C57.318 0 0 57.318 0 127.997c0 63.47 46.192 116.177 106.932 126.43l35.307-50.271c-24.243-3.574-42.82-24.58-42.82-49.94 0-27.915 22.618-50.533 50.533-50.533 27.915 0 50.533 22.618 50.533 50.533 0 25.36-18.577 46.366-42.82 49.94l35.307 50.271C254.302 244.174 256 236.256 256 227.997c0-70.679-57.318-127.997-127.997-127.997zm.005 33.666c-51.997 0-94.331 42.334-94.331 94.331 0 11.129 1.937 21.815 5.484 31.764l48.235-19.934c4.061-11.498 15.065-19.746 27.945-19.746 16.387 0 29.666 13.279 29.666 29.666s-13.279 29.666-29.666 29.666c-1.154 0-2.295-.067-3.418-.197l-22.347 32.013c5.484 1.154 11.129 1.771 16.932 1.771 51.997 0 94.331-42.334 94.331-94.331s-42.334-94.331-94.331-94.331z"/>
            </svg>
            Sign up with Steam
          </Button>

          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300 dark:border-gray-600"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-white dark:bg-gray-800 text-gray-500">OR</span>
            </div>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
              label="Email"
              type="email"
              placeholder="your@email.com"
              error={errors.email?.message}
              {...register('email', {
                required: 'Email is required',
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: 'Invalid email address',
                },
              })}
            />
            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              error={errors.password?.message}
              {...register('password', {
                required: 'Password is required',
                minLength: {
                  value: 6,
                  message: 'Password must be at least 6 characters',
                },
              })}
            />
            <Input
              label="Confirm Password"
              type="password"
              placeholder="••••••••"
              error={errors.confirmPassword?.message}
              {...register('confirmPassword', {
                required: 'Please confirm your password',
                validate: (value) => value === password || 'Passwords do not match',
              })}
            />

            {error && (
              <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
                {error}
              </div>
            )}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Creating Account...' : 'Create Account'}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Already have an account?{' '}
              <Link href="/login" className="text-blue-600 hover:underline">
                Sign in
              </Link>
            </p>
          </div>
        </CardBody>
      </Card>
    </div>
  );
}
