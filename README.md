# 📄 DocuWhiZ – PDF Q&A Chatbot

DocuWhiZ is a full-stack PDF-based chatbot powered by Gemini AI. Users can upload PDFs, and ask questions in a chat interface. The bot responds with contextual answers extracted from the document. Built with Go (Gin), React, and Vite, with streaming Gemini responses and markdown rendering.

---

## ✨ Features

- 📤 Upload any PDF
- 💬 Chat with the document using AI
- ⚡ Fast, streaming answers from Gemini
- 🧠 Context-aware replies based on uploaded document
- 📄 Frontend built with React + Vite
- 🚀 Backend in Go with Gin framework
- 📦 Containerized using Podman or Docker

---

## 🗂️ Project Structure

```
.
├── backend
│   ├── main.go
│   ├── routes.go
│   ├── utils.go
│   └── ...
├── frontend
│   ├── index.html
│   └── src/
│       ├── App.jsx
│       └── components/
├── Dockerfile
├── .env
└── README.md
```

---

## 🚀 Running Locally

### Prerequisites

- [Go 1.21+](https://golang.org)
- [Node.js 18+](https://nodejs.org)
- `pdftotext` installed (`sudo apt install poppler-utils`)
- [Podman](https://podman.io) or Docker

---

### 1️⃣ Clone Repo

```bash
git clone https://github.com/yourname/docuwhiz.git
cd docuwhiz
```

---

### 2️⃣ Set Environment Variables

Create a `.env` file:

```
GEMINI_API_KEY=your_google_gemini_key
PORT=8000
```

---

### 3️⃣ Run Locally with Podman

```bash
podman build -t docuwhiz .
podman run -p 8000:8000 --env-file .env -v $(pwd)/uploads:/app/uploads docuwhiz
```

---

### 4️⃣ Access in Browser

Visit: [http://localhost:8000](http://localhost:8000)

---

## 🧪 Testing Locally Without Podman

### Backend

```bash
cd backend
go run main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

---

## 🌐 Deployment

You can deploy on:

- Fly.io (⚠️ needs credit card)
- Railway
- Render
- DigitalOcean App Platform
- Self-hosted VPS

---

## 🛠️ TODO

- [ ] Add authentication
- [ ] Allow multiple PDF chat history
- [ ] Improve error handling
- [ ] Add file size limit & restrictions

---

## 🧠 Built With

- 🟦 Go + Gin
- ⚛️ React + Vite
- 🧠 Gemini API
- 🧪 Markdown + Streaming
- 🐳 Podman

---

## 📄 License

MIT © 2025 Zeeshan & Team
