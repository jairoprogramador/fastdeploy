package result

// InfraResultEntity represents the result of an infrastructure operation
type InfraResultEntity struct {
	Result any
	Error  error
}

// NewResult creates a new successful InfraResultEntity
func NewResult(result any) InfraResultEntity {
	return InfraResultEntity{
		Result: result,
		Error:  nil,
	}
}

// NewError creates a new error InfraResultEntity
func NewError(err error) InfraResultEntity {
	return InfraResultEntity{
		Result: nil,
		Error:  err,
	}
}

// IsSuccess checks if the operation was successful
func (r InfraResultEntity) IsSuccess() bool {
	return r.Error == nil
}