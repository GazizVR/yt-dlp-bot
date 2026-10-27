package domain

type Action string

const SendAudio Action = "sendaudio"
const RetryAudio Action = "retryaudio"
const RetryVideo Action = "retryvideo"

type User struct {
	ID        int64
	Videopath string
	MediaURL  string
	Action    Action
}
