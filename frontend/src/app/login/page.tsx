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
import { authService } from '@/services/auth.service';

interface LoginFormData extends LoginRequest {
  code?: string;
}

export default function LoginPage() {
  const router = useRouter();
  const { login, verify2FA } = useAuth();
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [requires2FA, setRequires2FA] = useState(false);
  const [useBackupCode, setUseBackupCode] = useState(false);
  const [credentials, setCredentials] = useState<LoginRequest | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>();

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

  const onSubmit = async (data: LoginFormData) => {
    setError('');
    setLoading(true);

    try {
      if (requires2FA && credentials) {
        if (!data.code) {
          setError(useBackupCode ? 'Please enter your backup code' : 'Please enter your 2FA code');
          setLoading(false);
          return;
        }

        // Use backup code or regular 2FA
        if (useBackupCode) {
          const response = await authService.verifyBackupCode({
            email: credentials.email,
            password: credentials.password,
            code: data.code,
          });
          if (response.token && response.user) {
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));
            router.push('/dashboard');
          }
        } else {
          await verify2FA({
            email: credentials.email,
            password: credentials.password,
            code: data.code,
          });
          router.push('/dashboard');
        }
      } else {
        const result = await login({ email: data.email, password: data.password });
        if (result.requires2FA) {
          setRequires2FA(true);
          setCredentials({ email: data.email, password: data.password });
        } else if (result.success) {
          // Use replace instead of push to prevent back button issues
          router.replace('/dashboard');
        } else {
          setError('Login failed - please try again');
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
          {!requires2FA && (
            <>
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
                Login with Steam
              </Button>

              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300 dark:border-gray-600"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white dark:bg-gray-800 text-gray-500">OR</span>
                </div>
              </div>
            </>
          )}

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
              <>
                <Input
                  label={useBackupCode ? 'Backup Code' : '2FA Code'}
                  type="text"
                  placeholder={useBackupCode ? 'XXXX-XXXX' : '000000'}
                  maxLength={useBackupCode ? 9 : 6}
                  error={errors.code?.message}
                  {...register('code', {
                    required: useBackupCode ? 'Backup code is required' : '2FA code is required',
                    pattern: useBackupCode ? {
                      value: /^[A-Z0-9]{4}-[A-Z0-9]{4}$/i,
                      message: 'Backup code format: XXXX-XXXX',
                    } : {
                      value: /^\d{6}$/,
                      message: 'Code must be 6 digits',
                    },
                  })}
                />
                <button
                  type="button"
                  className="text-sm text-blue-600 hover:underline -mt-2"
                  onClick={() => setUseBackupCode(!useBackupCode)}
                >
                  {useBackupCode ? 'Use authenticator code instead' : 'Use backup code instead'}
                </button>
              </>
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
