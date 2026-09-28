package teachsimpleclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mindfulBot/models"
	"net/http"
)

// ListBots - возвращает список ботов из simpleteach.
func (c *Client) ListBots(ctx context.Context) ([]models.Bot, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, http.NoBody)
	if err != nil {
		log.Println("simpleteach request failed")

		return []models.Bot{}, fmt.Errorf("simpleteach request failed")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("simpleteach do request failed")

		return []models.Bot{}, fmt.Errorf("simpleteach do request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("simpleteach unexpected status code: %d\n", resp.StatusCode)

		return []models.Bot{}, fmt.Errorf("simpleteach unexpected status code: %d\n", resp.StatusCode)
	}

	var dto []BotResponse
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		log.Println("simpleteach: decode response: ", err)

		return nil, fmt.Errorf("simpleteach: decode response: %w", err)
	}

	return toBots(dto), nil
}
