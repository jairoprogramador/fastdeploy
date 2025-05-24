package port
/* 
import (
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/repository"
	"sync"
)

type SonarqubeService struct {
	sonarqubeRepo repository.SonarqubeRepository
}

var (
	instanceSonarqubeService     *SonarqubeService
	instanceOnceSonarqubeService sync.Once
)

func GetSonarqubeService(sonarqubeRepo repository.SonarqubeRepository) *SonarqubeService {
	instanceOnceSonarqubeService.Do(func() {
		instanceSonarqubeService = &SonarqubeService{
			sonarqubeRepo: sonarqubeRepo,
		}
	})
	return instanceSonarqubeService
}

func (s *SonarqubeService) Add() *logger.Logger {
	return s.sonarqubeRepo.Add()
}
 */