package users

type SessionData struct {
	Shop      string          `json:"shop"`
	GuideStep map[string]bool `json:"guide_step"`
	GuideShow bool            `json:"guide_show"`
}
