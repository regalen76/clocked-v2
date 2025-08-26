package handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Holiday struct {
	Summary string `json:"summary"`
}
type HolidaysMap map[string]Holiday

type HolidaysWithInfo struct {
	Dates HolidaysMap
	Info  map[string]any
}

func fetchHolidays(ctx context.Context) (*HolidaysWithInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://raw.githubusercontent.com/guangrei/APIHariLibur_V2/main/holidays.json", nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("failed to fetch holidays: " + res.Status)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := &HolidaysWithInfo{Dates: make(HolidaysMap)}
	for k, v := range raw {
		if k == "info" {
			var info map[string]any
			_ = json.Unmarshal(v, &info)
			out.Info = info
			continue
		}
		var h Holiday
		if err := json.Unmarshal(v, &h); err == nil {
			out.Dates[k] = h
		}
	}
	return out, nil
}

// Fiber handler: GET /api/holidays
func GetHolidays(c *fiber.Ctx) error {
	data, err := fetchHolidays(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			// Keep shape similar to your load() error path:
			"holidays":    fiber.Map{}, // empty map
			"attendances": []any{},     // not used yet
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Fetched Holidays", "data": data.Dates})
}
