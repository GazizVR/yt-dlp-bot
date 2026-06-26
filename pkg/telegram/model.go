package telegram

type Update struct {
	Id      int64   `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id          int64          `json:"message_id"`
	Text        string         `json:"text"`
	LinkPreview LinkPreviewOps `json:"link_preview_options"`
	Chat        Chat           `json:"chat"`
}

type LinkPreviewOps struct {
	URL string `json:"url"`
}

type Chat struct {
	Id int64 `json:"id"`
}
