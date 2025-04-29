package service

import (
	"deploy/internal/domain/repository"
    "deploy/internal/domain/model"
	"sync"
)

type PublishService struct {
    publishRepo repository.PublishRepository
}

var (
    instancePublishService     *PublishService
    instanceOncePublishService sync.Once
    mutexPublishService        sync.Mutex
)

func GetPublishService(publishRepo repository.PublishRepository) *PublishService {
    instanceOncePublishService.Do(func() {
        instancePublishService = &PublishService{
            publishRepo: publishRepo,
        }
    })
    return instancePublishService
}

func (s *PublishService) SetPublishService(publishRepo repository.PublishRepository) {
    mutexPublishService.Lock()
    defer mutexPublishService.Unlock()
    
    s.publishRepo = publishRepo
}

func (s *PublishService) Build() *model.Response {
	if resp := s.publishRepo.Prepare(); resp.Error != nil {
		return resp
	}
	return s.publishRepo.Build()
}

func (s *PublishService) Package(response *model.Response) *model.Response {
	return s.publishRepo.Package(response)
}

func (s *PublishService) Deliver(response *model.Response) *model.Response {
    if resp := s.publishRepo.Deliver(response); resp.Error != nil {
		return resp
	}
	return s.publishRepo.Validate(response)
}
