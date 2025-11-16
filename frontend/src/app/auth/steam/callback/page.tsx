'use client';

import { useEffect, useState, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { authService } from '@/services/auth.service';
import { Card, CardBody } from '@/components/ui/Card';

export default function SteamCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setUser, setToken } = useAuth();
  const [error, setError] = useState<string>('');
  const [processing, setProcessing] = useState(true);
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

        // Extract token and user from URL
        const token = searchParams.get('token');
        const encodedUser = searchParams.get('user');

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

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900 px-4">
      <Card className="w-full max-w-md">
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
      </Card>
    </div>
  );
}
