import React, { useState } from 'react';
import axios from 'axios';
import './FileUpload.css';

function FileUpload({ onUpload }) {
  const [file, setFile] = useState(null);
  const [status, setStatus] = useState('');

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    setStatus('');
  };

  const handleUpload = async () => {
    if (!file || file.type !== 'application/pdf') {
      setStatus('❌ Please select a valid PDF file.');
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await axios.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      const fileId = res.data.fileId;
      onUpload(fileId); // Send to parent component
      setStatus('✅ Upload successful!');
    } catch (err) {
      console.error(err);
      setStatus('❌ Upload failed: ' + (err.response?.data?.error || err.message));
    }
  };

  return (
    <div className="upload-container">
      <input
        type="file"
        accept="application/pdf"
        onChange={handleFileChange}
        id="fileInput"
        hidden
      />
      <label htmlFor="fileInput" className="upload-button">
        📄 Choose PDF
      </label>
      {file && <span className="file-name">{file.name}</span>}

      <button
        onClick={handleUpload}
        className="submit-button"
        disabled={!file}
      >
        {file ? '⬆️ Upload' : '📂 Select PDF First'}
      </button>

      {status && <p className="status-message">{status}</p>}
    </div>
  );
}

export default FileUpload;

