package result

type DomainResult struct {
	Result  any
	Message string
	Error   error
}

func NewResultApp(result any) DomainResult {
	return DomainResult{
		Result:  result,
		Error:   nil,
		Message: "",
	}
}

func NewErrorApp(err error) DomainResult {
	return DomainResult{
		Result:  nil,
		Error:   err,
		Message: "",
	}
}

func (r DomainResult) IsSuccess() bool {
	return r.Error == nil
}
