# RX DISPATCH — Mapa del repositorio

La raíz de Git y del módulo Go es `RXDistpach/`.

| Ubicación | Contenido |
| --- | --- |
| `cmd/reader/` | Servicio inicial y pruebas unitarias junto al código |
| `cmd/<servicio>/` | Futuras entradas audit, delivery, gateway, image, result, security, storage y study |
| `internal/contracts/` | Contratos de comunicación |
| `internal/models/` | Modelos de dominio |
| `internal/shared/` | Constantes compartidas |
| `internal/transport/` | Clientes de comunicación |
| `docs/` | Única ubicación de documentación del proyecto |
| `scripts/` | Herramientas de mantenimiento |
| `go.mod` | Única definición del módulo `rx-dispatch` |

## Reglas para nuevas incorporaciones

- Agregar contenido en la carpeta existente que corresponda; no duplicar la raíz ni los documentos.
- Los imports `rx-dispatch/internal/...` identifican el módulo Go, no una carpeta anidada.
- Crear carpetas pendientes solo cuando se agregue su implementación.
- `assets/` se reserva para recursos fuente; `internal/webassets/` para su integración Go, sin mantener copias manuales.
- `bin/`, `.pids/` y otras salidas locales quedan excluidas de Git.
- `config.json`, `internal/config/` y pruebas de integración se incorporarán cuando sean necesarios.
- `go.sum` lo administra Go al incorporar dependencias; no crearlo vacío.
- Ejecutar `scripts/check_structure.ps1` para detectar otra raíz, módulos anidados y archivos idénticos.

Los árboles y componentes del manifiesto representan el diseño previsto. Este mapa describe la organización actual.
