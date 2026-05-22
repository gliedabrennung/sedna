import { useEffect, useRef, useCallback } from 'react';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';

const MAX_RECONNECT_DELAY = 30_000;
const BASE_DELAY = 1_000;

export function useWebSocket() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const userId = useAuthStore((s) => s.user?.id);
  const wsRef = useRef<WebSocket | null>(null);
  const retriesRef = useRef(0);
  const mountedRef = useRef(true);

  useEffect(() => {
    mountedRef.current = true;
    if (!isAuthenticated || !userId) return;

    const connect = () => {
      if (!mountedRef.current) return;

      const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
      const ws = new WebSocket(`${protocol}//${location.host}/ws`);
      wsRef.current = ws;

      ws.onopen = () => {
        retriesRef.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (!data.from || !data.message) return;

          const myId = useAuthStore.getState().user?.id;
          if (data.from === myId) return;

          const partnerId = data.from;
          useChatStore.getState().addMessage(partnerId, {
            from_id: data.from,
            to_id: data.to,
            content: data.message,
            created_at: new Date().toISOString(),
          });
        } catch {
          /* malformed message */
        }
      };

      ws.onclose = () => {
        wsRef.current = null;
        if (!mountedRef.current) return;
        const delay = Math.min(BASE_DELAY * 2 ** retriesRef.current, MAX_RECONNECT_DELAY);
        retriesRef.current++;
        setTimeout(connect, delay);
      };

      ws.onerror = () => ws.close();
    };

    connect();

    return () => {
      mountedRef.current = false;
      if (wsRef.current) {
        wsRef.current.onclose = null;
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [isAuthenticated, userId]);

  const sendMessage = useCallback((toId: number, message: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ to: toId, message }));
    }
  }, []);

  return { sendMessage };
}
