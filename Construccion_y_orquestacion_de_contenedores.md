# Construcción y orquestación de contenedores

## 1. Descripción de FinFlow

FinFlow es una startup fintech en crecimiento que opera una billetera digital y una plataforma de pagos diseñada para individuos, comercios y PYMEs. La plataforma centraliza transferencias entre usuarios (P2P), saldos de billeteras digitales y el procesamiento de pagos mediante una API REST.

Con el objetivo de facilitar el desarrollo, mantenimiento y despliegue de la plataforma, el proyecto se organiza mediante una estrategia de control de versiones basada en GitHub Flow. Esta organización permite mantener separados el código fuente de las aplicaciones, la infraestructura y los recursos de despliegue, favoreciendo un proceso de desarrollo colaborativo y alineado con prácticas modernas de DevOps.

## 2. Descripción de la arquitectura del cluster

La infraestructura de FinFlow se despliega sobre un único clúster de Kubernetes utilizando K3s como distribución ligera. Esta decisión permite contar con un entorno de orquestación completo, optimizando el consumo de recursos sin perder las funcionalidades propias de Kubernetes.

Con el objetivo de aislar los diferentes ambientes de trabajo, el clúster se organiza mediante cuatro namespaces, donde son destinados a una función específica.

Cada uno de los entornos de aplicación mantiene sus propios recursos de Kubernetes, tales como Ingress, Services y Pods, permitiendo que las aplicaciones se ejecuten de forma aislada mientras comparten la misma infraestructura física del clúster.

### 2.1 Desarrollo

Utilizado para el desarrollo e integración inicial de nuevas funcionalidades.

![Diagrama entorno de dev](./Diagramas/Diagrama%20entorno%20Infra.jpg)

### 2.2 Staging

Entorno de validación donde se realizan pruebas previas a la publicación de nuevas versiones.

![Diagrama entorno de Staging](./Diagramas/Diagrama%20arquitectura%20staging.jpg)

### 2.3 Producción

Ambiente destinado a la ejecución de la versión estable de la plataforma.

![Diagrama Entorno de Produccion](./Diagramas/Diagrama%20arquitectura%20produccion.jpg)

### 2.4 ArgoCD

Namespace reservado para la herramienta encargada de implementar la estrategia GitOps y gestionar la sincronización automática entre los repositorios Git y el estado del clúster.

## 3. Aplicaciones implementadas

La plataforma FinFlow se encuentra compuesta por un conjunto de servicios que trabajan de manera integrada para brindar las funcionalidades de la billetera digital. Cada componente cumple una responsabilidad específica dentro de la arquitectura, permitiendo mantener una separación clara entre la presentación, la lógica de negocio, la persistencia de datos, la observabilidad y la configuración dinámica de funcionalidades.

Los diferentes servicios implementados son:

- **Backend:** Desarrollado en Go utilizando el framework Gin, implementa la lógica de negocio de la plataforma. Expone una API REST encargada de gestionar usuarios, billeteras virtuales, transferencias entre cuentas y procesamiento de pagos. Además, publica un endpoint de métricas compatible con Prometheus para facilitar el monitoreo de la aplicación.

- **Frontend:** Implementado con React, constituye la interfaz web utilizada por los usuarios para interactuar con la plataforma. Su principal responsabilidad es presentar la información y consumir la API REST expuesta por el backend mediante solicitudes HTTP.

- **PostgreSQL:** Cada entorno dispone de una instancia aislada de PostgreSQL, garantizando la separación de los datos entre Desarrollo, Staging y Producción. La base de datos cuenta con almacenamiento persistente a través de PVC, evitando la pérdida de información ante el reinicio o la recreación de los pods. 

- **Redis:** Se utiliza Redis como sistema de caché entre el backend y PostgreSQL para reducir los tiempos de acceso a la información más consultada. Ante una solicitud, el backend verifica primero la existencia del dato en la caché, si no se encuentra (cache miss), realiza la consulta sobre PostgreSQL y almacena el resultado en Redis para futuras peticiones.

