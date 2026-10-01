package teachsimpleclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// dayOrder is the canonical weekday ordering used to build Config.Days,
// independent of Go's randomized map iteration order.
var dayOrder = []string{"пн", "вт", "ср", "чт", "пт", "сб", "вс"}

type Config struct {
	BotID       int
	Name        string
	Days        []string
	Slots       map[string][]string
	Welcome     string
	Reminder1h  string
	Reminder24h string
	Paylink     string
}

type slotDTO struct {
	Day  string `json:"day"`
	Time string `json:"time"`
}

type messageDTO struct {
	WelcomeText     string `json:"welcome_text"`
	Reminder1hText  string `json:"reminder_1h_text"`
	Reminder24hText string `json:"reminder_24h_text"`
	Paylink         string `json:"paylink"`
}

type botConfigDTO struct {
	BotID   int        `json:"bot_id"`
	Name    string     `json:"name"`
	Slots   []slotDTO  `json:"slots"`
	Message messageDTO `json:"message"`
}

// Client is a plain per-instance struct (not a package-level singleton) so a
// future multi-bot dispatcher can hold one Client per bot token.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Fetch(ctx context.Context) (*Config, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/bot-config", nil)
	if err != nil {
		return nil, fmt.Errorf("proteachclient: build request: %w", err)
	}
	req.Header.Set("X-Bot-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("proteachclient: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proteachclient: unexpected status %d", resp.StatusCode)
	}

	var dto botConfigDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("proteachclient: decode response: %w", err)
	}

	return toConfig(dto), nil
}

func toConfig(dto botConfigDTO) *Config {
	slots := make(map[string][]string)
	for _, s := range dto.Slots {
		slots[s.Day] = append(slots[s.Day], s.Time)
	}

	days := make([]string, 0, len(slots))
	for _, day := range dayOrder {
		if _, ok := slots[day]; ok {
			days = append(days, day)
		}
	}

	return &Config{
		BotID:       dto.BotID,
		Name:        dto.Name,
		Days:        days,
		Slots:       slots,
		Welcome:     dto.Message.WelcomeText,
		Reminder1h:  dto.Message.Reminder1hText,
		Reminder24h: dto.Message.Reminder24hText,
		Paylink:     dto.Message.Paylink,
	}
}
