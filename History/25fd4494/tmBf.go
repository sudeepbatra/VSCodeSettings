package nse

type CorporateAction struct {
	Symbol  string `json:"symbol"`
	Series  string `json:"series"`
	Ind     string `json:"ind"`
	FaceVal string `json:"faceVal"`
	Subject string `json:"subject"`
	ExDate  string `json:"exDate"`
	RecDate string `json:"recDate"`
	Comp    string `json:"comp"`
	Isin    string `json:"isin"`
}
