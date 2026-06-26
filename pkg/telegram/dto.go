package telegram

type CommonResp struct {
	Ok bool `json:"ok"`
}

type GetUpdatesResp struct {
	CommonResp
	Result []Update `json:"result"`
}

type SendMessageResp struct {
	CommonResp
	Result Message `json:"result"`
}

type ErrorResponse struct {
	CommonResp
	Code        uint16 `json:"error_code"`
	Description string `json:"description"`
}
