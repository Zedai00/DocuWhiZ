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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var extractedPDFText string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/upload", uploadPDFHandler)
	router.POST("/chat", func(c *gin.Context) {
		handleGeminiChat(c, apiKey)
	})

	router.Run(":8000")
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
	file, err := c.FormFile("pdf")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read PDF: %v", err)
		return
	}

	savePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %v", err)
		return
	}

	text, err := extractTextFromPDF(savePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to extract text: %v", err)
		return
	}

	extractedPDFText = text
	c.HTML(http.StatusOK, "index.html", gin.H{"message": "PDF uploaded and text extracted successfully!"})
}

func handleGeminiChat(c *gin.Context, apiKey string) {
	query := c.PostForm("query")
	if query == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"message": "No query provided."})
		return
	}
	if extractedPDFText == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"message": "Please upload a PDF first."})
		return
	}

	prompt := fmt.Sprintf("Answer the following based on this PDF:\n\n%s\n\nQuestion: %s", extractedPDFText, query)
	answer, err := callGeminiAPI(apiKey, prompt)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"message": "Gemini API call failed."})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"response": answer})
}

// ======= Gemini REST API =======

func callGeminiAPI(apiKey, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-2.5-pro:generateContent?key=%s", apiKey)

	// Build request body
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

	// Marshal JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error: %s", string(bodyBytes))
	}

	// Decode JSON response
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid response from Gemini")
	}

	// Return the response text
	return result.Candidates[0].Content.Parts[0].Text, nil
}
