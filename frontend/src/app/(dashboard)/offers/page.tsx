'use client';

import { useEffect, useState } from 'react';
import { offersService } from '@/services/offers.service';
import { skinsService } from '@/services/skins.service';
import { OfferWithDetails, UserSkinWithDetails, OfferStatus } from '@/types';
import { Button } from '@/components/ui/Button';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { SkinCard } from '@/components/SkinCard';
import { handleApiError } from '@/lib/api';
import { Plus, X, Check, ArrowRight } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

export default function OffersPage() {
  const { user } = useAuth();
  const [allOffers, setAllOffers] = useState<OfferWithDetails[]>([]);
  const [myOffers, setMyOffers] = useState<OfferWithDetails[]>([]);
  const [mySkins, setMySkins] = useState<UserSkinWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showCreateOffer, setShowCreateOffer] = useState(false);
  const [selectedOfferedSkin, setSelectedOfferedSkin] = useState<number | null>(null);
  const [selectedRequestedSkin, setSelectedRequestedSkin] = useState<number | null>(null);
  const [viewMode, setViewMode] = useState<'all' | 'my'>('all');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setError('');
      const [allOffersData, myOffersData, mySkinsData] = await Promise.all([
        offersService.getAllOffers(),
        offersService.getMyOffers(),
        skinsService.getMySkins(),
      ]);
      setAllOffers(allOffersData);
      setMyOffers(myOffersData);
      setMySkins(mySkinsData);
    } catch (err) {
      setError(handleApiError(err));
    } finally {
      setLoading(false);
    }
  };

  const handleCreateOffer = async () => {
    if (!selectedOfferedSkin || !selectedRequestedSkin) {
      setError('Please select both skins for the trade');
      return;
    }

    try {
      setError('');
      await offersService.createOffer({
        skin_offered_id: selectedOfferedSkin,
        skin_requested_id: selectedRequestedSkin,
      });
      setShowCreateOffer(false);
      setSelectedOfferedSkin(null);
      setSelectedRequestedSkin(null);
      await loadData();
    } catch (err) {
      setError(handleApiError(err));
    }
  };

  const handleAcceptOffer = async (offerId: number) => {
    try {
      setError('');
      await offersService.acceptOffer(offerId);
      await loadData();
    } catch (err) {
      setError(handleApiError(err));
    }
  };

  const handleCancelOffer = async (offerId: number) => {
    try {
      setError('');
      await offersService.cancelOffer(offerId);
      await loadData();
    } catch (err) {
      setError(handleApiError(err));
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  const displayOffers = viewMode === 'all' ? allOffers : myOffers;
  const openOffers = displayOffers.filter((o) => o.status === OfferStatus.Open);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Trade Offers</h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Browse and create P2P trade offers
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant={viewMode === 'all' ? 'primary' : 'secondary'}
            onClick={() => setViewMode('all')}
          >
            All Offers
          </Button>
          <Button
            variant={viewMode === 'my' ? 'primary' : 'secondary'}
            onClick={() => setViewMode('my')}
          >
            My Offers
          </Button>
          <Button onClick={() => setShowCreateOffer(!showCreateOffer)}>
            <Plus className="w-4 h-4 mr-2" />
            Create Offer
          </Button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {showCreateOffer && (
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Create New Offer</h2>
          </CardHeader>
          <CardBody className="space-y-4">
            <div>
              <h3 className="font-medium mb-3">I want to offer:</h3>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {mySkins.map((userSkin) => (
                  <div
                    key={userSkin.id}
                    onClick={() => setSelectedOfferedSkin(userSkin.skin_id)}
                    className={`cursor-pointer transition-all ${
                      selectedOfferedSkin === userSkin.skin_id
                        ? 'ring-4 ring-blue-500 rounded-lg'
                        : ''
                    }`}
                  >
                    <SkinCard skin={userSkin.skin} quantity={userSkin.quantity} />
                  </div>
                ))}
              </div>
              {mySkins.length === 0 && (
                <p className="text-gray-500 text-sm">You don&apos;t have any skins to offer</p>
              )}
            </div>

            {selectedOfferedSkin && (
              <div>
                <h3 className="font-medium mb-3">I want to receive:</h3>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  {allOffers
                    .map((o) => o.skin_offered)
                    .filter(
                      (skin, index, self) =>
                        self.findIndex((s) => s.id === skin.id) === index
                    )
                    .map((skin) => (
                      <div
                        key={skin.id}
                        onClick={() => setSelectedRequestedSkin(skin.id)}
                        className={`cursor-pointer transition-all ${
                          selectedRequestedSkin === skin.id
                            ? 'ring-4 ring-green-500 rounded-lg'
                            : ''
                        }`}
                      >
                        <SkinCard skin={skin} />
                      </div>
                    ))}
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <Button
                onClick={handleCreateOffer}
                disabled={!selectedOfferedSkin || !selectedRequestedSkin}
              >
                Create Offer
              </Button>
              <Button variant="secondary" onClick={() => setShowCreateOffer(false)}>
                Cancel
              </Button>
            </div>
          </CardBody>
        </Card>
      )}

      <div className="space-y-4">
        {openOffers.length === 0 ? (
          <Card>
            <CardBody className="text-center py-12">
              <h3 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">
                No offers available
              </h3>
              <p className="text-gray-500">
                {viewMode === 'my'
                  ? 'Create your first offer to start trading'
                  : 'Be the first to create an offer'}
              </p>
            </CardBody>
          </Card>
        ) : (
          openOffers.map((offer) => {
            const isMyOffer = offer.owner_id === user?.id;
            const canAccept = !isMyOffer && mySkins.some((ms) => ms.skin_id === offer.skin_requested_id);

            return (
              <Card key={offer.id}>
                <CardBody>
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                        {isMyOffer ? 'Your offer' : `From: ${offer.owner_email}`}
                      </p>
                      <div className="flex items-center gap-4">
                        <div className="w-40">
                          <SkinCard skin={offer.skin_offered} />
                        </div>
                        <ArrowRight className="w-8 h-8 text-gray-400" />
                        <div className="w-40">
                          <SkinCard skin={offer.skin_requested} />
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2 ml-4">
                      {isMyOffer ? (
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => handleCancelOffer(offer.id)}
                        >
                          <X className="w-4 h-4 mr-1" />
                          Cancel
                        </Button>
                      ) : (
                        <Button
                          variant="success"
                          size="sm"
                          onClick={() => handleAcceptOffer(offer.id)}
                          disabled={!canAccept}
                        >
                          <Check className="w-4 h-4 mr-1" />
                          Accept
                        </Button>
                      )}
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
