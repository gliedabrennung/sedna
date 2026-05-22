import { create } from 'zustand';
import type { User, Message } from '../types';

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

export const useChatStore = create<ChatState>((set) => ({
  activePartner: null,
  setActivePartner: (user) => set({ activePartner: user }),
  messages: {},
  setMessages: (partnerId, msgs) =>
    set((state) => ({
      messages: { ...state.messages, [partnerId]: msgs },
    })),
  addMessage: (partnerId, msg) =>
    set((state) => {
      const existing = state.messages[partnerId] || [];
      if (msg.message_id && existing.some((m) => m.message_id === msg.message_id)) {
        return state;
      }
      return {
        messages: { ...state.messages, [partnerId]: [...existing, msg] },
      };
    }),
  recentChats: [],
  setRecentChats: (chats) => set({ recentChats: chats }),
  addRecentChat: (user) =>
    set((state) => {
      const filtered = state.recentChats.filter((c) => c.id !== user.id);
      return { recentChats: [user, ...filtered] };
    }),
}));
