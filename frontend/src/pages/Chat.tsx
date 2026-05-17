import { useState, useEffect, useRef } from 'react';
import { Send, LogOut, Check, CheckCheck } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useChat } from '../hooks/useChat';
import AudioRecorder from '../components/AudioRecorder';
import { request } from '../services/api';

export default function Chat() {
  const { logout } = useAuth();
  const { 
    contacts, 
    activeContactId, 
    setActiveContactId, 
    activeContact, 
    messages, 
    sendMessage, 
    user,
    unreadCounts
  } = useChat();

  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault();
    const success = sendMessage(inputValue);
    if (success) {
      setInputValue('');
    }
  };

  const handleAudioSend = async (audioBlob: Blob) => {
    try {
      const formData = new FormData();
      formData.append('audio', audioBlob, 'audio.webm');
      
      const response = await request('/api/messages/upload', {
        method: 'POST',
        body: formData
      });
      
      if (response && response.url) {
        sendMessage(response.url, 'audio');
      }
    } catch (err) {
      console.error('Failed to upload audio:', err);
      alert('Failed to send voice message.');
    }
  };

  return (
    <div className="app-container">
      <div className="sidebar">
        <div className="sidebar-header">
          <h2>Chats</h2>
          <button onClick={logout} className="btn" style={{ padding: '8px' }} title="Logout">
            <LogOut size={18} />
          </button>
        </div>
        <div className="sidebar-content">
          {contacts.map(contact => (
            <div 
              key={contact.id} 
              className={`contact-item ${activeContactId === contact.id ? 'active' : ''}`}
              onClick={() => setActiveContactId(contact.id)}
            >
              <div className="contact-avatar">
                {contact.username.charAt(0).toUpperCase()}
              </div>
              <div className="contact-info" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                <div className="contact-name">{contact.username}</div>
                {unreadCounts && unreadCounts[contact.id] > 0 && (
                  <div style={{
                    backgroundColor: '#34b7f1',
                    color: 'white',
                    borderRadius: '50%',
                    padding: '2px 8px',
                    fontSize: '12px',
                    fontWeight: 'bold'
                  }}>
                    {unreadCounts[contact.id]}
                  </div>
                )}
              </div>
            </div>
          ))}
          {contacts.length === 0 && (
            <div style={{ padding: '20px', textAlign: 'center', color: 'var(--text-muted)' }}>
              No other users found
            </div>
          )}
        </div>
      </div>
      
      <div className="chat-area">
        {activeContactId ? (
          <>
            <div className="chat-header">
              <div className="contact-avatar" style={{ width: 36, height: 36, fontSize: 14 }}>
                {activeContact?.username.charAt(0).toUpperCase()}
              </div>
              <h3 style={{ marginLeft: 12 }}>{activeContact?.username}</h3>
            </div>
            
            <div className="chat-messages">
              {messages.map(msg => (
                <div key={msg.id} className={`message ${msg.from === user?.id ? 'sent' : 'received'}`}>
                  {msg.type === 'audio' ? (
                    <audio controls src={msg.message} style={{ maxWidth: '200px' }} />
                  ) : (
                    msg.message
                  )}
                  <div style={{ 
                    display: 'flex', 
                    justifyContent: 'flex-end', 
                    alignItems: 'center', 
                    gap: '4px',
                    fontSize: '11px',
                    opacity: 0.7,
                    marginTop: '4px'
                  }}>
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    {msg.from === user?.id && msg.status && (
                      <span style={{ display: 'flex', alignItems: 'center' }}>
                        {msg.status === 'sent' && <Check size={14} />}
                        {msg.status === 'delivered' && <CheckCheck size={14} />}
                        {msg.status === 'read' && <CheckCheck size={14} style={{ color: '#34b7f1' }} />}
                      </span>
                    )}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
            
            <div className="chat-input-area">
              <form className="chat-input-form" onSubmit={handleSend}>
                <input 
                  type="text" 
                  className="form-input" 
                  placeholder="Type a message..." 
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                />
                <button type="submit" className="btn" disabled={!inputValue.trim()}>
                  <Send size={18} />
                </button>
                <AudioRecorder 
                  onAudioReady={handleAudioSend} 
                  disabled={!activeContactId} 
                />
              </form>
            </div>
          </>
        ) : (
          <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)' }}>
            Select a chat to start messaging
          </div>
        )}
      </div>
    </div>
  );
}
