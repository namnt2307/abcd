package topviews

type ViewsContentStruct struct {
	Id_content string `json:"id_content"`
	Views      int64  `json:"views"`
	Date       string `json:"date"`
	Ranking    int    `json:"ranking"`
}
