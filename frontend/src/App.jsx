import React, { useState } from 'react';
import FileUpload from './components/FileUpload';
import ChatBox from './components/ChatBox';
import './App.css';

function App() {
  const [selectedFileId, setSelectedFileId] = useState(null);
  const [initialMessage, setInitialMessage] = useState("👋 Welcome to DocuWhiZ! Please upload a PDF to get started.");

  const handleUpload = (id) => {
    setSelectedFileId(id);
    setInitialMessage("✅ PDF uploaded. You can start chatting now.");
  };

  return (
    <div className="app-container">
      <div className="left-panel">
        <h2>📄 Upload PDF</h2>
        <FileUpload onUpload={handleUpload} />
      </div>

      <div className="right-panel">
        <h2>💬 Chat with DocuWhiZ</h2>
        <ChatBox selectedFileId={selectedFileId} initialMessage={initialMessage} />
      </div>
    </div>
  );
}

export default App;

