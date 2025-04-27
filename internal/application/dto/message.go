package dto

import "deploy/internal/domain/model"

type Message struct {
    Message string
    Error   error
}

func GetNewMessageFromResponse(response *model.Response) *Message {
	return &Message {
		Message: response.Message,
		Error: response.Error,
	}
}

func GetNewMessage(message string) *Message {
	return &Message {
		Message: message,
		Error: nil,
	}
}
