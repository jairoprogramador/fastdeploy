<div align="center">
  <!-- <img src="doc/img/FastDeploy.jpg" alt="FastDeploy Logo" width="150"/> -->
  <h1>fastdeploy (fd)</h1>
  <p><strong>Despliega cualquier tecnología en cualquier plataforma con solo 3 comandos.</strong></p>
  <p><i>La infraestructura se convierte en una plantilla.</i></p>

  <p>
    <a href="https://github.com/jairoprogramador/fastdeploy/releases">
      <img src="https://img.shields.io/github/v/release/jairoprogramador/fastdeploy?style=for-the-badge" alt="Latest Release">
    </a>
    <a href="https://github.com/jairoprogramador/fastdeploy/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/jairoprogramador/fastdeploy?style=for-the-badge" alt="License">
    </a>
  </p>
</div>

---

**`fastdeploy` (o `fd`)** es una herramienta CLI diseñada para eliminar la complejidad de los procesos de despliegue. Olvídate de los scripts frágiles, los largos `READMEs` y la pregunta "¿cómo se desplegaba esto?". Con `fastdeploy`, estandarizas tus despliegues usando plantillas reutilizables, permitiendo que cualquier desarrollador, en cualquier equipo, pueda desplegar cualquier aplicación de forma segura y predecible.

**Define tu proceso de despliegue una vez, y ejecútalo miles de veces con dos simples comandos.**

## ✨ Características Principales

*   **⚙️ Agnostico a la Tecnología:** ¿Java, Node.js, Python, Go? ¿Terraform, Docker, Kubernetes? `fastdeploy` orquesta cualquier herramienta que puedas ejecutar en un shell.
*   **📄 Infraestructura como plantilla:** Centraliza la lógica de tus despliegues (steps, variables, entornos) en un repositorio de plantillas. Estandariza las buenas prácticas y evoluciona tu infraestructura sin tocar tus microservicios.
*   **🚀 Despliegues en dos pasos:** Clona tu microservicio y ejecuta `fd init`, y `fd deploy`. Eso es todo.
*   **✅ Verificación continua:** El estado de cada despliegue se guarda, permitiendo validaciones y evitando ejecuciones accidentales en entornos incorrectos.
*   **💻 Experiencia de desarrollador primero:** Comandos intuitivos, feedback claro y la abstracción perfecta para que los desarrolladores se centren en lo que importa: el código.

## 🚀 Instalación

Instala `fastdeploy` en segundos.

### macOS (Homebrew)

```sh
brew install jairoprogramador/fastdeploy/fastdeploy
```

### Linux

