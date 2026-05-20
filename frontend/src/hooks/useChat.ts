import { useState, useEffect, useCallback, useRef } from 'react';
import { useAuth } from '../context/AuthContext';
import type { User as AuthUser } from '../context/AuthContext';
import { WebSocketService } from '../services/websocket';
import { request } from '../services/api';

export interface Message {
  id: string;
  from: number;
  to: number;
  type?: string;
  message: string;
  status?: string;
  timestamp: Date;
}

export function useChat() {
  const { user, token } = useAuth();
  const [ws, setWs] = useState<WebSocketService | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [contacts, setContacts] = useState<AuthUser[]>([]);
  const [activeContactId, setActiveContactId] = useState<number | null>(null);
  const [isLoadingContacts, setIsLoadingContacts] = useState(true);
  const [unreadCounts, setUnreadCounts] = useState<Record<number, number>>({});
  
  const activeContactIdRef = useRef<number | null>(null);
  useEffect(() => {
    activeContactIdRef.current = activeContactId;
  }, [activeContactId]);

  // Fetch contacts and unread counts
  useEffect(() => {
    let isMounted = true;
    const fetchData = async () => {
      setIsLoadingContacts(true);
      try {
        const [usersData, unreadData] = await Promise.all([
          request('/api/users'),
          request('/api/messages/unread')
        ]);
        if (isMounted) {
          const otherUsers = (usersData.users || []).filter((u: AuthUser) => u.id !== user?.id);
          setContacts(otherUsers);
          setUnreadCounts(unreadData.counts || {});
        }
      } catch (err) {
        console.error('Failed to fetch data', err);
      } finally {
        if (isMounted) setIsLoadingContacts(false);
      }
    };
    
    if (user) {
      fetchData();
    }

    return () => {
      isMounted = false;
    };
  }, [user]);

  // Fetch History
  useEffect(() => {
    let isMounted = true;
    const fetchHistory = async () => {
      if (!activeContactId) {
        setMessages([]);
        return;
      }
      try {
        const data = await request(`/api/messages?contact_id=${activeContactId}`);
        if (isMounted) {
          const msgs: Message[] = (data.messages || []).map((m: any) => ({
            id: String(m.id),
            from: m.from_id,
            to: m.to_id,
            type: m.type || 'text',
            message: m.content,
            status: m.status || 'sent',
            timestamp: new Date(m.created_at)
          }));
          setMessages(msgs);
        }
      } catch (err) {
        console.error('Failed to fetch history', err);
      }
    };
    
    fetchHistory();
    return () => { isMounted = false; };
  }, [activeContactId]);

  // Connect to WebSocket
  useEffect(() => {
    if (token) {
      const wsService = new WebSocketService(token);
      wsService.connect();
      setWs(wsService);

      const unsubscribe = wsService.onMessage((msg: any) => {
        if (msg.action === 'status_update') {
          setMessages(prev => prev.map(m => {
            if (m.to === msg.from && m.status !== 'read') {
              return { ...m, status: msg.status };
            }
            return m;
          }));
          return;
        }

        const newMessage: Message = {
          id: Math.random().toString(36).substr(2, 9),
          from: msg.from,
          to: msg.to,
          type: msg.type || 'text',
          message: msg.message,
          status: msg.status || 'delivered',
          timestamp: new Date()
        };

        if (msg.from !== user?.id) {
          if (msg.from === activeContactIdRef.current) {
            newMessage.status = 'read';
            wsService.sendStatusUpdate(msg.from, 'read');
          } else {
            wsService.sendStatusUpdate(msg.from, 'delivered');
            setUnreadCounts(prev => ({ ...prev, [msg.from]: (prev[msg.from] || 0) + 1 }));
          }
        }

        setMessages(prev => [...prev, newMessage]);
      });

      return () => {
        unsubscribe();
        wsService.disconnect();
      };
    }
  }, [token, user?.id]);

  // Handle opening a chat to mark as read
  useEffect(() => {
    if (activeContactId && ws && unreadCounts[activeContactId] > 0) {
      ws.sendStatusUpdate(activeContactId, 'read');
      setUnreadCounts(prev => ({ ...prev, [activeContactId]: 0 }));
      setMessages(prev => prev.map(m => {
        if (m.from === activeContactId && m.status !== 'read') {
          return { ...m, status: 'read' };
        }
        return m;
      }));
    }
  }, [activeContactId, ws, unreadCounts]);

  const sendMessage = useCallback((content: string, type: string = 'text') => {
    if (!content.trim() || !activeContactId || !ws || !user) return false;

    ws.sendMessage(activeContactId, content, type);
    
    const newMessage: Message = {
      id: Math.random().toString(36).substr(2, 9),
      from: user.id,
      to: activeContactId,
      type,
      message: content,
      status: 'sent',
      timestamp: new Date()
    };
    setMessages(prev => [...prev, newMessage]);
    return true;
  }, [ws, activeContactId, user]);

  const activeContact = contacts.find(c => c.id === activeContactId) || null;
  
  const visibleMessages = messages.filter(
    m => (m.from === user?.id && m.to === activeContactId) || 
         (m.from === activeContactId && m.to === user?.id)
  );

  return {
    contacts,
    activeContactId,
    setActiveContactId,
    activeContact,
    messages: visibleMessages,
    sendMessage,
    isLoadingContacts,
    user,
    unreadCounts
  };
}
