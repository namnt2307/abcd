package schema

type ArtistContentObjOutputStruct struct {
	Seo struct {
		Url         string `json:"url" `
		Description string `json:"decription" `
		Title       string `json:"title" `
	} `json:"seo" `

	Items []struct {
		Id                string                   `json:"id" `
		Group_id          string                   `json:"group_id" `
		Type              int                      `json:"type" `
		Title             string                   `json:"title" `
		Known_as          string                   `json:"known_as" `
		Status            int                      `json:"status" `
		Is_new_update     bool                     `json:"is_new_update" `
		Is_premium        int                      `json:"is_premium" `
		Is_new            int                      `json:"is_new" `
		Short_description string                   `json:"short_description" `
		Long_description  string                   `json:"long_description" `
		Caption_langs     string                   `json:"caption_langs" `
		Caption           string                   `json:"caption" `
		Resolution        int                      `json:"resolution" `
		Runtime           int                      `json:"runtime" `
		Slug              string                   `json:"slug" `
		Slug_seo          string                   `json:"slug_seo" `
		Released          bool                     `json:"released" `
		Avg_rate          float64                  `json:"avg_rate" `
		Total_rate        int                      `json:"total_rate" `
		Release_year      int                      `json:"release_year" `
		Publish_date      string                   `json:"publish_date" `
		Release_date      string                   `json:"release_date" `
		Images            ImagesOutputObjectStruct `json:"images" `
		Is_watchlater     bool                     `json:"is_watchlater" `
		Episode           int                      `json:"episode" `
		Current_episode   string                   `json:"current_episode" `
		Seo               struct {
			Url         string `json:"url" `
			Description string `json:"decription" `
			Title       string `json:"title" `
		} `json:"seo" `
		// Intro_start  int64 `json:"intro_start" `
		// Intro_end    int64 `json:"intro_end" `
		// Outtro_start int64 `json:"outtro_start" `
		// Outtro_end   int64 `json:"outtro_end" `
	}

	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
}