Puedes descargar el paquete `.deb` o `.rpm` desde la [página de Releases](https://github.com/jairoprogramador/fastdeploy/releases) y usar tu gestor de paquetes.

```sh
# Para sistemas basados en Debian/Ubuntu
sudo dpkg -i fastdeploy_*.deb

# Para sistemas basados en Red Hat/Fedora
sudo rpm -i fastdeploy_*.rpm
```

Alternativamente, puedes descargar el binario directamente:
```sh
curl -sL https://github.com/jairoprogramador/fastdeploy/releases/latest/download/fastdeploy_linux_amd64.tar.gz | tar xz
sudo mv fd /usr/local/bin/
```

### Windows

1.  Descarga el archivo `fastdeploy_windows_***64.zip` desde la [página de Releases](https://github.com/jairoprogramador/fastdeploy/releases).
2.  Descomprime el archivo.
3.  Añade el ejecutable `fd.exe` a tu variable de entorno `PATH`.


## 🏁 Guía de Inicio Rápido: Desplegando un Microservicio Java

Vamos a desplegar un microservicio Java que utiliza **Terraform** para provisionar la infraestructura en **Azure** y se empaqueta con **Docker**.

Toda la lógica de este despliegue está definida en nuestra plantilla de ejemplo:
➡️ **[jairoprogramador/mydeploy](https://github.com/jairoprogramador/mydeploy)**

Este repositorio de plantillas contiene los `steps`, `variables` y la definición de los `environments` (ej: `sandbox`, `stagin`, `produccion`).

### Paso 1: Inicializa tu Proyecto

Clona o crear el proyecto de microservicio que quieres desplegar. Una vez dentro del directorio del proyecto, ejecuta:

```sh
fd init
```

`fastdeploy` detectará que no está inicializado y te hará un par de preguntas para crear el archivo de configuración local `fdconfig.yaml`. Este archivo vincula tu proyecto con la plantilla de despliegue.

```yaml
# .fdconfig.yaml (Ejemplo generado)
project:
  id: 9238fa29be....
  name: "test"
  version: "1.0.0"
  team: "shikigami"
  description: "Mi proyecto de ejemplo"
  organization: "fastdeploy"

template:
  url: "https://github.com/jairoprogramador/mydeploy.git"
  ref: "main"
runtime:
    image: Dockerfile
    tag: latest
    build:
        args:
            - name: "FASTDEPLOY_VERSION"
              value: "1.0.10"
            - name: "MAVEN_VERSION"
              value: "3.9.12"
    run:
        volumes:
            - host: /home/user/.m2/
              container: /home/ubuntu/.m2
            - host: /home/user/myproject
              container: /home/fastdeploy/app
            - host: /home/user/dirFastDeploy
              container: /home/ubuntu/.fastdeploy
        envs:
            - name: "ARM_CLIENT_ID"
              value: "$ARM_CLIENT_ID"
            - name: "ARM_CLIENT_SECRET"
              value: "$ARM_CLIENT_SECRET"
            - name: "ARM_TENANT_ID"
              value: "$ARM_TENANT_ID"
            - name: "ARM_SUBSCRIPTION_ID"
              value: "$ARM_SUBSCRIPTION_ID"
```

### Paso 2: Prueba el despliegue en un entorno

Antes de desplegar, puedes validar que todo está bien. El comando `fd test [environment]` ejecuta los comandos definidos en la plantilla referentes a las pruebas.

```sh
# Ejecuta los pasos de prueba para el entorno 'sand'
fd test sand
```

Esto podría, por ejemplo, compilar el proyecto, ejecutar los test unitarios, las pruebas de seguridad, validar versiones, verificar pull request, etc, sin desplegarlo.

### Paso 3: Despliega

Una vez que las pruebas pasan, estás listo para desplegar. El comando `fd deploy [environment]` ejecuta la secuencia completa de pasos definidos en la plantilla, por ejemplo para el entorno de sandbox.

```sh
# Despliega en el entorno 'sand'
fd deploy sand
```
`fastdeploy` orquestará todo el proceso:
1.  Clonará la plantilla de despliegue.
2.  Ejecutará los pasos para aprovisionar recursos.
3.  Empaquetará y subirá la imagen del proyecto.
4.  Desplegará la aplicación en el ambiente elegido.

¡Y listo! Tu microservicio está desplegado.

## 📚 Comandos Básicos

| Comando | Descripción |
| :--- | :--- |
| `fd init` | Inicializa un proyecto creando el archivo `fdconfig.yaml`. |
| `fd [step] [env]` | Ejecuta hasta el `step` indicado en el entorno `env`. |
| `fd test [env]` | Ejecuta hasta el paso `test` en el entorno `env`. Verificamos la calidad del proyecto. |
| `fd supply [env]` | Ejecuta hasta el paso `supply` en el entorno `env`. Aprovisionamos la infraestructura necesaria. |
| `fd package [env]` | Ejecuta hasta el paso `package` en el entorno `env`. Empaquetamos el proyecto para su despliegue. |
| `fd deploy [env]` | Ejecuta hasta el paso `deploy` en el entorno `env`. Es el ultimo paso, desplegamos el projecto en el entorno indicado. |

**Flags comunes:**
*   `--yes` o `-y`: Salta las confirmaciones interactivas, para `fd init`
<!-- *   `--skip-test`: Omite los pasos de `test`.
*   `--skip-supply`: Omite los pasos de `supply`. -->

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Si tienes ideas, sugerencias o encuentras un error, por favor abre un [issue](https://github.com/jairoprogramador/fastdeploy/issues) o envía un [pull request](https://github.com/jairoprogramador/fastdeploy/pulls).

## 📄 Licencia

`fastdeploy` está distribuido bajo la [Apache License 2.0](https://github.com/jairoprogramador/fastdeploy/blob/main/LICENSE).
