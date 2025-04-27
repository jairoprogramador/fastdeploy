package model

import "deploy/internal/domain"

type Response struct {
    Message string
    Error   error
    Data    map[string]string
}

func GetNewResponseError(err error) *Response {
	return &Response {
		Message:    "",
		Error:     err,
		Data: make(map[string]string),
	}
}

func GetNewResponseMessage(message string) *Response {
	return &Response {
		Message:    message,
		Error:     nil,
		Data: make(map[string]string),
	}
}

func GetNewResponse() *Response {
	return &Response {
		Message:    "",
		Error:     nil,
		Data: make(map[string]string, 5),
	}
}

func (s *Response) SetCommitHash(commitHash string) {
	s.Data[constants.CommitHashKey] = commitHash
}

func (s *Response) SetImageId(imageId string) {
	s.Data[constants.ImageKey] = imageId
}

func (s *Response) GetCommitHash() string {
	return s.Data[constants.CommitHashKey]
}

func (s *Response) GetImageId() string {
	return s.Data[constants.ImageKey]
}

