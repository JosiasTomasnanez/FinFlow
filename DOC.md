# Introducción a la Metodología DevOps - Grupo 9

## Integrantes

* Josias Ñañez
* Fernando Ereño
* Franco Navarro
* Jerónimo Massaro
* Marcos Gabriel Reynoso

**Empresa:** FinFlow
**Fecha de presentación:** 2026

---

# Introduccion

Somos una startup fintech que implementa metodologías ágiles y prácticas DevOps con el objetivo de acelerar la entrega de valor, mejorar la calidad del software y garantizar la estabilidad de la plataforma.

Este documento describe:

* El funcionamiento de la empresa.
* Las tecnologías utilizadas.
* La interacción entre los diferentes roles.
* Los procesos de desarrollo y despliegue.
* La implementación de metodologías ágiles.
* La integración de prácticas DevOps y DevSecOps.

---

# 1. Introducción

FinFlow es una startup fintech especializada en la gestión de billeteras virtuales y transferencias electrónicas.

La plataforma permite a usuarios y organizaciones:

* Gestionar cuentas digitales.
* Realizar transferencias.
* Administrar usuarios.
* Gestionar movimientos financieros.
* Integrar futuros servicios financieros digitales.

El modelo de negocio se basa en la prestación de servicios financieros digitales, cobrando comisiones por transacciones y servicios asociados.

La empresa cuenta con aproximadamente veinte integrantes distribuidos entre áreas de producto, desarrollo, calidad, operaciones e infraestructura.

FinFlow adopta Scrum como metodología ágil principal y una cultura DevOps orientada a la automatización, la colaboración y la mejora continua.

---

# 2. Funcionamiento de la Empresa

La plataforma se encuentra diseñada bajo una arquitectura de servicios desacoplados, permitiendo evolucionar hacia una arquitectura de microservicios a medida que aumenten las necesidades de escalabilidad.

Actualmente el sistema se ejecuta mediante contenedores Docker coordinados por Docker Compose.

La estrategia tecnológica busca:

* Escalabilidad horizontal.
* Alta disponibilidad.
* Automatización de procesos.
* Seguridad integrada desde el desarrollo.
* Observabilidad y monitoreo continuo.

La organización promueve una cultura DevOps donde desarrollo, operaciones, calidad y seguridad trabajan de manera colaborativa, eliminando silos organizacionales.

Además, se aplica el enfoque Shift-Left, incorporando controles de calidad y seguridad desde las primeras etapas del ciclo de desarrollo.

---

# 3. Tecnologías Utilizadas

## 3.1 Desarrollo

### Backend

* Go
* Garlic Framework

### Frontend

* JavaScript
* React

### Base de Datos

* PostgreSQL (planificado para próximas versiones)

---

## 3.2 Control de Versiones

* Git
* GitHub

Modelo de trabajo:

* GitHub Flow

---

## 3.3 CI/CD

* GitHub Actions

Actualmente se encuentran implementados pipelines para:

* Build automático
* Validación de formato
* Análisis estático de código
* Quality Gates

---

## 3.4 Contenedores

### Estado Actual

* Docker
* Docker Compose

La aplicación puede ejecutarse mediante:

```bash
docker compose up -d
```

### Evolución Planificada

* Kubernetes para orquestación de contenedores.

---

## 3.5 Infraestructura Cloud

### Estado Actual

* Render (entornos de prueba)

### Arquitectura Objetivo

* AWS

La infraestructura futura estará orientada a alta disponibilidad, escalabilidad y automatización de despliegues.

---

## 3.6 Monitoreo y Observabilidad

### Planificado

* Prometheus para recolección de métricas.
* Grafana para dashboards y observabilidad.

Indicadores monitoreados:

* Disponibilidad.
* Latencia.
* Errores.
* Consumo de recursos.
* Estado de despliegues.

---

## 3.7 Seguridad (DevSecOps)

### Implementado

#### SonarQube Cloud

Permite:

* Detección de bugs.
* Detección de vulnerabilidades.
* Detección de code smells.
* Aplicación de Quality Gates.

### Planificado

#### Semgrep

Análisis estático de seguridad mediante reglas configurables.

#### Dependabot

Monitoreo de dependencias y vulnerabilidades conocidas.

#### OWASP ZAP

Análisis dinámico de seguridad sobre aplicaciones en ejecución.

---

## 3.8 Testing y Performance

### Estado Actual

Actualmente no existen pruebas automatizadas implementadas.

### Planificado

* Unit Testing.
* Cobertura de código.
* k6 para pruebas de carga y estrés.

---

# 4. Procesos

## 4.1 Estrategia de Branching

