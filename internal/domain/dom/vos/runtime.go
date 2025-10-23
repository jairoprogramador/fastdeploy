package vos

const (
	DefaultImageSource = "fastdeploy/runner-java17-springboot"
	DefaultImageTag    = "latest"
)

const (
	DefaultProjectMountPath = "/home/fastdeploy/app"
	DefaultStateMountPath   = "/home/fastdeploy/.fastdeploy"
)

type Image struct {
	source string
	tag    string
}

func NewImage(source, tag string) Image {
	if tag == "" {
		tag = DefaultImageTag
	}
	if source == "" {
		source = DefaultImageSource
	}
	return Image{source: source, tag: tag}
}

func (i Image) Source() string { return i.source }
func (i Image) Tag() string { return i.tag }

type Volumes struct {
	projectMountPath string
	stateMountPath   string
}

func NewVolumes(projectMountPath, stateMountPath string) Volumes {
	if projectMountPath == "" {
		projectMountPath = DefaultProjectMountPath
	}
	if stateMountPath == "" {
		stateMountPath = DefaultStateMountPath
	}
	return Volumes{projectMountPath: projectMountPath, stateMountPath: stateMountPath}
}

func (v Volumes) ProjectMountPath() string { return v.projectMountPath }
func (v Volumes) StateMountPath() string { return v.stateMountPath }


type Runtime struct {
	image   Image
	volumes Volumes
}

func NewRuntime(image Image, volumes Volumes) Runtime{
	return Runtime{
		image: image,
		volumes: volumes,
	}
}

func (r Runtime) Image() Image { return r.image }
func (r Runtime) Volumes() Volumes { return r.volumes }