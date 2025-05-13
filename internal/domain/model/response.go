package model

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


