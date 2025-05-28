package result

type DomainResult struct {
	Message string
	Error   error
}

func NewMessageApp(message string) DomainResult {
	return DomainResult{
		Error:   nil,
		Message: message,
	}
}

func NewErrorApp(err error) DomainResult {
	return DomainResult{
		Error:   err,
		Message: "",
	}
}

func (r DomainResult) IsSuccess() bool {
	return r.Error == nil
}
