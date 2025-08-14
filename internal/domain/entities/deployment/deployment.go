package deployment

type Deployment struct {
	version DeploymentVersion
}

func NewDeployment(version DeploymentVersion) Deployment {
	return Deployment{
		version: version,
	}
}

func (d Deployment) GetVersion() DeploymentVersion {
	return d.version
}

func (d Deployment) IncrementVersion() Deployment {
	newVersion := d.version.Increment()
	return Deployment{
		version: newVersion,
	}
}

func (d Deployment) IsValid() bool {
	return d.version.Value() != ""
}

func (d Deployment) Equals(other Deployment) bool {
	return d.version.Value() == other.version.Value()
}

func (d Deployment) String() string {
	return d.version.String()
}
