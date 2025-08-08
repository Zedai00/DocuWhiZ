import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import ReactMarkdown from 'react-markdown';
import './ChatBox.css';

function ChatBox({ selectedFileId, initialMessage }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [typing, setTyping] = useState(false);
  const messagesEndRef = useRef(null);
  const hasShownWelcome = useRef(false);

  useEffect(() => {
    scrollToBottom();
  }, [messages, typing]);

  useEffect(() => {
    if (!initialMessage) return;

    if (initialMessage.includes('Welcome') && !hasShownWelcome.current) {
      setMessages([{ sender: 'docuwhiz', text: initialMessage }]);
      hasShownWelcome.current = true;
    } else if (
      initialMessage.includes('PDF uploaded') &&
      !messages.some((msg) => msg.text.includes('PDF uploaded'))
    ) {
      setMessages((prev) => [
        ...prev,
        { sender: 'docuwhiz', text: initialMessage },
      ]);
    }
  }, [initialMessage]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const simulateStreaming = (text, onUpdate, onFirstChunk, onDone) => {
    const words = text.split(' ');
    let index = 0;
    let hasStarted = false;

    const interval = setInterval(() => {
      if (index < words.length) {
        const partial = words.slice(0, index + 1).join(' ');
        onUpdate(partial);

        if (!hasStarted && partial.length > 0) {
          hasStarted = true;
          onFirstChunk();
        }

        index++;
      } else {
        clearInterval(interval);
        onDone();
      }
    }, 80);
  };

  const sendMessage = async () => {
    const trimmedInput = input.trim();
    if (!trimmedInput) return;

    if (!selectedFileId) {
      setMessages((prev) => [
        ...prev,
        { sender: 'docuwhiz', text: '⚠️ Please upload a PDF first.' },
      ]);
      return;
    }

    const userMessage = { sender: 'user', text: trimmedInput };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setTyping(true);

    const placeholderIndex = messages.length + 1;
    setMessages((prev) => [...prev, { sender: 'docuwhiz', text: '' }]);

    try {
      const response = await axios.post('/api/chat', {
        fileId: selectedFileId,
        message: trimmedInput,
      });

      const fullAnswer = response.data.response;

      simulateStreaming(
        fullAnswer,
        (partial) => {
          setMessages((prev) =>
            prev.map((msg, i) =>
              i === placeholderIndex ? { ...msg, text: partial } : msg
            )
          );
        },
        () => setTyping(false),
        () => { }
      );
    } catch (error) {
      console.error(error);
      setTyping(false);
      setMessages((prev) => [
        ...prev,
        {
          sender: 'docuwhiz',
          text: '⚠️ Error fetching response. Try again later.',
        },
      ]);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className="chatbox-container">
      <div className="chat-messages">
        {messages.map((msg, idx) =>
          msg.text.trim() ? (
            <div
              key={idx}
              className={`chat-message ${msg.sender === 'user' ? 'user-message' : 'bot-message'}`}
            >
              {msg.sender === 'docuwhiz' ? (
                <div className="markdown-content">
                  <ReactMarkdown>{msg.text}</ReactMarkdown>
                </div>
              ) : (
                msg.text
              )}
            </div>
          ) : null
        )}

        {typing && (
          <div className="typing-indicator">
            <div className="spinner" />
            DocuWhiZ is typing...
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      <div className="chat-input">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Ask your question about the PDF..."
        />
        <button onClick={sendMessage}>Send</button>
      </div>
    </div>
  );
}

export default ChatBox;

