package service

import (
	"deploy/internal/domain/repository"
    "deploy/internal/domain/model"
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

func (s *SonarqubeService) Add() *model.Response {
    return s.sonarqubeRepo.Add()
}