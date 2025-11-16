'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { skinsService } from '@/services/skins.service';
import { offersService } from '@/services/offers.service';
import { tradesService } from '@/services/trades.service';
import { UserSkinWithDetails, OfferWithDetails, TradeWithDetails } from '@/types';
import Link from 'next/link';
import { Package, RefreshCcw, Trophy, TrendingUp } from 'lucide-react';

export default function DashboardPage() {
  const { user } = useAuth();
  const [stats, setStats] = useState({
    totalSkins: 0,
    activeOffers: 0,
    completedTrades: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const [mySkins, myOffers, myTrades] = await Promise.all([
        skinsService.getMySkins(),
        offersService.getMyOffers(),
        tradesService.getMyTrades(),
      ]);

      setStats({
        totalSkins: (mySkins || []).reduce((acc, item) => acc + item.quantity, 0),
        activeOffers: (myOffers || []).filter((o) => o.status === 'open').length,
        completedTrades: (myTrades || []).length,
      });
    } catch (error) {
      console.error('Failed to load stats:', error);
      setStats({
        totalSkins: 0,
        activeOffers: 0,
        completedTrades: 0,
      });
    } finally {
      setLoading(false);
    }
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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Welcome back, {user?.email}!
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Skins</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.totalSkins}</p>
            </div>
            <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-full">
              <Package className="w-8 h-8 text-blue-600 dark:text-blue-300" />
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Active Offers</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.activeOffers}</p>
            </div>
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-full">
              <RefreshCcw className="w-8 h-8 text-green-600 dark:text-green-300" />
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Completed Trades</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.completedTrades}</p>
            </div>
            <div className="p-3 bg-purple-100 dark:bg-purple-900 rounded-full">
              <Trophy className="w-8 h-8 text-purple-600 dark:text-purple-300" />
            </div>
          </CardBody>
        </Card>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Quick Actions</h2>
          </CardHeader>
          <CardBody className="space-y-4">
            <Link href="/skins">
              <Button variant="primary" className="w-full">
                <Package className="w-4 h-4 mr-2" />
                Manage My Skins
              </Button>
            </Link>
            <Link href="/offers">
              <Button variant="secondary" className="w-full">
                <TrendingUp className="w-4 h-4 mr-2" />
                Browse Offers
              </Button>
            </Link>
            <Link href="/trades">
              <Button variant="secondary" className="w-full">
                <Trophy className="w-4 h-4 mr-2" />
                View Trade History
              </Button>
            </Link>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Getting Started</h2>
          </CardHeader>
          <CardBody>
            <ol className="list-decimal list-inside space-y-2 text-gray-700 dark:text-gray-300">
              <li>Add skins to your inventory</li>
              <li>Browse available offers from other users</li>
              <li>Create your own trade offers</li>
              <li>Accept offers that interest you</li>
              <li>View your completed trades</li>
            </ol>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
