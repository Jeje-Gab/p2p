'use client';

import { useEffect, useState } from 'react';
import { tradesService } from '@/services/trades.service';
import { TradeWithDetails } from '@/types';
import { Card, CardBody } from '@/components/ui/Card';
import { SkinCard } from '@/components/SkinCard';
import { handleApiError } from '@/lib/api';
import { ArrowRight, Calendar } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

export default function TradesPage() {
  const { user } = useAuth();
  const [trades, setTrades] = useState<TradeWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    loadTrades();
  }, []);

  const loadTrades = async () => {
    try {
      setError('');
      const tradesData = await tradesService.getMyTrades();
      setTrades(tradesData);
    } catch (err) {
      setError(handleApiError(err));
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Trade History</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          View your completed trades
        </p>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {trades.length === 0 ? (
          <Card>
            <CardBody className="text-center py-12">
              <h3 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">
                No trades yet
              </h3>
              <p className="text-gray-500">
                Your completed trades will appear here
              </p>
            </CardBody>
          </Card>
        ) : (
          trades.map((trade) => {
            const iReceived = trade.to_user_id === user?.id;
            const skinReceived = iReceived
              ? trade.offer.skin_offered
              : trade.offer.skin_requested;
            const skinGiven = iReceived
              ? trade.offer.skin_requested
              : trade.offer.skin_offered;
            const tradedWith = iReceived ? trade.from_user_email : trade.to_user_email;

            return (
              <Card key={trade.id}>
                <CardBody>
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                      <Calendar className="w-4 h-4" />
                      {formatDate(trade.executed_at)}
                    </div>
                    <div className="text-sm text-gray-600 dark:text-gray-400">
                      Traded with: <span className="font-medium">{tradedWith}</span>
                    </div>
                  </div>

                  <div className="flex items-center gap-4">
                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        You gave:
                      </p>
                      <div className="w-48">
                        <SkinCard skin={skinGiven} />
                      </div>
                    </div>

                    <ArrowRight className="w-8 h-8 text-gray-400" />

                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        You received:
                      </p>
                      <div className="w-48">
                        <SkinCard skin={skinReceived} />
                      </div>
                    </div>
                  </div>
                </CardBody>
              </Card>
            );
          })
        )}
      </div>
    </div>
  );
}
