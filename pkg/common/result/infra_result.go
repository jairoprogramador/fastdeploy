package result

type InfraResult struct {
	Result any
	Error  error
}

func NewResult(result any) InfraResult {
	return InfraResult{
		Result: result,
		Error:  nil,
	}
}

func NewError(err error) InfraResult {
	return InfraResult{
		Result: nil,
		Error:  err,
	}
}

func (r InfraResult) IsSuccess() bool {
	return r.Error == nil
}
