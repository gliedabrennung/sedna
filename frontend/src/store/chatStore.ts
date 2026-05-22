import { create } from 'zustand';
import type { User, Message } from '@/types';

const MAX_MESSAGES_PER_CHAT = 500;
const RECENT_CHATS_KEY = 'messenger_recent_chats';

interface ChatState {
  activePartner: User | null;
  setActivePartner: (user: User | null) => void;
  messages: Record<number, Message[]>;
  setMessages: (partnerId: number, messages: Message[]) => void;
  addMessage: (partnerId: number, message: Message) => void;
  recentChats: User[];
  setRecentChats: (chats: User[]) => void;
  addRecentChat: (user: User) => void;
}

function loadRecentChats(): User[] {
  try {
    const raw = localStorage.getItem(RECENT_CHATS_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function persistRecentChats(chats: User[]) {
  localStorage.setItem(RECENT_CHATS_KEY, JSON.stringify(chats.slice(0, 50)));
}

export const useChatStore = create<ChatState>((set) => ({
  activePartner: null,
  setActivePartner: (user) => set({ activePartner: user }),

  messages: {},
  setMessages: (partnerId, msgs) =>
    set((state) => ({
      messages: {
        ...state.messages,
        [partnerId]: msgs.length > MAX_MESSAGES_PER_CHAT
          ? msgs.slice(msgs.length - MAX_MESSAGES_PER_CHAT)
          : msgs,
      },
    })),

  addMessage: (partnerId, msg) =>
    set((state) => {
      const existing = state.messages[partnerId] || [];
      if (msg.message_id && existing.some((m) => m.message_id === msg.message_id)) {
        return state;
      }
      const updated = [...existing, msg];
      return {
        messages: {
          ...state.messages,
          [partnerId]: updated.length > MAX_MESSAGES_PER_CHAT
            ? updated.slice(updated.length - MAX_MESSAGES_PER_CHAT)
            : updated,
        },
      };
    }),

  recentChats: loadRecentChats(),
  setRecentChats: (chats) => {
    persistRecentChats(chats);
    set({ recentChats: chats });
  },
  addRecentChat: (user) =>
    set((state) => {
      const filtered = state.recentChats.filter((c) => c.id !== user.id);
      const updated = [user, ...filtered];
      persistRecentChats(updated);
      return { recentChats: updated };
    }),
}));
