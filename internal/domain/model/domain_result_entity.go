package model

type DomainResultEntity struct {
	Result  any
	Message string
	Error   error
}

func NewResultApp(result any) DomainResultEntity {
	return DomainResultEntity{
		Result:  result,
		Error:   nil,
		Message: "",
	}
}

func NewErrorApp(err error) DomainResultEntity {
	return DomainResultEntity{
		Result:  nil,
		Error:   err,
		Message: "",
	}
}

func (r DomainResultEntity) IsSuccess() bool {
	return r.Error == nil
}
