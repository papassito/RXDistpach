# GOBERNANZA DEL REPOSITORIO Y CICLO DE VIDA DEL CÓDIGO
**DOCUMENT STATUS:** `BASELINE`
**PURPOSE:** Definir las reglas normativas para la interacción con el repositorio Git/GitHub.

Este documento establece el protocolo para la gestión de ramas, commits, Pull Requests (PRs) y la trazabilidad entre el trabajo planificado y el código fuente. Su cumplimiento es obligatorio para mantener la integridad del proceso de reconstrucción definido en `MANIFEST.md`.

---

## 1. Estrategia de Ramas (Branching Strategy)

Se adopta un modelo de ramas estricto para reflejar las fases del manifiesto.

*   **`main`**: Rama protegida. Representa el último estado `CERTIFIED`. Solo acepta fusiones desde ramas `release/*` aprobadas. No se permite el push directo.
*   **`develop`**: Rama de integración. Consolida funcionalidades que han alcanzado el nivel de evidencia `TESTED`. Es la base para las ramas `release/*`.
*   **`feature/<issue-id>-<description>`**: Ramas para nuevo desarrollo. Parten de `develop`. El `issue-id` debe corresponder a una tarea definida.
    *   *Ejemplo:* `feature/RX-123-implement-reader-contract`
*   **`fix/<issue-id>-<description>`**: Ramas para corrección de errores sobre `develop`.
*   **`audit/<phase>-<auditor>`**: Ramas de solo lectura para realizar las fases de auditoría (`ZERO MUTATION`). No se fusionan.
    *   *Ejemplo:* `audit/phase2-gemini`

---

## 2. Convención de Mensajes de Commit

Se utilizará el estándar Conventional Commits. Esto es mandatorio para la generación automática de changelogs y la trazabilidad.

**Formato:** `<type>(<scope>): <subject>`

*   **Types:** `feat`, `fix`, `build`, `chore`, `ci`, `docs`, `perf`, `refactor`, `revert`, `style`, `test`.
*   **Scope:** El componente afectado (ej. `reader`, `audit`, `contracts`).

*Ejemplos:*
```
feat(reader): Implementar endpoint /analyze
docs(contracts): Actualizar la especificación de HealthResponse
fix(storage): Corregir condición de carrera en la purga de disco
```

---

## 3. Proceso de Pull Request (PR)

Todo cambio de código debe ingresar a `develop` o `main` a través de un PR.

1.  **Título:** Debe seguir la convención de commits.
2.  **Descripción:** Debe usar una plantilla que incluya:
    *   **Propósito:** ¿Qué problema resuelve este PR?
    *   **Solución:** Descripción técnica de los cambios.
    *   **Evidencia Adjunta:** Enlace a los resultados de pruebas, validaciones estáticas, etc.
    *   **Nivel de Evidencia Alcanzado:** (`IMPLEMENTED`, `TESTED`, etc.).
3.  **Revisión:** Requiere al menos una aprobación de un `CODEOWNER` definido.
4.  **Checks:** Debe pasar todas las verificaciones de CI (build, lint, tests unitarios) antes de poder ser fusionado.

---

## 4. Gestión de Issues y Trazabilidad

*   Cada `feature` o `fix` debe estar asociado a un `issue` en el tracker del proyecto.
*   Los `issues` se etiquetarán con la fase del `MANIFEST.md` a la que pertenecen (ej. `phase-2`, `phase-4`).
*   Los hallazgos de auditoría (`FINDINGS`) se registrarán como `issues` con la etiqueta `finding` y la severidad correspondiente.

---

## 5. Propiedad del Código (`CODEOWNERS`)

Se mantendrá un archivo `.github/CODEOWNERS` para asignar la responsabilidad de revisión de cada componente o documento a individuos o equipos específicos, automatizando el proceso de solicitud de revisión en los PRs.