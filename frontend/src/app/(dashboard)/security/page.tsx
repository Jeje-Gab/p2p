'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { authService } from '@/services/auth.service';
import { Button } from '@/components/ui/Button';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import QRCode from 'qrcode';

export default function SecurityPage() {
  const { user, refreshUser } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // 2FA Setup states
  const [showSetup, setShowSetup] = useState(false);
  const [qrCodeUrl, setQrCodeUrl] = useState('');
  const [secret, setSecret] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [showBackupCodes, setShowBackupCodes] = useState(false);

  // Disable 2FA state
  const [showDisable, setShowDisable] = useState(false);
  const [disableCode, setDisableCode] = useState('');

  const handleSetup2FA = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await authService.setup2FA();
      setSecret(response.secret);

      // Generate QR code image
      const qrCodeDataUrl = await QRCode.toDataURL(response.qr_url);
      setQrCodeUrl(qrCodeDataUrl);

      setShowSetup(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to setup 2FA');
    } finally {
      setLoading(false);
    }
  };

  const handleEnable2FA = async () => {
    if (!verificationCode || verificationCode.length !== 6) {
      setError('Please enter a valid 6-digit code');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const response = await authService.enable2FA(verificationCode);
      setBackupCodes(response.backup_codes);
      setShowBackupCodes(true);
      setShowSetup(false);
      setSuccess('2FA enabled successfully!');
      await refreshUser();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to enable 2FA');
    } finally {
      setLoading(false);
    }
  };

  const handleDisable2FA = async () => {
    if (!disableCode || disableCode.length !== 6) {
      setError('Please enter a valid 6-digit code');
      return;
    }

    setLoading(true);
    setError('');
    try {
      await authService.disable2FA(disableCode);
      setShowDisable(false);
      setDisableCode('');
      setSuccess('2FA disabled successfully');
      await refreshUser();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to disable 2FA');
    } finally {
      setLoading(false);
    }
  };

  const copyBackupCodes = () => {
    const text = backupCodes.join('\n');
    navigator.clipboard.writeText(text);
    setSuccess('Backup codes copied to clipboard!');
  };

  const downloadBackupCodes = () => {
    const text = `CS2 P2P Skins Trading - 2FA Backup Codes\n\nThese codes can only be used once. Save them in a secure place.\n\n${backupCodes.join('\n')}`;
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = '2fa-backup-codes.txt';
    a.click();
    URL.revokeObjectURL(url);
  };

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold mb-8">Security Settings</h1>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-green-600 dark:text-green-400 px-4 py-3 rounded mb-4">
          {success}
        </div>
      )}

      {/* Backup Codes Display */}
      {showBackupCodes && backupCodes.length > 0 && (
        <Card className="mb-6 border-yellow-500">
          <CardHeader>
            <h2 className="text-xl font-bold text-yellow-600">⚠️ Save Your Backup Codes</h2>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
              These codes will only be shown once. Save them in a secure place. Each code can only be used once.
            </p>
          </CardHeader>
          <CardBody>
            <div className="bg-gray-100 dark:bg-gray-800 p-4 rounded font-mono text-sm grid grid-cols-2 gap-2 mb-4">
              {backupCodes.map((code, index) => (
                <div key={index} className="py-1">
                  {index + 1}. {code}
                </div>
              ))}
            </div>
            <div className="flex gap-2">
              <Button onClick={copyBackupCodes} variant="secondary">
                📋 Copy All
              </Button>
              <Button onClick={downloadBackupCodes} variant="secondary">
                💾 Download
              </Button>
              <Button onClick={() => setShowBackupCodes(false)} className="ml-auto">
                I&apos;ve Saved Them
              </Button>
            </div>
          </CardBody>
        </Card>
      )}

      {/* Two-Factor Authentication Card */}
      <Card>
        <CardHeader>
          <h2 className="text-2xl font-bold">Two-Factor Authentication (2FA)</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
            Add an extra layer of security to your account using Google Authenticator or similar apps.
          </p>
        </CardHeader>
        <CardBody>
          <div className="flex items-center justify-between mb-4">
            <div>
              <p className="font-semibold">Status:</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {user.twofa_enabled ? (
                  <span className="text-green-600 font-semibold">✓ Enabled</span>
                ) : (
                  <span className="text-gray-500">Disabled</span>
                )}
              </p>
            </div>
            {!user.twofa_enabled && !showSetup && (
              <Button onClick={handleSetup2FA} disabled={loading}>
                {loading ? 'Loading...' : 'Enable 2FA'}
              </Button>
            )}
            {user.twofa_enabled && !showDisable && (
              <Button onClick={() => setShowDisable(true)} variant="secondary">
                Disable 2FA
              </Button>
            )}
          </div>

          {/* Setup 2FA Form */}
          {showSetup && (
            <div className="border-t pt-4 space-y-4">
              <div>
                <h3 className="font-semibold mb-2">Step 1: Scan QR Code</h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                  Open Google Authenticator and scan this QR code:
                </p>
                {qrCodeUrl && (
                  <div className="flex justify-center bg-white p-4 rounded">
                    <img src={qrCodeUrl} alt="2FA QR Code" className="w-64 h-64" />
                  </div>
                )}
              </div>

              <div>
                <h3 className="font-semibold mb-2">Step 2: Enter Verification Code</h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                  Enter the 6-digit code from your authenticator app:
                </p>
                <div className="w-full">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Verification Code
                  </label>
                  <input
                    type="text"
                    placeholder="000000"
                    maxLength={6}
                    value={verificationCode}
                    onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, ''))}
                    className="w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-white border-gray-300 dark:border-gray-600"
                  />
                </div>
              </div>

              <div className="flex gap-2">
                <Button onClick={handleEnable2FA} disabled={loading || verificationCode.length !== 6}>
                  {loading ? 'Verifying...' : 'Enable 2FA'}
                </Button>
                <Button
                  onClick={() => {
                    setShowSetup(false);
                    setVerificationCode('');
                    setQrCodeUrl('');
                    setSecret('');
                  }}
                  variant="secondary"
                >
                  Cancel
                </Button>
              </div>
            </div>
          )}

          {/* Disable 2FA Form */}
          {showDisable && (
            <div className="border-t pt-4 space-y-4">
              <div>
                <h3 className="font-semibold mb-2">Disable Two-Factor Authentication</h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                  Enter your current 2FA code to disable two-factor authentication:
                </p>
                <div className="w-full">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    2FA Code
                  </label>
                  <input
                    type="text"
                    placeholder="000000"
                    maxLength={6}
                    value={disableCode}
                    onChange={(e) => setDisableCode(e.target.value.replace(/\D/g, ''))}
                    className="w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-white border-gray-300 dark:border-gray-600"
                  />
                </div>
              </div>

              <div className="flex gap-2">
                <Button onClick={handleDisable2FA} disabled={loading || disableCode.length !== 6} variant="secondary">
                  {loading ? 'Disabling...' : 'Disable 2FA'}
                </Button>
                <Button onClick={() => setShowDisable(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          )}
        </CardBody>
      </Card>

      {/* Password Change Card (placeholder) */}
      <Card className="mt-6">
        <CardHeader>
          <h2 className="text-2xl font-bold">Password</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
            Change your account password
          </p>
        </CardHeader>
        <CardBody>
          <p className="text-sm text-gray-500">Password change feature coming soon...</p>
        </CardBody>
      </Card>
    </div>
  );
}
