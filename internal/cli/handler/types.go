package handler

import "deploy/internal/domain/model"

type IsInitAppFunc func() (*model.ProjectEntity, error)
