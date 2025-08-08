package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SessionData struct {
	PDFText string
}

var (
	sessions = make(map[string]*SessionData)
	mutex    = &sync.Mutex{}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Falling back to system environment variables.")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API route
	// --- CORRECTED STATIC FILE SERVING ---
	// 1. Serve the SPA's main HTML file directly from the root URL.
	router.StaticFile("/", "./dist/index.html")

	// 2. Serve the entire 'assets' directory from the '/assets' URL path.
	// This will correctly match requests like '/assets/index-DPuVEiKz.js'.
	router.Static("/assets", "./dist/assets")

	// 3. The NoRoute handler is a fallback for SPA routing.
	// Any non-API, non-static-file URL will serve the main index.html file.
	router.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	router.POST("/api/upload", uploadPDFHandler)
	router.POST("/api/chat", func(c *gin.Context) {
		handleGeminiChat(c, apiKey)
	})

	router.Run(":" + port)
}

func extractTextFromPDF(pdfPath string) (string, error) {
	cmd := exec.Command("pdftotext", "-layout", pdfPath, "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %v", err)
	}
	return out.String(), nil
}

func uploadPDFHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file required"})
		return
	}

	// Ensure 'uploads' directory exists
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
		return
	}

	// Save file
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join(uploadsDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDF"})
		return
	}

	text, err := extractTextFromPDF(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract PDF text"})
		return
	}

	if err := os.Remove(savePath); err != nil {
		log.Printf("Warning: failed to delete uploaded file %s: %v", savePath, err)
	}

	// Save to session
	mutex.Lock()
	sessions[filename] = &SessionData{PDFText: text}
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "PDF uploaded and text extracted",
		"fileId":  filename,
	})
}

func handleGeminiChat(c *gin.Context, apiKey string) {
	var payload struct {
		FileID  string `json:"fileId"`
		Message string `json:"message"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	mutex.Lock()
	data, exists := sessions[payload.FileID]
	mutex.Unlock()

	if !exists || data.PDFText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF text not found for session"})
		return
	}

	prompt := fmt.Sprintf("Answer based on this PDF:\n\n%s\n\nQuestion: %s", data.PDFText, payload.Message)
	answer, err := callGeminiAPI(apiKey, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gemini API call failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": answer})
}

func callGeminiAPI(apiKey, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-2.5-pro:generateContent?key=%s", apiKey)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini error: %s", string(body))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				}
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid response from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
