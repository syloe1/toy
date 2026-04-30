package model

type CalculateRequest struct {
	Left     float64 `json:"left"`
	Right    float64 `json:"right"`
	Operator string  `json:"operator"`
}
type CalculateResponse struct {
	Result float64 `json:"result"`
}
type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}
