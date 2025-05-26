package model

type InfraResultEntity struct {
	Result any
	Error  error
}

func NewResult(result any) InfraResultEntity {
	return InfraResultEntity{
		Result: result,
		Error:  nil,
	}
}

func NewError(err error) InfraResultEntity {
	return InfraResultEntity{
		Result: nil,
		Error:  err,
	}
}

func (r InfraResultEntity) IsSuccess() bool {
	return r.Error == nil
}
