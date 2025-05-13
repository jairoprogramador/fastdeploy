package template

const ComposeTemplate = `services:
  {{.NameDelivery}}:
    image: {{.CommitHash}}:{{.Version}}
    container_name: {{.NameDelivery}}
    ports:
      - "{{.Port}}:8080"
    restart: unless-stopped`