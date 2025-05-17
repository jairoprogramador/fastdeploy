package template

const ComposeTemplate = `services:
  {{.NameDelivery}}:
    build:
      context: {{.PathHomeDirectory}}
      dockerfile: {{.PathDockerDirectory}}/Dockerfile
    image: {{.CommitHash}}:{{.Version}}
    container_name: {{.NameDelivery}}
    ports:
      - "{{.Port}}:8080"
    restart: unless-stopped`