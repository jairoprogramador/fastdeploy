package dto

import "deploy/internal/domain/model"

type ResponseDto struct {
    Message string
    Error   error
}

func GetDtoWithModel(response *model.Response) *ResponseDto {
	return &ResponseDto {
		Message: response.Message,
		Error: response.Error,
	}
}

func GetDtoWithMessage(message string) *ResponseDto {
	return &ResponseDto {
		Message: message,
		Error: nil,
	}
}

func GetDtoWithError(err error) *ResponseDto {
	return &ResponseDto {
		Message: "",
		Error: err,
	}
}
