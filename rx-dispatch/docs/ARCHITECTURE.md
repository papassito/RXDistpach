# ARQUITECTURA DE RX DISPATCH BY KLIK
## 100% GO · CERO MONOLÍTICO · MODULAR · DESACOPLADO

El sistema RX DISPATCH ha sido diseñado y construido estrictamente bajo el mandato arquitectónico de cero monolito, eliminando cualquier punto único de ejecución o empaquetado de servicios en un único proceso.

---

## 1. Topología de Procesos Independientes

Cada capacidad principal existe como un ejecutable Go autónomo con ciclo de vida, configuración de puertos y responsabilidad exclusiva:

```text
                                [ CLIENTES / PACS / INTERFAZ ]
                                               │
                                      (HTTP / JSON REST)
                                               ▼
                              ┌─────────────────────────────────┐
                              │           RX GATEWAY            │
                              │           (Port 8080)           │
                              └────────┬──────────────┬─────────┘
                                       │              │
                   ┌───────────────────┼──────────────┼───────────────────┐
                   ▼                   ▼              ▼                   ▼
        ┌─────────────────────┐┌──────────────┐┌──────────────┐┌─────────────────────┐
        │     RX SECURITY     ││   RX STUDY   ││  RX STORAGE  ││      RX AUDIT       │
        │     (Port 8081)     ││ (Port 8082)  ││ (Port 8083)  ││     (Port 8088)     │
        └─────────────────────┘└──────┬───────┘└──────┬───────┘└─────────────────────┘
                                      │               │
                                      ▼               ▼
                               ┌──────────────┐┌──────────────┐
                               │   RX IMAGE   ││  RX READER   │
                               │ (Port 8084)  ││ (Port 8085)  │
                               └──────┬───────┘└──────┬───────┘
                                      │               │
                                      ▼               ▼
                               ┌──────────────┐┌──────────────┐
                               │  RX RESULT   ││ RX DELIVERY  │
                               │ (Port 8086)  ││ (Port 8087)  │
                               └──────────────┘└──────────────┘
```

---

## 2. Matriz de Responsabilidades y Binarios

| Componente | Binario | Puerto | Responsabilidad Única |
|---|---|---|---|
| **RX Gateway** | `rx-gateway` | 8080 | Ingress API, enrutamiento, coordinación de llamadas y topología de salud. |
| **RX Security** | `rx-security` | 8081 | Emisión y validación de tokens, autorización de scopes y control de accesos. |
| **RX Study** | `rx-study` | 8082 | Registro, metadata demográfica, identificación y ciclo de vida de los estudios. |
| **RX Storage** | `rx-storage` | 8083 | Almacén inmutable de imágenes originales, derivados y paquetes. Verificación SHA-256. |
| **RX Image** | `rx-image` | 8084 | Localización, cálculo de integridad criptográfica y generación de representaciones derivadas no destructivas. |
| **RX Reader** | `rx-reader` | 8085 | Lectura genérica automatizada de la imagen, sin carácter diagnóstico, con leyenda legal estricta. |
| **RX Result** | `rx-result` | 8086 | Ensamblado del resultado final con lectura, manifiesto de archivos, advertencias y checksum. |
| **RX Delivery** | `rx-delivery` | 8087 | Despacho multicanal (Email, SMS, Portal), token de seguimiento y estados de envío. |
| **RX Audit** | `rx-audit` | 8088 | Registro inmutable de eventos operacionales, trazabilidad y cumplimiento normativo. |

---

## 3. Principio de Preservación de Datos Radiológicos

Las imágenes originales NUNCA sufren modificaciones:
1. `RX STORAGE` recibe el binario y lo sella como `isImmutable: true` con su hash SHA-256.
2. `RX IMAGE` genera una representación derivada para inspección visual y análisis.
3. `RX READER` procesa la representación derivada y nunca tiene acceso de escritura sobre el original.

---

## 4. Principio de Lectura No Diagnóstica

Todo resultado generado por `RX READER` incluye obligatoriamente:
- `TYPE: GENERIC_AUTOMATED`
- `STATUS: WITHOUT_MEDICAL_SIGNATURE`
- `MEDICAL_REPORT: NOT_INCLUDED`
- **Leyenda Legal Obligatoria**:
  > *"La información presentada corresponde a una lectura genérica automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico oficial. Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio médico correspondiente."*
