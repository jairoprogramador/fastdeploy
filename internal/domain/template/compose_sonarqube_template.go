package template

const ComposeSonarqubeTemplate = `services:
  sonarqube:
    image: sonarqube:community
    container_name: sonarqube-fastdeploy
    ports:
      - "9000:9000"
    volumes:
      - {{.HomeDir}}/.fastdeploy/sonarqube/volumes/data:/opt/sonarqube/data
      - {{.HomeDir}}/.fastdeploy/sonarqube/volumes/logs:/opt/sonarqube/logs
      - {{.HomeDir}}/.fastdeploy/sonarqube/volumes/extensions:/opt/sonarqube/extensions`