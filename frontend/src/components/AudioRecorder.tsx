import { useState, useRef } from 'react';
import { Mic, Square, Loader } from 'lucide-react';

interface AudioRecorderProps {
  onAudioReady: (audioBlob: Blob) => void;
  disabled?: boolean;
}

export default function AudioRecorder({ onAudioReady, disabled }: AudioRecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chunksRef = useRef<BlobPart[]>([]);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream);
      mediaRecorderRef.current = mediaRecorder;
      chunksRef.current = [];

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) {
          chunksRef.current.push(e.data);
        }
      };

      mediaRecorder.onstop = () => {
        setIsProcessing(true);
        const audioBlob = new Blob(chunksRef.current, { type: 'audio/webm' });
        onAudioReady(audioBlob);
        setIsProcessing(false);
        // Stop all tracks to release mic
        stream.getTracks().forEach(track => track.stop());
      };

      mediaRecorder.start();
      setIsRecording(true);
    } catch (err) {
      console.error('Error accessing microphone:', err);
      alert('Could not access microphone. Please check permissions.');
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
  };

  if (isProcessing) {
    return (
      <button type="button" className="btn" disabled>
        <Loader size={18} className="animate-spin" />
      </button>
    );
  }

  if (isRecording) {
    return (
      <button 
        type="button" 
        className="btn" 
        style={{ background: 'var(--danger)' }} 
        onClick={stopRecording}
      >
        <Square size={18} />
      </button>
    );
  }

  return (
    <button 
      type="button" 
      className="btn" 
      style={{ background: 'var(--border)', color: 'var(--text-main)' }}
      onClick={startRecording}
      disabled={disabled}
      title="Hold or click to record"
    >
      <Mic size={18} />
    </button>
  );
}
