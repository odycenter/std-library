package dto

type ErrorMsg struct {
	Code int `json:"code"`
	*ErrorData
}

type ErrorData struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
