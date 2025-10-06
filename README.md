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

**`fastdeploy` (o `fd`)** es una herramienta CLI diseñada para eliminar la complejidad y la repetición de los procesos de despliegue. Olvídate de los scripts frágiles, los largos `READMEs` y la pregunta "¿cómo se desplegaba esto?". Con `fastdeploy`, estandarizas tus despliegues usando plantillas reutilizables, permitiendo que cualquier desarrollador, en cualquier equipo, pueda desplegar cualquier aplicación de forma segura y predecible.

**Define tu proceso de despliegue una vez, y ejecútalo miles de veces con tre simples comandos.**

## ✨ Características Principales

*   **⚙️ Agnostico a la Tecnología:** ¿Java, Node.js, Python, Go? ¿Terraform, Docker, Kubernetes? `fastdeploy` orquesta cualquier herramienta que puedas ejecutar en un shell.
*   **📄 Infraestructura como Plantilla:** Centraliza la lógica de tus despliegues (steps, variables, entornos) en un repositorio de plantillas. Estandariza las buenas prácticas y evoluciona tu infraestructura sin tocar tus microservicios.
*   **🚀 Despliegues en 3 Pasos:** Clona tu microservicio y ejecuta `fd init`, `fd test`, y `fd deploy`. Eso es todo.
*   **✅ Verificación Continua:** El estado de cada despliegue se guarda, permitiendo validaciones y evitando ejecuciones accidentales en entornos incorrectos.
*   **💻 Experiencia de Desarrollador Primero:** Comandos intuitivos, feedback claro y la abstracción perfecta para que los desarrolladores se centren en lo que importa: el código.

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
curl -sL https://github.com/jairoprogramador/fastdeploy/releases/latest/download/fastdeploy_Linux_x86_64.tar.gz | tar xz
sudo mv fd /usr/local/bin/
```

### Windows

1.  Descarga el archivo `.zip` desde la [página de Releases](https://github.com/jairoprogramador/fastdeploy/releases).
2.  Descomprime el archivo.
3.  Añade el ejecutable `fd.exe` a tu `PATH`.


## 🏁 Guía de Inicio Rápido: Desplegando un Microservicio Java

Vamos a desplegar un microservicio Java que utiliza **Terraform** para provisionar la infraestructura en **Azure** (ACR, AKS) y se empaqueta con **Docker**.

Toda la lógica de este despliegue está definida en nuestra plantilla de ejemplo:
➡️ **[jairoprogramador/mydeploy](https://github.com/jairoprogramador/mydeploy)**

Este repositorio de plantillas contiene los `steps`, `variables` y la definición de los `environments` (ej: `sandbox`, `stagin`, `produccion`).

### Paso 1: Inicializa tu Proyecto

Clona el microservicio que quieres desplegar. Una vez dentro del directorio, ejecuta:

```sh
fd init
```

`fastdeploy` detectará que no está inicializado y te hará un par de preguntas para crear el archivo de configuración local `.fastdeploy/dom.yaml`. Este archivo vincula tu proyecto con la plantilla de despliegue.

```yaml
# .fastdeploy/dom.yaml (Ejemplo generado)
product:
  name: "core-banking"
project:
  name: "accounts-ms"
  revision: "v1.0.0"
technology:
  type: "backend"
  solution: microservicio
  stack: springboot
template:
  repository_url: "https://github.com/jairoprogramador/mydeploy.git"
  ref: "main"
```

### Paso 2: Prueba el Despliegue en un Entorno

Antes de desplegar, puedes validar que todo está en el entorno de desarrollo. El comando `test` ejecuta los comandos definidos en la plantilla referentes a las pruebas.

```sh
# Ejecuta los pasos de prueba para el entorno 'sand'
fd test sand
```

Esto podría, por ejemplo, ejecutar los test unitarios, la conexión con Azure, validar la versión de Terraform, compilar el proyecto, verificar pull request sin desplegarlo.

### Paso 3: Despliega

Una vez que las pruebas pasan, estás listo para desplegar. El comando `deploy` ejecuta la secuencia completa de pasos definidos en la plantilla, por ejemplo para el entorno de sandbox.

```sh
# Despliega en el entorno 'sand'
fd deploy sand
```
`fastdeploy` orquestará todo el proceso:
1.  Clonará la plantilla `mydeploy`.
2.  Ejecutará `terraform apply` para provisionar ACR y AKS.
3.  Construirá la imagen Docker de tu microservicio.
4.  Subirá la imagen al Azure Container Registry (ACR).
5.  Desplegará la aplicación en Azure Kubernetes Service (AKS).

¡Y listo! Tu microservicio está desplegado.

## 📚 Comandos Básicos

| Comando | Descripción |
| :--- | :--- |
| `fd init` | Inicializa un proyecto creando el archivo `.fastdeploy/dom.yaml`. |
| `fd [step] [env]` | Ejecuta un despliegue hasta el `step` indicado en el entorno `env`. |
| `fd test [env]` | Ejecuta solo los pasos de verificación (`test`) en el entorno `env`. |
| `fd supply [env]` | Ejecuta los pasos de aprovisionamiento de infraestructura (`supply`). |
| `fd deploy [env]` | Ejecuta todos los pasos hasta el despliegue final (`deploy`). |

**Flags comunes:**
*   `--yes` o `-y`: Salta las confirmaciones interactivas, para `fd init`
<!-- *   `--skip-test`: Omite los pasos de `test`.
*   `--skip-supply`: Omite los pasos de `supply`. -->

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Si tienes ideas, sugerencias o encuentras un error, por favor abre un [issue](https://github.com/jairoprogramador/fastdeploy/issues) o envía un [pull request](https://github.com/jairoprogramador/fastdeploy/pulls).

## 📄 Licencia

`fastdeploy` está distribuido bajo la [Licencia MIT](https://github.com/jairoprogramador/fastdeploy/blob/main/LICENSE).
