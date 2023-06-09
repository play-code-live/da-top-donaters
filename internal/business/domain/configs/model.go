package configs

type Config struct {
	ChannelID     string   `json:"channel_id"`
	Title         string   `json:"title,omitempty"`
	DonatersCount int      `json:"donaters_count,omitempty"`
	NamesToIgnore []string `json:"names_to_ignore"`
	TopCount      int      `json:"top_count"`
}
