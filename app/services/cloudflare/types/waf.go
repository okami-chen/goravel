package types

type Rule struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Action      string `json:"action"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	LastUpdated string `json:"last_updated"`
	Ref         string `json:"ref"`
	Enabled     bool   `json:"enabled"`
}
