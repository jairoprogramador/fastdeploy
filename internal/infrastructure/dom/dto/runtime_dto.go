package dto

type RuntimeDTO struct {
	Image ImageDTO `yaml:"image"`
	Volumes VolumesDTO `yaml:"volumes"`
}

type ImageDTO struct {
	Source string `yaml:"source"`
	Tag    string `yaml:"tag"`
}

type VolumesDTO struct {
	ProjectMountPath string `yaml:"project_mount_path"`
	StateMountPath   string `yaml:"state_mount_path"`
}