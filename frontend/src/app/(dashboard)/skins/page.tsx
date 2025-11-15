'use client';

import { useEffect, useState } from 'react';
import { skinsService } from '@/services/skins.service';
import { UserSkinWithDetails, Skin } from '@/types';
import { SkinCard } from '@/components/SkinCard';
import { Button } from '@/components/ui/Button';
import { Card, CardBody, CardHeader } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { handleApiError } from '@/lib/api';
import { Plus, Minus, Search } from 'lucide-react';

export default function SkinsPage() {
  const [mySkins, setMySkins] = useState<UserSkinWithDetails[]>([]);
  const [allSkins, setAllSkins] = useState<Skin[]>([]);
  const [filteredSkins, setFilteredSkins] = useState<Skin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [searchTerm, setSearchTerm] = useState('');
  const [showAddSkin, setShowAddSkin] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  useEffect(() => {
    if (searchTerm) {
      const filtered = allSkins.filter(
        (skin) =>
          skin.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          skin.weapon_type.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setFilteredSkins(filtered);
    } else {
      setFilteredSkins(allSkins);
    }
  }, [searchTerm, allSkins]);

  const loadData = async () => {
    try {
      setError('');
      const [mySkinsData, allSkinsData] = await Promise.all([
        skinsService.getMySkins(),
        skinsService.getAllSkins(),
      ]);
      setMySkins(mySkinsData);
      setAllSkins(allSkinsData);
      setFilteredSkins(allSkinsData);
    } catch (err) {
      setError(handleApiError(err));
    } finally {
      setLoading(false);
    }
  };

  const handleAddSkin = async (skinId: number) => {
    try {
      setError('');
      await skinsService.addSkinToInventory(skinId, 1);
      await loadData();
    } catch (err) {
      setError(handleApiError(err));
    }
  };

  const handleRemoveSkin = async (skinId: number) => {
    try {
      setError('');
      await skinsService.removeSkinFromInventory(skinId, 1);
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

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">My Skins</h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Manage your CS2 skins inventory
          </p>
        </div>
        <Button onClick={() => setShowAddSkin(!showAddSkin)}>
          {showAddSkin ? 'View My Skins' : 'Add Skins'}
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {showAddSkin ? (
        <div>
          <Card className="mb-6">
            <CardBody>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                <Input
                  type="text"
                  placeholder="Search skins..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </CardBody>
          </Card>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredSkins.map((skin) => {
              const ownedSkin = mySkins.find((ms) => ms.skin_id === skin.id);
              return (
                <SkinCard
                  key={skin.id}
                  skin={skin}
                  quantity={ownedSkin?.quantity}
                  actions={
                    <Button
                      variant="primary"
                      size="sm"
                      className="w-full"
                      onClick={() => handleAddSkin(skin.id)}
                    >
                      <Plus className="w-4 h-4 mr-1" />
                      Add to Inventory
                    </Button>
                  }
                />
              );
            })}
          </div>

          {filteredSkins.length === 0 && (
            <div className="text-center py-12 text-gray-500">
              No skins found matching your search.
            </div>
          )}
        </div>
      ) : (
        <>
          {mySkins.length === 0 ? (
            <Card>
              <CardBody className="text-center py-12">
                <Package className="w-16 h-16 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">
                  No skins in your inventory
                </h3>
                <p className="text-gray-500 mb-4">
                  Start by adding some skins to your collection
                </p>
                <Button onClick={() => setShowAddSkin(true)}>
                  <Plus className="w-4 h-4 mr-2" />
                  Add Skins
                </Button>
              </CardBody>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {mySkins.map((userSkin) => (
                <SkinCard
                  key={userSkin.id}
                  skin={userSkin.skin}
                  quantity={userSkin.quantity}
                  actions={
                    <Button
                      variant="danger"
                      size="sm"
                      className="w-full"
                      onClick={() => handleRemoveSkin(userSkin.skin_id)}
                    >
                      <Minus className="w-4 h-4 mr-1" />
                      Remove
                    </Button>
                  }
                />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}

// Import Package icon
import { Package } from 'lucide-react';
