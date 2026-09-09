### 5. `MAP.md`

```markdown
# MAPA TOPOLÓGICO DE REPOSITORIO (SOURCE TREE VERIFICATION)
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`

El presente árbol declara la estructura certificable del código fuente. Toda desviación, supresión o anexo de directorios fantasma detectados durante el escrutinio del repositorio constituirá un hallazgo de No Conformidad en la Fase de Auditoría 1A.

```text
rx-dispatch/
├── assets/                  # Recursos estáticos UI e íconos base
├── bin/                     # Artifact output (excluido del ámbito de auditoría de fuente)
├── cmd/                     # PUNTOS DE ENTRADA (Main entrypoints)
│   ├── audit/main.go
│   ├── delivery/main.go
│   ├── gateway/main.go
│   ├── image/main.go
│   ├── reader/main.go
│   ├── result/main.go
│   ├── security/main.go
│   ├── storage/main.go
│   └── study/main.go
├── internal/                # CORE DEL DOMINIO (Lógica de negocio aislada)
│   ├── dicom/               # TCP SCP Handler & Association Managers
│   ├── models/              # Contratos de estructuras (Study, Job, AuditEvent)
│   └── webassets/           # Dependencias estáticas embebidas
├── scripts/                 # Infraestructura como código y automatización
├── winres/                  # Ensamblado de binarios PE (Windows)
├── config.json              # Mapeo maestro de parámetros de ejecución
├── go.mod                   # Manifiesto estricto de dependencias
└── README.md                # Documento normativo raíz