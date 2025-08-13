package project

type ProjectInitialize interface {
	Initialize() (*ProjectEntity, error)
	IsInitialized() bool
}
