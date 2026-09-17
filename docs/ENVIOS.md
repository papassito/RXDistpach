# FLUJOS DE DESPACHO Y ENTREGA
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  

Este documento describe el flujo de trabajo para el despacho y la entrega de resultados.

## 1. Componente Responsable: `rx-delivery`

*   **Componente:** `rx-delivery`
*   **Puerto:** `8087/TCP`
*   **Responsabilidad:** Gestionar el envío del resultado final a los destinatarios configurados (ej. Webhook, Email, SMS).
*   **Estado:** `DEFINED` en la arquitectura.

## 2. Flujo de Despacho

1.  **Recepción de Tarea:** El servicio `rx-result` (o un orquestador como `rx-gateway`) invoca a `rx-delivery` con la carga útil final y el destino.
2.  **Encolamiento (Store-and-Forward):** `rx-delivery` recibe la instrucción de forma síncrona y la encola en un almacén persistente local. Esto asegura que la tarea de envío no se pierda si el destino externo no está disponible.
3.  **Intento de Entrega:** El servicio intenta enviar el resultado al destino externo a través de una conexión segura (HTTPS/TLS).
4.  **Manejo de Errores y Reintentos:**
    *   **Éxito (ACK):** Si la entrega es exitosa (ej. HTTP 200 OK del webhook), la tarea se marca como completada.
    *   **Fallo:** Si la entrega falla, el servicio no elimina la tarea. En su lugar, aplica una estrategia de reintentos con retroceso exponencial (*Exponential Backoff*) para no saturar al sistema de destino.
5.  **Auditoría:** Cada intento de entrega, ya sea exitoso o fallido, debe generar un evento de auditoría en `rx-audit`.

Este mecanismo de "Cola Resiliente" está diseñado para mitigar interrupciones de red (WAN) y garantizar la entrega eventual de los resultados.