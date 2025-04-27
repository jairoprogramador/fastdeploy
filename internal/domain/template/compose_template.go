package template

const ComposeTemplate = `services:
  {{.NameDelivery}}:
    image: {{.CommitHash}}
    container_name: {{.NameDelivery}}
    ports:
      - "{{.Port}}:8080"
    restart: unless-stopped`