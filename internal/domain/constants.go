package constants

const RootDirectory   = ".deploy"
const NameProjectFile = "deploy.yml"

const DockerfileTemplateFilePath    = RootDirectory + "/docker/Dockerfile.template"
const DockerfileFilePath    = RootDirectory + "/docker/Dockerfile"
const DockercomposeTemplateFilePath = RootDirectory + "/docker/Dockercompose.template"
const DockercomposeFilePath = RootDirectory + "/docker/compose.yaml"

const MessageRunInit = "you must initialize the project, run the command: deploy init" 

const MessagePreviouslyInitializedProject = "the project is already initialized"
const MessageErrorInitializingProject = "error initializing project" 
const MessageSuccessInitializingProject = "project initialized successfully"

const MessageThereAreUnconfirmedChanges = "debe confirmar todos sus cambios antes de publicar"

const MessageErrorNoPackedFile = "error no packed files found"
const MessageErrorFileNotFound = "error file found"
const MessageErrorCreatingContainer = "error creating container"
const MessageErrorNoPortHost = "error no host port found"
const MessageSuccessPublish = "service %d available in: http://localhost:%s/"

const FileNameKey             = "fileName"
const CommitMessageKey        = "commitMessage"
const CommitHashKey           = "commitHash"
const CommitAuthorKey         = "commitAuthor"
const NameDeliveryKey         = "nameDelivery"
const ImageKey                = "commitHash"
const PortKey                 = "port"
const TeamKey                 = "Team"
const OrganizationKey         = "organization"
