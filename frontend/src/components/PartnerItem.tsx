import { memo, type FC } from 'react';
import clsx from 'clsx';
import { Avatar } from './ui/Avatar';
import type { User } from '../types';

interface PartnerItemProps {
  partner: User;
  isActive: boolean;
  onClick: () => void;
}

const PartnerItemBase: FC<PartnerItemProps> = ({ partner, isActive, onClick }) => {
  return (
    <button
      onClick={onClick}
      className={clsx(
        'w-full p-4 flex items-center gap-3 hover:bg-zinc-800/40 transition-colors text-left border-b border-zinc-900/50',
        isActive && 'bg-zinc-800/80'
      )}
    >
      <Avatar username={partner.username} />
      <div className="flex-1 min-w-0">
        <div className="text-sm font-medium text-zinc-100 truncate">
          {partner.username}
        </div>
        <div className="text-xs text-zinc-500 truncate mt-0.5">
          ID: {partner.id}
        </div>
      </div>
    </button>
  );
};

export const PartnerItem = memo(PartnerItemBase);
