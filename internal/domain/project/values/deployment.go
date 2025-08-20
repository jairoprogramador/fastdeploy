package values

type Deployment struct {
	version VersionDeployment
}

func NewDeployment(version VersionDeployment) Deployment {
	return Deployment{
		version: version,
	}
}

func NewDefaultDeployment() Deployment {
	return NewDeployment(NewDefaultVersionDeployment())
}

func (d Deployment) GetVersion() VersionDeployment {
	return d.version
}

func (d Deployment) IncrementVersion() Deployment {
	newVersion := d.version.Increment()
	return Deployment{
		version: newVersion,
	}
}

func (d Deployment) Equals(other Deployment) bool {
	return d.version.Value() == other.version.Value()
}


