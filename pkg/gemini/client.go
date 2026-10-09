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

	"gadget-marketplace/models"
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

func (c *GeminiClient) GenerateGadgetRecommendation(userPrompt string, catalog []models.Product) (string, error) {
	catalogSummary := formatCatalogForAI(catalog)

	if c.apiKey == "" {
		return c.generateFallbackRecommendation(userPrompt, catalog), nil
	}

	systemContext := fmt.Sprintf(`Anda adalah Asisten AI Spesialis Gadget Marketplace yang cerdas, sopan, dan berpengalaman.

Berikut adalah KATALOG PRODUK REALS YANG TERSEDIA DAN READY STOCK DI TOKO KAMI SAAT INI:
%s

Instruksi Penting:
1. Analisis kebutuhan dan budget pengguna dari pertanyaan mereka.
2. REKOMENDASIKAN PRODUK YANG TERSEDIA DI KATALOG TOKO KAMI DI ATAS!
3. Sebutkan Product ID, Nama Produk, Harga (dalam Rp), dan alasan mengapa produk tersebut cocok.
4. Beritahu pengguna bahwa mereka bisa langsung memesan produk tersebut di marketplace dengan product_id terkait.
5. Jawab dengan bahasa Indonesia yang ramah, jelas, dan informatif.

Pertanyaan Pengguna: %s`, catalogSummary, userPrompt)

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
		log.Printf("Warning: Gemini API connection error (%v). Serving catalog AI fallback.", err)
		return c.generateFallbackRecommendation(userPrompt, catalog), nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Warning: Gemini API returned status %d. Serving catalog AI recommendation fallback.", resp.StatusCode)
		return c.generateFallbackRecommendation(userPrompt, catalog), nil
	}

	var geminiResp GeminiGenerateResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return c.generateFallbackRecommendation(userPrompt, catalog), nil
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return c.generateFallbackRecommendation(userPrompt, catalog), nil
}

func formatCatalogForAI(catalog []models.Product) string {
	if len(catalog) == 0 {
		return "- Katalog toko saat ini belum di-sync."
	}

	var sb strings.Builder
	for _, p := range catalog {
		sb.WriteString(fmt.Sprintf("- Product ID: %d | Nama: %s | Kategori: %s | Brand: %s | Harga: Rp %.0f | Stok: %d unit\n",
			p.ID, p.Name, p.Category, p.Brand, p.Price, p.Stock))
	}
	return sb.String()
}

func (c *GeminiClient) generateFallbackRecommendation(userPrompt string, catalog []models.Product) string {
	if len(catalog) == 0 {
		return "🤖 **Rekomendasi AI Gadget Assistant**:\n\nKatalog toko saat ini masih kosong. Silakan jalankan `POST /api/v1/products/sync` terlebih dahulu untuk mengisi katalog produk toko kami!"
	}

	lowerPrompt := strings.ToLower(userPrompt)
	var recommended []models.Product

	for _, p := range catalog {
		if strings.Contains(lowerPrompt, strings.ToLower(p.Category)) ||
			strings.Contains(lowerPrompt, strings.ToLower(p.Name)) ||
			strings.Contains(lowerPrompt, strings.ToLower(p.Brand)) {
			recommended = append(recommended, p)
		}
	}

	if len(recommended) == 0 && len(catalog) > 0 {
		count := 3
		if len(catalog) < 3 {
			count = len(catalog)
		}
		recommended = catalog[:count]
	}

	var sb strings.Builder
	sb.WriteString("🤖 **Rekomendasi AI Gadget Assistant (Katalog Produk Real Toko Kami)**:\n\n")
	for i, p := range recommended {
		sb.WriteString(fmt.Sprintf("%d. **%s** (Brand: %s)\n", i+1, p.Name, p.Brand))
		sb.WriteString(fmt.Sprintf("   - **Product ID**: %d\n", p.ID))
		sb.WriteString(fmt.Sprintf("   - **Harga Beli**: Rp %.2f\n", p.Price))
		sb.WriteString(fmt.Sprintf("   - **Stok Tersedia**: %d unit\n\n", p.Stock))
	}

	sb.WriteString("💡 *Anda dapat langsung membeli gadget di atas dengan memasukkan Product ID terkait ke endpoint `POST /api/v1/orders/checkout`!*")
	return sb.String()
}
