import { useState, useCallback, type FC } from 'react';
import { Search, LogOut } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { useSearchUsers } from '../hooks/useSearchUsers';
import { PartnerItem } from './PartnerItem';
import { Avatar } from './ui/Avatar';
import { Input } from './ui/Input';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import type { User } from '../types';

export const Sidebar: FC = () => {
  const [search, setSearch] = useState('');
  const { results: searchResults, isLoading: isSearching } = useSearchUsers(search);
  
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  
  const activePartner = useChatStore((state) => state.activePartner);
  const setActivePartner = useChatStore((state) => state.setActivePartner);
  const recentChats = useChatStore((state) => state.recentChats);
  const addRecentChat = useChatStore((state) => state.addRecentChat);
  
  const navigate = useNavigate();

  const handleSelectPartner = useCallback(
    (partner: User) => {
      setActivePartner(partner);
      addRecentChat(partner);
      setSearch('');
    },
    [setActivePartner, addRecentChat]
  );

  const handleLogout = useCallback(async () => {
    try { await api.post('/auth/logout'); } catch {}
    logout();
    navigate('/login');
  }, [logout, navigate]);

  const displayList = search ? searchResults : recentChats;

  return (
    <div className="w-80 flex-shrink-0 bg-zinc-900 border-r border-zinc-800 flex flex-col h-screen">
      <div className="p-4 border-b border-zinc-800 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Avatar username={user?.username} />
          <div>
            <div className="text-sm font-medium text-zinc-100">{user?.username}</div>
            <div className="text-xs text-zinc-500">Online</div>
          </div>
        </div>
        <button
          onClick={handleLogout}
          className="p-2 text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800 rounded-lg transition-colors"
        >
          <LogOut size={18} />
        </button>
      </div>

      <div className="p-4">
        <div className="relative">
          <Search
            className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500"
            size={16}
          />
          <Input
            type="text"
            placeholder="Search users..."
            className="pl-9"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        {isSearching && (
          <div className="p-4 text-center text-xs text-zinc-500">
            Searching...
          </div>
        )}
        
        {!isSearching && displayList.length === 0 && (
          <div className="p-4 text-center text-sm text-zinc-500">
            {search ? 'No users found' : 'No recent chats'}
          </div>
        )}

        {!isSearching &&
          displayList.map((partner) => (
            <PartnerItem
              key={partner.id}
              partner={partner}
              isActive={activePartner?.id === partner.id}
              onClick={() => handleSelectPartner(partner)}
            />
          ))}
      </div>
    </div>
  );
};