- **Monitoreo y observabilidad:** La plataforma incorpora una solución de monitoreo basada en Prometheus y Grafana. El backend expone métricas mediante un endpoint dedicado, las cuales son recolectadas por Prometheus y posteriormente visualizadas en dashboards de Grafana, permitiendo supervisar el estado y rendimiento de los pods levantados en la infraestructura.

- **Unleash - Feature Flags:** La gestión dinámica de funcionalidades se realiza mediante Unleash, permitiendo habilitar o deshabilitar características de la aplicación sin necesidad de generar un nuevo despliegue. Esta estrategia facilita la liberación gradual de funcionalidades y reduce el riesgo asociado a la publicación de nuevas versiones.

- **Argo Rollout:** Permite implementar estrategias de despliegue como Rolling update, Canary y Blue/Green. De esta manera las nuevas versiones de las aplicaciones pueden publicarse de forma gradual reduciendo el riesgo de fallos durante una actualización.

- **Secretos:** Administra los secretos de la infraestructura. Las credenciales son cifradas antes de almacenarse en el repositorio y solo pueden ser descifradas por el cluster.

- **Keda:** Proporciona el escalado automático de las aplicaciones mediante la creación dinámica de réplicas.

## 4. Helm

Con el objetivo de minimizar las configuraciones manuales y mantener una única definición de los recursos de Kubernetes, el proyecto utiliza Helm como gestor de paquetes para la automatización de los despliegues.

Cada aplicación cuenta con su correspondiente Helm Chart, el cual define la estructura de los recursos de Kubernetes mediante templates. Estas plantillas contienen la definición de los Deployments, Services, Ingress, ConfigMaps, PersistentVolumeClaims y demás recursos necesarios para el funcionamiento de la aplicación.

Los principales archivos de configuración corresponden a:

- `values-infra.yaml`: configuración de los componentes del entorno de desarrollo local.
- `values-infra-aws.yaml`: configuración de los componentes del entorno de desarrollo para AWS.
- `values-staging.yaml`: parámetros del entorno de Staging de manera local.
- `values-staging-aws.yaml`: configuración específica del despliegue en AWS para Staging.
- `values-prod.yaml`: parámetros del entorno de Producción de manera local.
- `values-prod-aws.yaml`: configuración específica del despliegue en AWS para Producción.
- `values-keda.yaml`: configuración del escalado automático mediante Keda.

## 5. Keda

Además de los valores correspondientes a cada entorno, la infraestructura incorpora un archivo de configuración para Keda. Este componente permite definir políticas de escalado automático de los Deployments en función de métricas o eventos, incrementando o reduciendo dinámicamente la cantidad de réplicas según la carga de trabajo.

Al igual que el resto de la infraestructura, su configuración se encuentra parametrizada mediante Helm, lo que facilita modificar umbrales de escalado, cantidad mínima y máxima de réplicas, y demás parámetros sin alterar los manifiestos base.

## 6. Scripts

Con el objetivo de simplificar las tareas de reducir la ejecución de comandos manuales, el proyecto incorpora un conjunto de scripts que automatizan las operaciones más frecuentes para el levantamiento y mantenimiento de la infraestructura.

### 6.1 Consola (dev.sh)

Se desarrolló una consola interactiva que centraliza el acceso a las principales herramientas del proyecto y automatiza la realización de port-forward necesarios para acceder a los servicios desplegados en Kubernetes.

Entre sus principales funcionalidades se encuentran:

- Inicio de los servicios locales.
- Detener de los servicios locales.
- Apertura automática de Grafana, Unleash, ArgoCD, ArgoRollout, app funcional de staging y app funcional de producción.
- Inicio del servicio de k9s para la administración del cluster en general.

### 6.2 Secretos (sealed-secrets.sh)

Para garantizar una gestión segura de las credenciales, el proyecto incorpora un script encargado de generar automáticamente los Sealed Secrets utilizados por Kubernetes.

El proceso consiste en:

