import { useState, useEffect, useRef, useCallback, type FC, type FormEvent } from 'react';
import { Send, UserCircle } from 'lucide-react';
import { useVirtualizer } from '@tanstack/react-virtual';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { useWebSocket } from '../hooks/useWebSocket';
import { useChatMessages } from '../hooks/useChatMessages';
import { MessageBubble } from './MessageBubble';
import { Avatar } from './ui/Avatar';
import { Input } from './ui/Input';

export const ChatWindow: FC = () => {
  const user = useAuthStore((state) => state.user);
  const activePartner = useChatStore((state) => state.activePartner);
  const addMessage = useChatStore((state) => state.addMessage);

  const { messages, isLoading, hasMore, loadMore } = useChatMessages(activePartner?.id);
  const { sendMessage } = useWebSocket();
  const [text, setText] = useState('');

  const containerRef = useRef<HTMLDivElement>(null);
  const scrollHeightRef = useRef<number>(0);

  const rowVirtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => containerRef.current,
    estimateSize: () => 72,
    overscan: 10,
  });

  useEffect(() => {
    if (containerRef.current && scrollHeightRef.current > 0) {
      const delta = containerRef.current.scrollHeight - scrollHeightRef.current;
      containerRef.current.scrollTop = delta;
      scrollHeightRef.current = 0;
    }
  }, [messages]);

  useEffect(() => {
    if (messages.length > 0 && scrollHeightRef.current === 0) {
      rowVirtualizer.scrollToIndex(messages.length - 1, { align: 'end' });
    }
  }, [messages.length, rowVirtualizer]);

  const handleScroll = useCallback(() => {
    if (containerRef.current && containerRef.current.scrollTop === 0 && hasMore && !isLoading) {
      scrollHeightRef.current = containerRef.current.scrollHeight;
      loadMore();
    }
  }, [hasMore, isLoading, loadMore]);

  const handleSend = useCallback(
    (e: FormEvent) => {
      e.preventDefault();
      if (!text.trim() || !activePartner || !user) return;

      sendMessage(activePartner.id, text.trim());

      addMessage(activePartner.id, {
        message_id: `pending-${Date.now()}-${Math.random().toString(36).slice(2)}`,
        from_id: user.id,
        to_id: activePartner.id,
        content: text.trim(),
        created_at: new Date().toISOString(),
        isPending: true,
      });

      setText('');
      setTimeout(() => {
        rowVirtualizer.scrollToIndex(messages.length, { align: 'end' });
      }, 50);
    },
    [text, activePartner, user, sendMessage, addMessage, messages.length, rowVirtualizer]
  );

  if (!activePartner) {
    return (
      <div className="flex-1 bg-zinc-950 flex flex-col items-center justify-center text-zinc-500 select-none">
        <div className="w-16 h-16 rounded-full bg-zinc-900 flex items-center justify-center mb-4 border border-zinc-800">
          <UserCircle size={32} />
        </div>
        <p className="text-sm">Select a chat to start messaging</p>
      </div>
    );
  }

  const virtualItems = rowVirtualizer.getVirtualItems();

  return (
    <div className="flex-1 flex flex-col bg-zinc-950 h-screen">
      <div className="p-4 border-b border-zinc-800 bg-zinc-900/50 flex items-center gap-3">
        <Avatar username={activePartner.username} />
        <div>
          <div className="font-medium text-zinc-100">{activePartner.username}</div>
          <div className="text-xs text-zinc-500">User ID: {activePartner.id}</div>
        </div>
      </div>

      <div
        ref={containerRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto p-4 scroll-smooth"
      >
        {isLoading && (
          <div className="text-center text-xs text-zinc-500 py-2">
            Loading messages...
          </div>
        )}

        <div
          style={{
            height: `${rowVirtualizer.getTotalSize()}px`,
            width: '100%',
            position: 'relative',
          }}
        >
          {virtualItems.map((virtualItem) => {
            const msg = messages[virtualItem.index];
            return (
              <div
                key={msg.message_id || `${msg.from_id}-${msg.created_at}`}
                data-index={virtualItem.index}
                ref={rowVirtualizer.measureElement}
                style={{
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  width: '100%',
                  transform: `translateY(${virtualItem.start}px)`,
                  paddingBottom: '12px',
                }}
              >
                <MessageBubble
                  message={msg}
                  isMe={msg.from_id === user?.id}
                />
              </div>
            );
          })}
        </div>
      </div>

      <div className="p-4 bg-zinc-900 border-t border-zinc-800">
        <form onSubmit={handleSend} className="flex gap-2">
          <Input
            type="text"
            placeholder="Write a message..."
            className="flex-1 bg-zinc-950 border border-zinc-800 rounded-full px-5 py-2.5 text-sm"
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
          <button
            type="submit"
            disabled={!text.trim()}
            className="p-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0 flex items-center justify-center"
          >
            <Send size={18} />
          </button>
        </form>
      </div>
    </div>
  );
};
