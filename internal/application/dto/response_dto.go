package dto

import "deploy/internal/domain/model"

type ResponseDto struct {
    Message string
    Error   error
}

func GetNewResponseDtoFromModel(response *model.Response) *ResponseDto {
	return &ResponseDto {
		Message: response.Message,
		Error: response.Error,
	}
}

func GetNewResponseDto(message string) *ResponseDto {
	return &ResponseDto {
		Message: message,
		Error: nil,
	}
}
