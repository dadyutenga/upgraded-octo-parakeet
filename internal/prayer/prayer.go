package prayer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// PrayerTime holds the name and time for a single prayer.
type PrayerTime struct {
	Name string
	Time time.Time
}

// Aladhan API response structures
type aladhanResponse struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Data   aladhanData  `json:"data"`
}

type aladhanData struct {
	Timings aladhanTimings `json:"timings"`
	Date    aladhanDate    `json:"date"`
	Meta    aladhanMeta    `json:"meta"`
}

type aladhanTimings struct {
	Fajr    string `json:"Fajr"`
	Sunrise string `json:"Sunrise"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type aladhanDate struct {
	Readable string `json:"readable"`
}

type aladhanMeta struct {
	Method struct {
		Name string `json:"name"`
	} `json:"method"`
}

// cache stores fetched prayer times to avoid repeated API calls
var (
	cache      []PrayerTime
	cacheDate  string
	cacheMu    sync.Mutex
	cacheErr   error
	lastStatus string
)

// GetPrayerTimes fetches prayer times from the Aladhan API for Dar es Salaam
// using MWL (Muslim World League) method (method=3).
// Results are cached per day to minimize API calls.
func GetPrayerTimes(date time.Time) ([]PrayerTime, error) {
	dateStr := date.Format("02-01-2006")

	cacheMu.Lock()
	if cacheDate == dateStr && cache != nil {
		result := make([]PrayerTime, len(cache))
		copy(result, cache)
		cacheMu.Unlock()
		return result, nil
	}
	cacheMu.Unlock()

	// Aladhan API: method=3 is MWL (Muslim World League)
	url := fmt.Sprintf(
		"http://api.aladhan.com/v1/timingsByCity/%s?city=Dar%%20es%%20Salaam&country=Tanzania&method=3",
		dateStr,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResp aladhanResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API error: %s", apiResp.Status)
	}

	loc, _ := time.LoadLocation("Africa/Dar_es_Salaam")
	if loc == nil {
		loc = time.FixedZone("EAT", 3*60*60)
	}

	t := apiResp.Data.Timings
	prayers := []PrayerTime{
		{Name: "Fajr", Time: parseTime(t.Fajr, date, loc)},
		{Name: "Sunrise", Time: parseTime(t.Sunrise, date, loc)},
		{Name: "Dhuhr", Time: parseTime(t.Dhuhr, date, loc)},
		{Name: "Asr", Time: parseTime(t.Asr, date, loc)},
		{Name: "Maghrib", Time: parseTime(t.Maghrib, date, loc)},
		{Name: "Isha", Time: parseTime(t.Isha, date, loc)},
	}

	cacheMu.Lock()
	cache = prayers
	cacheDate = dateStr
	cacheErr = nil
	lastStatus = fmt.Sprintf("MWL (%s)", apiResp.Data.Meta.Method.Name)
	cacheMu.Unlock()

	return prayers, nil
}

// GetLastStatus returns the method info from the last successful fetch.
func GetLastStatus() string {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	return lastStatus
}

// parseTime parses "HH:MM" or "HH:MM (EAT)" format from the API response.
func parseTime(s string, date time.Time, loc *time.Location) time.Time {
	// API may return "HH:MM (TZ)" — strip the timezone suffix
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, " "); idx != -1 {
		s = s[:idx]
	}

	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return date
	}

	var h, m int
	fmt.Sscanf(parts[0], "%d", &h)
	fmt.Sscanf(parts[1], "%d", &m)

	year, month, day := date.Date()
	return time.Date(year, month, day, h, m, 0, 0, loc)
}

// Render returns a formatted string of prayer times for display.
func Render(prayers []PrayerTime, now time.Time, fetchErr error) string {
	var b strings.Builder

	b.WriteString("\033[1m\033[36m╔══════════════════════════════════════╗\033[0m\n")
	b.WriteString("\033[1m\033[36m║   🕌  Prayer Times - Dar es Salaam  ║\033[0m\n")
	b.WriteString("\033[1m\033[36m╚══════════════════════════════════════╝\033[0m\n")

	status := GetLastStatus()
	if status != "" {
		b.WriteString(fmt.Sprintf("  \033[90mMethod: %s\033[0m\n", status))
	}
	b.WriteString("\n")

	if fetchErr != nil {
		b.WriteString(fmt.Sprintf("  \033[31m⚠ %s\033[0m\n", fetchErr.Error()))
		b.WriteString("  \033[90mCheck internet connection and retry.\033[0m\n")
		return b.String()
	}

	if len(prayers) == 0 {
		b.WriteString("  \033[33mLoading prayer times...\033[0m\n")
		return b.String()
	}

	nextFound := false
	for _, p := range prayers {
		marker := "  "
		color := "\033[0m"
		if !nextFound && p.Time.After(now) {
			marker = "▶ "
			color = "\033[33m\033[1m"
			nextFound = true
		} else if p.Time.Before(now) {
			color = "\033[90m"
		}
		b.WriteString(fmt.Sprintf("  %s%s%-10s %s%s\n",
			color, marker, p.Name, p.Time.Format("15:04"), "\033[0m"))
	}

	if !nextFound {
		b.WriteString("\n  \033[32mAll prayers completed for today ✓\033[0m\n")
	}

	return b.String()
}
