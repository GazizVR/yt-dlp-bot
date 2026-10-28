package bot

import "errors"

var ErrEmptyCallbackData = errors.New("callback data must not be empty")
var ErrInvalidCallbackData = errors.New("callback data is invalid")
var ErrMessageNotFound = errors.New("callback query message not found")
