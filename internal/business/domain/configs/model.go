package configs

type Config struct {
	ChannelID     string `json:"channel_id"`
	Title         string `json:"title,omitempty"`
	DonatersCount int    `json:"donaters_count,omitempty"`
}
