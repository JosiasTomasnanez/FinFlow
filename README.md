# 💰 FinFlow - Metodología DevOps, Prácticas Ágiles e Ingeniería de Plataformas

---

## 1. 🏢 Descripción de la Empresa y Modelo de Negocio
FinFlow es una startup fintech en crecimiento que opera una billetera digital y una plataforma de pagos diseñada para individuos, comercios y PYMEs. La plataforma centraliza transferencias entre usuarios (P2P), saldos de billeteras digitales y el procesamiento de pagos mediante una API REST.

Nuestro modelo de ingresos se basa principalmente en comisiones por transacciones de comercios, planes comerciales premium y tarifas por operaciones de gestión de activos. La organización, compuesta por un equipo cross-functional de 20 personas, se divide en Desarrollo de Software, Plataforma/Operaciones, Producto y QA. Adoptamos **Scrum** junto con una madura **cultura DevOps/SRE** para acelerar la velocidad de entrega de funcionalidades, garantizando al mismo tiempo estabilidad y cumplimiento normativo de nivel financiero.

---

## 2. 🏗️ Concepto Arquitectónico y Evolución de la Infraestructura

### Arquitectura Actual (Fase MVP y Staging)
Para optimizar la velocidad del desarrollo inicial y minimizar los costos tempranos de infraestructura, FinFlow está diseñado actualmente como un **Monolito Modular**:
* **Backend:** Un servidor altamente eficiente escrito en **Go (Golang)** potenciado por el framework web **Gin-Gonic**, que encapsula la lógica de negocio de billeteras y pagos bajo módulos estructurados.
* **Frontend:** Una aplicación moderna de página única (SPA) construida con **React, TypeScript y Vite**, servida por el propio backend de Go o alojada de forma independiente.
* **Persistencia:** Instancias de **PostgreSQL** completamente aisladas por entorno, lo que garantiza una estricta separación de datos entre Staging y Producción sin la sobrecarga de grandes clusters de bases de datos durante la fase de validación.

### Roadmap Arquitectónico Empresarial (Estado Futuro)
A medida que la concurrencia de transacciones escale y los requisitos de cumplimiento normativo se vuelvan más estrictos, el roadmap arquitectónico de FinFlow contempla una migración fluida hacia:
1. **Transición a Microservicios:** Desacoplar el monolito modular en dominios especializados (Servicio de Autenticación, Ledger de Pagos, Gestión de Usuarios) con radios de impacto (*blast radiuses*) aislados.
2. **Capa de Datos Empresarial:** Consolidar en **PostgreSQL** para aprovechar el cumplimiento estricto de ACID a nivel empresarial, la gestión de conexiones concurrentes y el bloqueo a nivel de fila (*row-level locking*) para transacciones financieras seguras.
3. **Orquestación Nativa de la Nube:** Transición desde el alojamiento PaaS hacia **Amazon Web Services (AWS)**, gestionando las cargas de trabajo dentro de **Amazon EKS (Kubernetes)** con alta disponibilidad multirregión.

---

## 🛠️ 3. Stack Tecnológico y Ecosistema

### Aplicación Core
* **Backend:** Go (Golang) + Framework Gin Gonic
* **Frontend:** React + Vite + TypeScript

### DevOps y Gates de Calidad
* **Integración Continua (CI):** GitHub Actions
* **Pruebas Estáticas de Seguridad (SAST):** SonarCloud / SonarQube
* **Pruebas Dinámicas de Seguridad (DAST):** OWASP ZAP (Automatización Baseline)
* **Calidad de Código y Linters:** `golangci-lint`

### Observabilidad y Gestión de Releases
* **Ingesta de Métricas:** Prometheus (Scrapeando el endpoint interno `/metrics` de Gin)
* **Visualización de Datos:** Dashboards de Grafana (Implementación personalizada de múltiples servicios en Docker sobre Render)
* **Feature Flags y Despliegue Progresivo:** Servidor Unleash (Self-hosted) y SDK del Cliente en Go

---

## ⚙️ 4. Procesos y Pipeline de DevOps

### 🔄 Estrategia de Ramas (GitHub Flow) y Políticas del Repositorio
Para equilibrar la velocidad de una startup con el cumplimiento de software fintech, el repositorio implementa estrictamente **GitHub Flow**:
* `main`: La única fuente de verdad. Representa el código activo, estable y listo para integración.
* `feature/*` o `bugfix/*`: Ramas de corta duración creadas a partir de `main`. Una vez que el trabajo se verifica localmente, se abre un Pull Request (PR).
* **Políticas de Protección de Ramas:** Los pushes directos a `main` están estrictamente prohibidos. Fusionar (mergear) un PR requiere superar todos los gates de calidad automatizados (Build de CI, Linter, análisis de SonarCloud) y recibir revisiones obligatorias de los pares.

### 🚀 Estrategia de Entornos y Gestión de Releases con GitOps
Gestionamos dos entornos de infraestructura distintos, desacoplados del antipatrón de usar "entornos como ramas de Git":


```

[ Developer PR ] ➔ [ Merge to main ] ➔ [ Trigger Staging Deploy (Webhook) ] ➔ [ Run DAST / Security Scan ]
│
[ Production Release ] 🡠 [ Automated Tag Validation ] 🡠 [ Git Tag Created (vX.X.X) ] 🡠┘

```

1. **Entorno de Staging:**
   * **Disparador:** Cada PR aprobado y mergeado a `main` activa automáticamente un despliegue inmutable a Staging (Backend en Render, Frontend en Vercel).
   * **Bucle de Validación:** Una vez desplegado, se permite un tiempo de estabilización de 90 segundos. Luego, se ejecuta un escaneo automatizado **OWASP ZAP DAST** dinámicamente contra el entorno para verificar la postura de seguridad antes de la promoción.