FinFlow utiliza GitHub Flow.

Ramas utilizadas:

### main

Contiene el código estable.

### feature-*

Cada funcionalidad se desarrolla en una rama independiente.

Ejemplos:

* feature-login
* feature-transferencias
* feature-usuarios

No existen ramas develop ni staging.

---

## 4.2 Flujo de Commits

Se adopta Conventional Commits.

Ejemplos:

```text
feat: add login feature
fix: validate transfer amount
docs: update deployment guide
chore: update dependencies
```

Esto permite mantener trazabilidad y generar historial consistente.

---

## 4.3 Políticas del Repositorio

La rama principal posee protección activa.

Requisitos:

* Pull Request obligatorio.
* Dos aprobaciones obligatorias.
* Ejecución satisfactoria de pipelines.
* Bloqueo automático ante fallos de calidad.

---

## 4.4 Integración Continua (CI)

### Pipeline Build & Test

Ejecutado en:

* Pull Requests.
* Push sobre main.

Validaciones:

* Build del backend.
* Build del frontend.
* gofmt.
* golangci-lint.

### Pipeline SonarQube

Ejecutado en:

* Pull Requests.
* Push sobre main.

Validaciones:

* Bugs.
* Vulnerabilidades.
* Code Smells.
* Quality Gate.

---

## 4.5 Entrega Continua (CD)

FinFlow adopta una estrategia de Continuous Delivery.

Proceso:

1. Desarrollo en rama feature.
2. Pull Request.
3. Revisión de código.
4. Ejecución de pipelines.
5. Merge hacia main.
6. Despliegue automático a staging.
7. Generación de release.
8. Despliegue controlado a producción.

No se utiliza Continuous Deployment, ya que la liberación a producción requiere una decisión explícita mediante releases.

---

## 4.6 Versionado

Se utiliza Semantic Versioning (SemVer).

Formato:

```text
MAJOR.MINOR.PATCH
```

Ejemplos:

```text
v1.0.0
v1.1.0
v1.1.1
```

---

## 4.7 Estrategias de Despliegue

### Planificadas

* Rolling Updates.
* Feature Flags.

Estas estrategias permitirán minimizar riesgos durante la liberación de nuevas versiones.

---

## 4.8 Shift-Left

La calidad se integra desde etapas tempranas mediante:

* Revisiones de Pull Requests.
* SonarQube.
* Quality Gates.
* Linters.
* Validaciones automáticas en CI.

---

# 5. Metodologías Ágiles

FinFlow utiliza Scrum como marco principal de trabajo.

## Roles

### Product Owner

Define prioridades de negocio y gestiona el Product Backlog.

### Scrum Master

Facilita la metodología y elimina impedimentos.

### Equipo de Desarrollo

Implementa funcionalidades y participa en actividades técnicas y de calidad.

### DevOps Engineers

Gestionan pipelines, automatización e infraestructura.

### QA Engineers

Definen estrategias de validación y testing.

---

## Eventos Scrum

### Sprint Planning

Definición de objetivos y selección de historias.

### Daily Scrum

Sincronización diaria del equipo.

### Sprint Review

Presentación de resultados a stakeholders.

### Sprint Retrospective

Identificación de mejoras para el siguiente sprint.

---

## Integración Scrum + DevOps

Cada sprint genera software potencialmente desplegable.

Las funcionalidades desarrolladas pasan por:

* Pull Requests.
* Validaciones automáticas.
* Revisión de código.
* Entrega continua.

---

# 6. Site Reliability Engineering (SRE)

Como evolución futura de la plataforma se incorporarán prácticas SRE.

## SLI

* Disponibilidad.
* Latencia.
* Tasa de errores.
* Éxito de transacciones.

## SLO

* Disponibilidad mensual del 99.9%.
* Tasa de errores inferior al 1%.
* Latencia dentro de umbrales definidos.

## Error Budget

Permitirá equilibrar estabilidad operativa y velocidad de innovación.

---

# 7. Conclusión

La estrategia DevOps de FinFlow busca combinar automatización, calidad y seguridad para acelerar la entrega de valor sin comprometer la estabilidad de la plataforma.

Actualmente se encuentran implementadas prácticas como GitHub Flow, protección de ramas, Pull Requests obligatorias, análisis estático mediante SonarQube, validación de calidad y despliegue basado en contenedores.

Como evolución futura se incorporarán pruebas automatizadas, monitoreo avanzado, análisis dinámico de seguridad, Kubernetes, AWS y estrategias de despliegue progresivo, consolidando una plataforma alineada con las mejores prácticas modernas de DevOps, DevSecOps y SRE.

