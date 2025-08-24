package shopifys

type OnlineStoreTheme struct {
	ID            string `json:"id"`
	Role          string `json:"role"`
	CreatedAt     string `json:"createdAt"`
	Name          string `json:"name"`
	Prefix        string `json:"prefix"`
	Processing    bool   `json:"processing"`
	ProcessFailed bool   `json:"processFailed"`
	ThemeStoreID  string `json:"themeStoreId"`
	UpdatedAt     string `json:"updatedAt"`
	Files         struct {
		Nodes []struct {
			OnlineStoreThemeFile
		} `json:"nodes"`
	} `json:"files"`
}
type OnlineStoreThemeFile struct {
	Body struct {
		OnlineStoreThemeFileBodyUrl
		OnlineStoreThemeFileBodyBase64
		OnlineStoreThemeFileBoyText
	} `json:"body"`
	ChecksumMd5 string `json:"checksumMd5"`
	ContentType string `json:"contentType"`
	CreatedAt   string `json:"createdAt"`
	Filename    string `json:"filename"`
	Size        string `json:"size"`
	UpdatedAt   string `json:"updatedAt"`
}

type OnlineStoreThemeFileBodyBase64 struct {
	ContentBase64 string `json:"contentBase64,omitempty"`
}
type OnlineStoreThemeFileBoyText struct {
	Content string `json:"content,omitempty"`
}
type OnlineStoreThemeFileBodyUrl struct {
	Url string `json:"url,omitempty"`
}
type SettingsData struct {
	Current struct {
		Sections map[string]struct {
			Type     string                 `json:"type"`
			Settings map[string]interface{} `json:"settings"`
		} `json:"sections"`
		Colors map[string]string       `json:"colors_solid_button_labels,colors_accent_1,colors_accent_2,colors_text,colors_outline_button_labels,colors_background_1,colors_background_2"`
		Blocks map[string]BlockSetting `json:"blocks"`
	} `json:"current"`
	Presets map[string]interface{} `json:"presets"`
}

type BlockSetting struct {
	Type     string                 `json:"type"`
	Disabled bool                   `json:"disabled"`
	Settings map[string]interface{} `json:"settings"`
}
