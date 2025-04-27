package service

import "deploy/internal/domain"
import "deploy/internal/domain/model"
import "sync"


type InitializeService struct {
    projectService ProjectService
}

var (
    instanceInitializeService     *InitializeService
    instanceOnceInitializeService sync.Once
    mutexInitializeService        sync.Mutex
)

func GetInitializeService(projectService ProjectService) *InitializeService {
    instanceOnceInitializeService.Do(func() {
        instanceInitializeService = &InitializeService{
            projectService: projectService,
        }
    })
    return instanceInitializeService
}

func (s *InitializeService) SetProjectRepository(projectService ProjectService) {
    mutexInitializeService.Lock()
    defer mutexInitializeService.Unlock()
    
    s.projectService = projectService
}

func (s *InitializeService) Initialize() *model.Response {
	if s.projectService.Exists() {
		return model.GetNewResponseMessage(constants.MessagePreviouslyInitializedProject)
	}

	if resp := s.projectService.Create(); resp.Error != nil {
        return resp
    }
	return model.GetNewResponseMessage(constants.MessageSuccessInitializingProject)
}

