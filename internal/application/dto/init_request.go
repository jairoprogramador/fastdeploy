package dto

import "context"

type InitRequest struct {
	Ctx              context.Context
	SkipPrompt       bool
	WorkingDirectory string
}
