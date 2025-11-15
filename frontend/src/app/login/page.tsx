'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { useAuth } from '@/contexts/AuthContext';
import { LoginRequest } from '@/types';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';

interface LoginFormData extends LoginRequest {
  code?: string;
}

export default function LoginPage() {
  const router = useRouter();
  const { login, verify2FA } = useAuth();
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [requires2FA, setRequires2FA] = useState(false);
  const [credentials, setCredentials] = useState<LoginRequest | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>();

  const onSubmit = async (data: LoginFormData) => {
    setError('');
    setLoading(true);

    try {
      if (requires2FA && credentials) {
        if (!data.code) {
          setError('Please enter your 2FA code');
          setLoading(false);
          return;
        }
        await verify2FA({
          email: credentials.email,
          password: credentials.password,
          code: data.code,
        });
        router.push('/dashboard');
      } else {
        const result = await login({ email: data.email, password: data.password });
        if (result.requires2FA) {
          setRequires2FA(true);
          setCredentials({ email: data.email, password: data.password });
        } else {
          router.push('/dashboard');
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900 px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <h1 className="text-2xl font-bold text-center">CS2 P2P Skins Trading</h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 text-center mt-2">
            {requires2FA ? 'Enter your 2FA code' : 'Sign in to your account'}
          </p>
        </CardHeader>
        <CardBody>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {!requires2FA ? (
              <>
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
              </>
            ) : (
              <Input
                label="2FA Code"
                type="text"
                placeholder="000000"
                maxLength={6}
                error={errors.code?.message}
                {...register('code', {
                  required: '2FA code is required',
                  pattern: {
                    value: /^\d{6}$/,
                    message: 'Code must be 6 digits',
                  },
                })}
              />
            )}

            {error && (
              <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
                {error}
              </div>
            )}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Loading...' : requires2FA ? 'Verify' : 'Sign In'}
            </Button>

            {requires2FA && (
              <Button
                type="button"
                variant="secondary"
                className="w-full"
                onClick={() => {
                  setRequires2FA(false);
                  setCredentials(null);
                }}
              >
                Back to Login
              </Button>
            )}
          </form>

          {!requires2FA && (
            <>
              <div className="mt-6 text-center">
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Don&apos;t have an account?{' '}
                  <Link href="/register" className="text-blue-600 hover:underline">
                    Sign up
                  </Link>
                </p>
              </div>
            </>
          )}
        </CardBody>
      </Card>
    </div>
  );
}
