package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
}

type GeminiContentPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiContentPart `json:"parts"`
}

type GeminiGenerateRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type GeminiGenerateResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *GeminiClient) GenerateGadgetRecommendation(userPrompt string) (string, error) {
	if c.apiKey == "" {
		return c.generateFallbackRecommendation(userPrompt), nil
	}

	systemContext := "Anda adalah Asisten AI Spesialis Gadget Marketplace yang cerdas, sopan, dan berpengalaman. " +
		"Tugas Anda adalah memberikan rekomendasi gadget (laptop, smartphone, kamera, console) dan perbandingan harga " +
		"yang relevan, hemat, serta sesuai dengan kebutuhan dan budget pengguna. Jawab dengan bahasa Indonesia yang jelas, menarik, dan informatif.\n\n" +
		"Pertanyaan Pengguna: " + userPrompt

	reqBody := GeminiGenerateRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiContentPart{
					{Text: systemContext},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", c.apiKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Warning: Gemini API connection error (%v). Serving smart fallback recommendation.", err)
		return c.generateFallbackRecommendation(userPrompt), nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Warning: Gemini API returned status %d. Serving smart AI recommendation fallback.", resp.StatusCode)
		return c.generateFallbackRecommendation(userPrompt), nil
	}

	var geminiResp GeminiGenerateResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return c.generateFallbackRecommendation(userPrompt), nil
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return c.generateFallbackRecommendation(userPrompt), nil
}

func (c *GeminiClient) generateFallbackRecommendation(userPrompt string) string {
	lowerPrompt := strings.ToLower(userPrompt)

	if strings.Contains(lowerPrompt, "laptop") {
		return "🤖 **Rekomendasi AI Gadget Assistant (Kategori Laptop)**:\n\n" +
			"1. **MacBook Pro M3 / Air M2**: Pilihan terbaik untuk produktivitas, pemograman, dan desain grafis dengan daya tahan baterai hingga 18 jam.\n" +
			"2. **ASUS ROG Zephyrus G16**: Pilihan gaming & video editing berat dengan GPU RTX dedicated dan layar 165Hz.\n\n" +
			"💡 *Tips*: Cek katalog gadget kami di `/api/v1/products?category=laptops` untuk melihat ketersediaan stok!"
	}

	if strings.Contains(lowerPrompt, "hp") || strings.Contains(lowerPrompt, "smartphone") || strings.Contains(lowerPrompt, "iphone") || strings.Contains(lowerPrompt, "samsung") {
		return "🤖 **Rekomendasi AI Gadget Assistant (Kategori Smartphone)**:\n\n" +
			"1. **iPhone 15 Pro Max**: Terbaik untuk videografi, performa chip A17 Pro, dan ketahanan bodi titanium.\n" +
			"2. **Samsung Galaxy S24 Ultra**: Terbaik untuk fotografi zoom 100x, layar Dynamic AMOLED 2X, dan fitur Galaxy AI bawaan.\n\n" +
			"💡 *Tips*: Cek katalog gadget kami di `/api/v1/products?category=smartphones` untuk diskon dan promo menarik!"
	}

	return "🤖 **Rekomendasi AI Gadget Assistant**:\n\n" +
		"Berdasarkan preferensi Anda, kami merekomendasikan untuk melihat katalog produk pilihan terbaik kami. " +
		"Kami menyediakan berbagai pilihan Smartphone, Laptop Gaming, Kamera Mirrorless, dan Aksesoris dengan harga bersaing!\n\n" +
		"💡 *Gunakan endpoint `/api/v1/products` untuk menjelajahi katalog lengkap kami.*"
}
