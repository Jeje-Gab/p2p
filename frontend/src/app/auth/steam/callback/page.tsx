'use client';

import { useEffect, useState, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { authService } from '@/services/auth.service';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';

export default function SteamCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setUser, setToken } = useAuth();
  const [error, setError] = useState<string>('');
  const [processing, setProcessing] = useState(true);
  const [requires2FA, setRequires2FA] = useState(false);
  const [userEmail, setUserEmail] = useState<string>('');
  const [code, setCode] = useState('');
  const [verifying, setVerifying] = useState(false);
  const hasProcessed = useRef(false);

  useEffect(() => {
    const handleCallback = async () => {
      // Prevent multiple executions
      if (hasProcessed.current) {
        return;
      }
      hasProcessed.current = true;

      try {
        console.log('Callback URL params:', Object.fromEntries(searchParams.entries()));

        // Check for error parameter
        const errorParam = searchParams.get('error');
        if (errorParam) {
          console.error('Steam auth error:', errorParam);
          setError('Steam authentication failed');
          setProcessing(false);
          return;
        }

        // Check if 2FA is required
        const requires2FAParam = searchParams.get('requires_2fa');
        const encodedUser = searchParams.get('user');

        if (requires2FAParam === 'true' && encodedUser) {
          // Decode user data
          const userJson = atob(encodedUser);
          const user = JSON.parse(userJson);

          setUserEmail(user.email);
          setRequires2FA(true);
          setProcessing(false);
          return;
        }

        // Extract token and user from URL
        const token = searchParams.get('token');

        console.log('Token:', token ? 'present' : 'missing');
        console.log('User:', encodedUser ? 'present' : 'missing');

        if (!token || !encodedUser) {
          console.error('Missing params - token:', !!token, 'user:', !!encodedUser);
          setError('Authentication failed - missing token or user data. Please try again.');
          setProcessing(false);
          return;
        }

        // Decode user data
        const userJson = atob(encodedUser);
        const user = JSON.parse(userJson);

        console.log('Decoded user:', user);

        // Store token and user info
        setToken(token);
        setUser(user);

        console.log('Redirecting to dashboard...');
        // Redirect to dashboard
        router.replace('/dashboard');
      } catch (err) {
        console.error('Callback error:', err);
        setError(err instanceof Error ? err.message : 'Steam authentication failed');
        setProcessing(false);
      }
    };

    handleCallback();
  }, [searchParams, router, setUser, setToken]);

  const handleVerify2FA = async () => {
    if (!code || code.length !== 6) {
      setError('Please enter a valid 6-digit code');
      return;
    }

    setVerifying(true);
    setError('');

    try {
      const response = await authService.verify2FAOAuth(userEmail, code);

      if (response.token && response.user) {
        setToken(response.token);
        setUser(response.user);
        router.replace('/dashboard');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid 2FA code');
    } finally {
      setVerifying(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900 px-4">
      <Card className="w-full max-w-md">
        {requires2FA ? (
          <>
            <CardHeader>
              <h1 className="text-2xl font-bold text-center">Two-Factor Authentication</h1>
              <p className="text-sm text-gray-600 dark:text-gray-400 text-center mt-2">
                Enter your 2FA code to complete Steam login
              </p>
            </CardHeader>
            <CardBody>
              <div className="space-y-4">
                <div className="w-full">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    2FA Code
                  </label>
                  <input
                    type="text"
                    placeholder="000000"
                    maxLength={6}
                    value={code}
                    onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                    className="w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-white border-gray-300 dark:border-gray-600"
                    autoFocus
                  />
                </div>

                {error && (
                  <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
                    {error}
                  </div>
                )}

                <Button
                  onClick={handleVerify2FA}
                  className="w-full"
                  disabled={verifying || code.length !== 6}
                >
                  {verifying ? 'Verifying...' : 'Verify'}
                </Button>

                <Button
                  onClick={() => router.push('/login')}
                  variant="secondary"
                  className="w-full"
                >
                  Cancel
                </Button>
              </div>
            </CardBody>
          </>
        ) : (
          <CardBody className="text-center py-8">
            {processing ? (
              <>
                <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-white mx-auto mb-4"></div>
                <h2 className="text-2xl font-bold mb-2">Processing Steam Login...</h2>
                <p className="text-gray-600 dark:text-gray-400">
                  Please wait while we verify your Steam account
                </p>
              </>
            ) : (
              <>
                <div className="text-red-600 text-6xl mb-4">✗</div>
                <h2 className="text-2xl font-bold mb-2">Authentication Failed</h2>
                <p className="text-gray-600 dark:text-gray-400 mb-4">{error}</p>
                <button
                  onClick={() => router.push('/login')}
                  className="text-blue-600 hover:underline"
                >
                  Return to Login
                </button>
              </>
            )}
          </CardBody>
        )}
      </Card>
    </div>
  );
}