2. **Entorno de Producción:**
   * **Disparador:** Se activa exclusivamente mediante Tags de Release de Git (`v*.*.*`). Esto establece un registro de auditoría claro.
   * **Mecanismo de Despliegue:** Aunque la arquitectura en la nube actual se basa en disparadores de promoción directa, el próximo roadmap de AWS introduce **Rolling Updates** nativos a través del ciclo de vida de los Pods de Kubernetes para garantizar despliegues con cero tiempo de inactividad (*zero-downtime*).
   * **Feature Flags:** Impulsado por **Unleash**. Las funcionalidades de negocio se desacoplan de los despliegues. El código se libera de forma segura a producción en estado latente y se habilita progresivamente para testers internos (alpha), grupos canario y, finalmente, al 100% de la base de clientes.

---

## 👥 5. Framework Ágil (Integración de Scrum + DevOps)
FinFlow opera en **Sprints de dos semanas** utilizando Scrum complementado por el **framework CALMS** (Cultura, Automatización, Lean, Medición, Compartir).

* **Ceremonies & Planning:** Los SREs evalúan el **Presupuesto de Errores (Error Budget)** actual durante el Sprint Planning. Los desarrolladores planifican la cobertura de la automatización de pruebas, y los requisitos se estructuran en Historias de Usuario estimadas en Story Points.
* **Feedback Continuo:** Las reuniones diarias (Daily Standups) de 15 minutos desbloquean las tareas entre áreas. Las retrospectivas fomentan la mejora continua de los procesos y aplican estrictamente una **cultura de post-mortems sin culpas (blameless)** tras cualquier incidente mayor en staging o producción.

---

## 🛡️ 6. Ingeniería de Confiabilidad del Sitio (SRE) y Observability

La salud de la plataforma se monitorea proactivamente a través de telemetría en vivo y métricas contractuales:
1. **Indicadores de Nivel de Servicio (SLI):** Medimos activamente los umbrales de latencia de las peticiones a la API, tasas de error HTTP 5xx, estado de conexión de la base de datos y métricas de ejecución de Prometheus.
2. **Objetivos de Nivel de Servicio (SLO):** Nos comprometemos con un estricto objetivo de **99.9% de disponibilidad mensual del sistema** y un límite de tasa de errores estrictamente por debajo del 1.0%.
3. **Presupuesto de Errores (Error Budget):** El marco definitivo para la toma de decisiones. Si las anomalías en producción agotan el Error Budget mensual asignado, los despliegues de producto se detienen instantáneamente y todo el equipo de ingeniería se enfoca al 100% en la reducción de deuda técnica y la estabilización de la infraestructura.

### Hubs de Telemetría en Vivo (Staging)
* **Target de la API de la Aplicación:** `https://finflow-backend-uv2f.onrender.com/metrics`
* **Plano de Control de Feature Flags:** `https://unleash-web-1e2b.onrender.com`
* **Dashboard de Métricas y Analíticas:** `https://prometheus-config-brn4.onrender.com` *(Ciclo de scraping vivo auto-sostenido)*

---

## 📦 Inicio Rápido para Desarrollo Local

### Estructura de Directorios
- `cmd/finflow/main.go`: Punto de entrada central de la API REST que inicializa los módulos de Go.
- `internal/service/`: Lógica del dominio de negocio, manejadores de ejecución de transferencias y gestión de estado.
- `frontend/`: Código fuente independiente de la aplicación de página única en React/Vite.

### Ejecutar el programa Localmente
1. **Compilar y empaquetar el Frontend:**
   ```bash
   cd frontend
   npm install
   npm run build

2. **Compilar y ejecutar el Backend en Go:**
  ```bash
  cd ..
  go build ./cmd/finflow
  ./finflow

  ```


3. **Acceder a la aplicación:** Abre `http://localhost:8080/` en tu navegador. La interfaz de React interactuará de forma nativa con los endpoints de la API de Go bajo `/api/*`.

Para ejecutar solo el frontend en modo de desarrollo con recarga en caliente (*hot-reload*):

  ```bash
  cd frontend
  npm run dev

  ```

## 📦 Inicio rapido usando Docker Compose

La forma más eficiente y confiable de levantar todo el ecosistema de FinFlow localmente es utilizando Docker Compose. Esta estrategia orquesta todos los módulos de la aplicación, las capas de persistencia y los pipelines de observabilidad detrás de un proxy inverso de Nginx, replicando la topología de red y seguridad de producción.

### Levantar el Ecosistema Local
Ejecuta el siguiente comando desde el directorio raíz para compilar y desplegar todos los contenedores en segundo plano:

  ```bash
  docker-compose up -d --build

  ```

### Matriz de Enrutamiento Local

Una vez que todos los contenedores se estabilizan, Nginx expone de forma segura cada servicio de la infraestructura bajo el puerto `80`, eliminando por completo los conflictos de CORS y el manejo complejo de puertos:

| Ruta del Gateway | Servicio de la Plataforma | Propósito Operativo |
| --- | --- | --- |
| `http://localhost/` | **React SPA Frontend** | Interfaz Web de Usuario lista para producción |
| `http://localhost/api/*` | **Go Backend (Gin)** | APIs REST Core y Motores del Ledger Financiero |
| `http://localhost/metrics` | **Prometheus Engine** | Ingesta de Telemetría y Scraping en Tiempo Real |
| `http://localhost/grafana` | **Plataforma Grafana** | Métricas de Rendimiento y Dashboards Analíticos |
| `http://localhost/unleash` | **Servidor Unleash** | Plano de Control y Gestión de Feature Flags |

*Para apagar de forma segura todo el stack de la infraestructura local, simplemente ejecuta:* `docker-compose down`
