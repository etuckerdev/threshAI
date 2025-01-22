import React, { useState, useEffect, useRef } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { python } from '@codemirror/lang-python';
import { cpp } from '@codemirror/lang-cpp';
import { rust } from '@codemirror/lang-rust';
import { java } from '@codemirror/lang-java';
import { oneDark } from '@codemirror/theme-one-dark';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const languageMap = {
  javascript: javascript,
  python: python,
  cpp: cpp,
  rust: rust,
  java: java,
};

const ChatTerminal = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const messagesEndRef = useRef(null);

  // Auto-scroll to bottom when new messages arrive
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Copy code to clipboard
  const copyToClipboard = async (code) => {
    try {
      await navigator.clipboard.writeText(code);
      alert('Code copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy code:', err);
    }
  };

  // Handle feedback submission
  const handleFeedback = async (messageId, isPositive) => {
    try {
      const response = await fetch(`${API_URL}/api/feedback`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messageId, isPositive }),
        credentials: 'include',
      });
      
      if (!response.ok) throw new Error('Feedback submission failed');
      
      // Update message to show feedback was received
      setMessages(prev => prev.map(msg => 
        msg.id === messageId 
          ? { ...msg, feedback: isPositive ? 'positive' : 'negative' }
          : msg
      ));
    } catch (err) {
      console.error('Failed to submit feedback:', err);
      alert('Failed to submit feedback. Please try again.');
    }
  };

  // Send message to backend
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!input.trim() || isProcessing) return;

    setIsProcessing(true);
    const userMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: input,
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');

    try {
      const response = await fetch(`${API_URL}/api/generate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          prompt: input,
          userId: 'user-' + Math.random().toString(36).substr(2, 9)
        }),
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('API request failed');
      }

      const data = await response.json();
      
      if (!data.code) {
        throw new Error('Invalid response format');
      }

      setMessages(prev => [...prev, {
        id: Date.now().toString(),
        type: 'system',
        content: data.code,
        language: data.language || 'javascript',
      }]);
    } catch (err) {
      console.error('Failed to generate response:', err);
      setMessages(prev => [...prev, {
        id: Date.now().toString(),
        type: 'error',
        content: 'Failed to generate response. Please try again.',
      }]);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="chat-terminal">
      <div className="messages">
        {messages.map(msg => (
          <div key={msg.id} className={`message ${msg.type}`}>
            {msg.type === 'system' && msg.content ? (
              <div className="code-block">
                <div className="code-header">
                  <span className="language">{msg.language}</span>
                  <button
                    onClick={() => copyToClipboard(msg.content)}
                    className="copy-button"
                  >
                    📋 Copy
                  </button>
                </div>
                <CodeMirror
                  value={msg.content}
                  height="auto"
                  theme={oneDark}
                  extensions={[languageMap[msg.language || 'javascript']?.()]}
                  editable={false}
                />
                <div className="feedback-buttons">
                  <button
                    onClick={() => handleFeedback(msg.id, true)}
                    className={msg.feedback === 'positive' ? 'active' : ''}
                  >
                    👍
                  </button>
                  <button
                    onClick={() => handleFeedback(msg.id, false)}
                    className={msg.feedback === 'negative' ? 'active' : ''}
                  >
                    👎
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-content">{msg.content}</div>
            )}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={handleSubmit} className="input-form">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Describe the code you want to generate..."
          disabled={isProcessing}
        />
        <button type="submit" disabled={isProcessing}>
          {isProcessing ? 'Generating...' : 'Generate'}
        </button>
      </form>
    </div>
  );
};

export default ChatTerminal;