import { memo, useMemo, type FC } from 'react';
import clsx from 'clsx';
import type { Message } from '../types';

interface MessageBubbleProps {
  message: Message;
  isMe: boolean;
}

const MAX_DISPLAY_LENGTH = 4000;

function sanitizeContent(raw: string): string {
  if (raw.length > MAX_DISPLAY_LENGTH) {
    return raw.slice(0, MAX_DISPLAY_LENGTH) + '…';
  }
  return raw;
}

const MessageBubbleBase: FC<MessageBubbleProps> = ({ message, isMe }) => {
  const content = sanitizeContent(message.content);

  const formattedTime = useMemo(() => {
    if (!message.created_at) return '';
    return new Date(message.created_at).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
  }, [message.created_at]);

  return (
    <div className={clsx('flex', isMe ? 'justify-end' : 'justify-start')}>
      <div
        className={clsx(
          'max-w-[70%] rounded-2xl px-4 py-2 shadow-sm',
          isMe ? 'bg-indigo-600 text-white rounded-br-none' : 'bg-zinc-800 text-zinc-100 rounded-bl-none',
          message.isPending && 'opacity-70'
        )}
      >
        <p className="text-sm break-words whitespace-pre-wrap">{content}</p>
        <div
          className={clsx(
            'text-[10px] mt-1 text-right select-none',
            isMe ? 'text-indigo-200' : 'text-zinc-500'
          )}
        >
          {formattedTime}
        </div>
      </div>
    </div>
  );
};

export const MessageBubble = memo(MessageBubbleBase);
