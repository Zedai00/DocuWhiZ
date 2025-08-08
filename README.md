# 📄 DocuWhiZ – AI PDF Q&A Chatbot

**DocuWhiZ** is an open-source, full-stack chatbot that lets you **chat with your PDFs**. Built with Go and React, it uses the **Gemini API** to provide accurate, streaming, markdown-rendered answers from your documents.

> 🧠 Upload any PDF. Ask anything. Get AI-powered contextual answers instantly.

---

## ✨ Features

- 📤 Upload and process any PDF document
- 💬 Chat in real time with document context
- ⚡ Gemini API streaming response support
- 🧠 Smart question-answering with LLM context
- 🖥️ Frontend: React + Vite (Markdown-rendered UI)
- 🧩 Backend: Go + Gin (PDF parsing & chat memory)
- 📦 Fully containerized (Podman / Docker)
- 🧾 Deployed at: [https://docuwhiz.onrender.com](https://docuwhiz.onrender.com)

---

## 🧠 Tech Stack

- 🔹 Go 1.21+ with Gin Web Framework
- ⚛️ React 18 + Vite
- 🧠 Gemini Pro API
- 📄 Poppler (`pdftotext`) for PDF parsing
- 🐳 Podman or Docker for deployment
- 🧪 Markdown for chat formatting

---

## 📁 Project Structure

```
docuwhiz/
├── backend/
│   ├── main.go          # Gin server entry
├── frontend/
│   ├── index.html
│   └── src/
│       ├── App.jsx
│       ├── components/
│       │   ├── ChatBox.jsx
│       │   └── FileUpload.jsx
│       └── styles/
├── Dockerfile
├── .env
└── README.md
```

---

## 🚀 Getting Started

### ✅ Prerequisites

- [Go 1.21+](https://golang.org)
- [Node.js 18+](https://nodejs.org)
- [`pdftotext`](https://poppler.freedesktop.org/) (`sudo apt install poppler-utils`)
- [Podman](https://podman.io) or Docker

---

### 1️⃣ Clone the Repo

```bash
git clone https://github.com/zedai00/docuwhiz.git
cd docuwhiz
```

---

### 2️⃣ Set Environment Variables

Create a `.env` file in the root directory:

```env
GEMINI_API_KEY=your_google_gemini_api_key
PORT=8000
```

---

### 3️⃣ Build & Run with Podman (or Docker)

```bash
podman build -t docuwhiz .
podman run -p 8000:8000 --env-file .env -v $(pwd)/uploads:/app/uploads docuwhiz
```

Visit 👉 [http://localhost:8000](http://localhost:8000)

---

## 🧪 Development Setup

### Backend (Go + Gin)

```bash
cd backend
go run main.go
```

### Frontend (React + Vite)

```bash
cd frontend
npm install
npm run dev
```

---

## 🌐 Deployment Options

You can deploy DocuWhiZ on:

| Platform                                      | Free Plan | Notes                       |
| --------------------------------------------- | --------- | --------------------------- |
| [Render](https://render.com)                  | ✅        | Used in production          |
| [Railway](https://railway.app)                | ✅        | Easy CI/CD setup            |
| [Fly.io](https://fly.io)                      | ⚠️        | Needs credit card           |
| [DigitalOcean](https://www.digitalocean.com/) | ✅        | Self-hosted or App Platform |
| Any VPS (e.g., EC2, Hetzner)                  | ✅        | Use Docker + nginx          |

---

## 🛠️ TODO

- [ ] User authentication & session history
- [ ] Support multiple PDFs
- [ ] Add chat memory across PDFs
- [ ] Upload limit handling
- [ ] Mobile responsive layout

---

## 👨‍💻 Maintainers

Built by **Zed**
🌐 [LinkedIn](https://linkedin.com/in/zedai00) | 🐙 [GitHub](https://github.com/zedai00)

---

## 📄 License

MIT License © 2025 Zed

---

> 🚀 Try it live: [https://docuwhiz.onrender.com](https://docuwhiz.onrender.com)
