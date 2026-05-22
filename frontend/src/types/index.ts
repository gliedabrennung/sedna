export interface User {
  id: number;
  username: string;
}

export interface Message {
  chat_id?: string;
  message_id?: string;
  from_id: number;
  to_id: number;
  content: string;
  created_at?: string;
  isPending?: boolean;
}

export interface ChatHistoryResponse {
  messages: Message[];
  next_cursor: string | null;
}
