package dto

// DailySmokingDTO store information about daily smoking
type DailySmokingDTO struct {
	Date   string   `json:"date"`   // date of the smoking, like 2024-03-17
	Count  int      `json:"count"`  // number of smoking for today
	Events []string `json:"events"` // list of smoking events on hour and minute, e.g ["08:00", "08:30"]
}
