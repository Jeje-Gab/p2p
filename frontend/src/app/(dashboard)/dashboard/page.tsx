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
        <Card className="hover:shadow-xl transition-shadow duration-200">
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Total Skins</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.totalSkins}</p>
              <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">Items in inventory</p>
            </div>
            <div className="p-3 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl shadow-lg">
              <Package className="w-8 h-8 text-white" />
            </div>
          </CardBody>
        </Card>

        <Card className="hover:shadow-xl transition-shadow duration-200">
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Active Offers</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.activeOffers}</p>
              <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">Open trades</p>
            </div>
            <div className="p-3 bg-gradient-to-br from-green-500 to-green-600 rounded-xl shadow-lg">
              <RefreshCcw className="w-8 h-8 text-white" />
            </div>
          </CardBody>
        </Card>

        <Card className="hover:shadow-xl transition-shadow duration-200">
          <CardBody className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Completed Trades</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.completedTrades}</p>
              <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">Successful deals</p>
            </div>
            <div className="p-3 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl shadow-lg">
              <Trophy className="w-8 h-8 text-white" />
            </div>
          </CardBody>
        </Card>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Quick Actions</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Choose an action to get started
            </p>
          </CardHeader>
          <CardBody className="space-y-3">
            <Link href="/skins" className="block group">
              <div className="flex items-center justify-between p-4 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-lg shadow-md hover:shadow-lg transition-all duration-200 transform hover:-translate-y-0.5">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-white/20 rounded-lg group-hover:bg-white/30 transition-colors">
                    <Package className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="font-semibold">Manage My Skins</p>
                    <p className="text-xs text-blue-100">View and organize your inventory</p>
                  </div>
                </div>
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>

            <Link href="/offers" className="block group">
              <div className="flex items-center justify-between p-4 bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white rounded-lg shadow-md hover:shadow-lg transition-all duration-200 transform hover:-translate-y-0.5">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-white/20 rounded-lg group-hover:bg-white/30 transition-colors">
                    <TrendingUp className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="font-semibold">Browse Offers</p>
                    <p className="text-xs text-green-100">Discover new trade opportunities</p>
                  </div>
                </div>
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>

            <Link href="/trades" className="block group">
              <div className="flex items-center justify-between p-4 bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white rounded-lg shadow-md hover:shadow-lg transition-all duration-200 transform hover:-translate-y-0.5">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-white/20 rounded-lg group-hover:bg-white/30 transition-colors">
                    <Trophy className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="font-semibold">View Trade History</p>
                    <p className="text-xs text-purple-100">Check your completed trades</p>
                  </div>
                </div>
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Getting Started</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Follow these steps to start trading
            </p>
          </CardHeader>
          <CardBody>
            <div className="space-y-3">
              <div className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-semibold">
                  1
                </div>
                <div className="flex-1">
                  <p className="text-gray-700 dark:text-gray-300 font-medium">Add skins to your inventory</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Start by adding items you want to trade</p>
                </div>
              </div>

              <div className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-green-500 text-white flex items-center justify-center text-sm font-semibold">
                  2
                </div>
                <div className="flex-1">
                  <p className="text-gray-700 dark:text-gray-300 font-medium">Browse available offers</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Explore trades from other users</p>
                </div>
              </div>

              <div className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-500 text-white flex items-center justify-center text-sm font-semibold">
                  3
                </div>
                <div className="flex-1">
                  <p className="text-gray-700 dark:text-gray-300 font-medium">Create your own offers</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Set up trades with your items</p>
                </div>
              </div>

              <div className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-orange-500 text-white flex items-center justify-center text-sm font-semibold">
                  4
                </div>
                <div className="flex-1">
                  <p className="text-gray-700 dark:text-gray-300 font-medium">Accept offers that interest you</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Review and complete trades</p>
                </div>
              </div>

              <div className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-pink-500 text-white flex items-center justify-center text-sm font-semibold">
                  5
                </div>
                <div className="flex-1">
                  <p className="text-gray-700 dark:text-gray-300 font-medium">View your completed trades</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Track your trading history</p>
                </div>
              </div>
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
