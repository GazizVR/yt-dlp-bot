package telegram

type CommonResponse struct {
	Ok bool `json:"ok"`
}

type UpdatesResponse struct {
	CommonResponse
	Result []Update `json:"result"`
}

type MessageResponse struct {
	CommonResponse
	Result Message `json:"result"`
}

type ErrorResponse struct {
	CommonResponse
	Code        uint16 `json:"error_code"`
	Description string `json:"description"`
}
