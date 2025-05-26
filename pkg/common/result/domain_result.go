package result

// DomainResultEntity represents the result of a domain operation
type DomainResultEntity struct {
	Result  any
	Message string
	Error   error
}

// NewResultApp creates a new successful DomainResultEntity
func NewResultApp(result any) DomainResultEntity {
	return DomainResultEntity{
		Result:  result,
		Error:   nil,
		Message: "",
	}
}

// NewErrorApp creates a new error DomainResultEntity
func NewErrorApp(err error) DomainResultEntity {
	return DomainResultEntity{
		Result:  nil,
		Error:   err,
		Message: "",
	}
}

// IsSuccess checks if the operation was successful
func (r DomainResultEntity) IsSuccess() bool {
	return r.Error == nil
}