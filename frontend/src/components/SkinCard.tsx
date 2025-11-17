import Image from 'next/image';
import { Skin } from '@/types';
import { Card, CardBody } from './ui/Card';

interface SkinCardProps {
  skin: Skin;
  quantity?: number;
  onClick?: () => void;
  actions?: React.ReactNode;
}

function getRarityClass(rarity: string): string {
  const rarityMap: Record<string, string> = {
    'Consumer Grade': 'rarity-consumer-grade',
    'Industrial Grade': 'rarity-industrial-grade',
    'Mil-Spec Grade': 'rarity-mil-spec-grade',
    'Restricted': 'rarity-restricted',
    'Classified': 'rarity-classified',
    'Covert': 'rarity-covert',
    'Exceedingly Rare': 'rarity-exceedingly-rare',
  };
  return rarityMap[rarity] || 'bg-gray-400 text-white';
}

export function SkinCard({ skin, quantity, onClick, actions }: SkinCardProps) {
  return (
    <Card className={`overflow-hidden ${onClick ? 'cursor-pointer hover:shadow-lg transition-shadow' : ''}`}>
      <div onClick={onClick}>
        <div className="relative h-48 bg-gradient-to-br from-gray-800 to-gray-900">
          <Image
            src={skin.image_url}
            alt={skin.name}
            fill
            className="object-contain p-4"
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            unoptimized
          />
          {quantity !== undefined && quantity > 1 && (
            <div className="absolute top-2 right-2 bg-black bg-opacity-75 text-white px-2 py-1 rounded-full text-sm font-bold">
              x{quantity}
            </div>
          )}
        </div>
        <CardBody className="space-y-2">
          <div className={`inline-block px-2 py-1 rounded text-xs font-semibold ${getRarityClass(skin.rarity)}`}>
            {skin.rarity}
          </div>
          <h3 className="font-semibold text-lg">{skin.name}</h3>
          <p className="text-sm text-gray-600 dark:text-gray-400">{skin.weapon_type}</p>
          {skin.float_value !== undefined && skin.float_value !== null && (
            <p className="text-xs text-gray-500">Float: {skin.float_value.toFixed(4)}</p>
          )}
        </CardBody>
      </div>
      {actions && <div className="px-6 pb-4">{actions}</div>}
    </Card>
  );
}
