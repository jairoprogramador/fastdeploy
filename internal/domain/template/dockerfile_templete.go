package template

const DockerfileTemplate = `FROM eclipse-temurin:17-jre-alpine
LABEL description='{{.CommitMessage}}'
LABEL created.by='fastDeploy'
LABEL maintainer='{{.CommitAuthor}}'
LABEL commit='{{.CommitHash}}'
LABEL team='{{.Team}}'
LABEL organization='{{.Organization}}'

WORKDIR /app

COPY {{.FileName}} app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]`
