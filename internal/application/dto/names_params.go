package dto

type NamesParams struct {
	projectName    string
	repositoryName string
}

func NewNamesParams(projectName, repositoryName string) NamesParams {
	return NamesParams{
		projectName:    projectName,
		repositoryName: repositoryName,
	}
}

func (r *NamesParams) ProjectName() string {
	return r.projectName
}

func (r *NamesParams) RepositoryName() string {
	return r.repositoryName
}