- Leer las variables sensibles desde archivos `.env` locales, los cuales nunca forman parte del repositorio.
- Generar los objetos Secret correspondientes para cada entorno.
- Sellar dichos secretos mediante kubeseal, utilizando el controlador Sealed Secrets instalado en el clúster.
- Generar automáticamente los manifiestos Helm que posteriormente serán desplegados mediante ArgoCD.

El script contempla la generación de secretos independientes para los entornos de Staging, Producción e Infraestructura, permitiendo mantener aisladas las credenciales de cada ambiente. Además, incorpora una validación del contexto actual del clúster antes de realizar cualquier operación, reduciendo el riesgo de generar secretos sobre un entorno incorrecto.

## 7. Seguridad

Con el objetivo de proteger la comunicación entre los componentes de la plataforma y resguardar la información sensible, la infraestructura implementa mecanismos de seguridad propios de Kubernetes. 
En particular, se utilizan Secrets para la gestión de credenciales y Network Policies para controlar el tráfico de red entre los diferentes Pods y servicios del clúster.

### 7.1 Network policy

Inicialmente se aplica una política Default Deny, que bloquea todo el tráfico de entrada y salida de los Pods. A partir de esta configuración, únicamente se habilitan de forma explícita las comunicaciones necesarias para el funcionamiento de la plataforma.

Las principales reglas implementadas son:

- Se permite el tráfico DNS para que los Pods puedan resolver los nombres de los servicios internos del clúster.
- El Frontend puede recibir conexiones provenientes del Ingress Controller, actuando como único punto de acceso público a la aplicación.
- El Backend únicamente acepta solicitudes provenientes del Frontend y del servicio de Prometheus encargado de recolectar métricas.
- El Backend solo puede establecer conexiones hacia PostgreSQL, Redis y Unleash, impidiendo cualquier otra comunicación saliente no autorizada.
- PostgreSQL únicamente acepta conexiones provenientes del Backend.
- Redis únicamente acepta conexiones provenientes del Backend.

De esta forma, cada componente de la arquitectura únicamente puede comunicarse con los servicios estrictamente necesarios para cumplir su función, reduciendo la superficie de ataque y limitando el impacto potencial ante una posible vulnerabilidad.

### 7.2 Secretos

Las credenciales y datos sensibles de la aplicación, como contraseñas de bases de datos, claves de acceso y variables de configuración, se administran mediante Kubernetes Secrets. Esta estrategia evita almacenar información confidencial dentro del código fuente o de las imágenes de los contenedores, permitiendo que las aplicaciones consuman estos valores de forma segura durante su ejecución.

## 8. AWS

Además del despliegue sobre un clúster local de K3s, la infraestructura fue diseñada para ser portable y desplegarse sobre Amazon Web Services (AWS). Gracias al uso de Helm Charts parametrizados y a la estrategia GitOps, los mismos manifiestos pueden adaptarse a un entorno cloud modificando únicamente los archivos de configuración (values).

Para ello, la arquitectura contempla el uso de los siguientes servicios de AWS:

- **Amazon EC2:** Instancias virtuales utilizadas para alojar los recursos necesarios de la infraestructura cuando se requiere un entorno administrado manualmente o componentes auxiliares.

- **Amazon EKS:** Servicio administrado de Kubernetes utilizado para ejecutar el clúster en la nube, manteniendo la misma arquitectura implementada en K3s.

- **Application Load Balancer (ALB):** Balanceador de carga encargado de distribuir el tráfico HTTP/HTTPS hacia los servicios publicados dentro del clúster.

- **AWS Secrets Manager:** Servicio utilizado para almacenar y administrar credenciales, contraseñas y otra información sensible de forma segura.

- **Amazon Elastic Container Registry (ECR):** Registro privado de contenedores donde se almacenan las imágenes Docker utilizadas por las aplicaciones antes de su despliegue.

- **NAT Gateway:** Permite que los recursos ubicados en subredes privadas puedan acceder a Internet para descargar imágenes, dependencias o actualizaciones, sin exponerlos directamente al tráfico entrante.