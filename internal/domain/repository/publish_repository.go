package repository

import "deploy/internal/domain/model"

type PublishRepository interface {
	Prepare() *model.Response
	Build() *model.Response
	Package(*model.Response) *model.Response
	Deliver(*model.Response) *model.Response
	Validate(*model.Response) *model.Response
}
