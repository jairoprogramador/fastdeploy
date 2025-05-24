package model

// InfrastructureResponse represents a standardized response from infrastructure implementations
// It contains a result of any type, an error if one occurred, and additional details if needed
type InfrastructureResponse struct {
	// Result is the actual result of the operation, which can be of any type
	Result any
	// Error is the error that occurred during the operation, if any
	Error error
	// Details contains additional information about the operation
	Details string
}

// NewResponse creates a new InfrastructureResponse with the given result and no error
func NewResponse(result any) InfrastructureResponse {
	return InfrastructureResponse{
		Result: result,
		Error:  nil,
	}
}

// NewResponseWithDetails creates a new InfrastructureResponse with the given result and details
func NewResponseWithDetails(result any, details string) InfrastructureResponse {
	return InfrastructureResponse{
		Result:  result,
		Error:   nil,
		Details: details,
	}
}

// NewErrorResponse creates a new InfrastructureResponse with the given error
func NewErrorResponse(err error) InfrastructureResponse {
	return InfrastructureResponse{
		Result: nil,
		Error:  err,
	}
}

// NewErrorResponseWithDetails creates a new InfrastructureResponse with the given error and details
func NewErrorResponseWithDetails(err error, details string) InfrastructureResponse {
	return InfrastructureResponse{
		Result:  nil,
		Error:   err,
		Details: details,
	}
}

// IsSuccess returns true if the response does not contain an error
func (r InfrastructureResponse) IsSuccess() bool {
	return r.Error == nil
}