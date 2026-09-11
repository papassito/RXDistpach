# RX DISPATCH — Repository Map

**DOCUMENT STATUS:** `CONSOLIDATED BASELINE`
**PURPOSE:** A map of the declared repository structure.

---

## 1. DECLARED DIRECTORY STRUCTURE

This map reflects the repository structure as declared in the project's documentation. The physical existence and content of each item are subject to phased audits.

```text
RXDistpach/
├── rx-dispatch/             # Go project root
│   ├── cmd/                 # Entrypoints for the 9 domain services
│   │   ├── audit/
│   │   ├── delivery/
│   │   ├── gateway/
│   │   ├── image/
│   │   ├── reader/
│   │   ├── result/
│   │   ├── security/
│   │   ├── storage/
│   │   └── study/
│   ├── internal/            # Shared internal packages
│   │   ├── config/
│   │   ├── contracts/
│   │   ├── models/
│   │   ├── shared/
│   │   ├── transport/
│   │   └── webassets/
│   ├── docs/                # Authoritative documentation (this folder)
│   ├── scripts/             # Build, execution, and validation scripts
│   ├── go.mod               # Go module dependencies
│   └── go.sum
├── assets/                  # Source for static UI assets (e.g., logos)
├── bin/                     # Output for compiled binaries
├── config.json              # Service topology configuration
└── README.md                # Project root README
```

## 2. NOTES ON STRUCTURE

*   **`docs/` vs. Root:** This `docs/` directory is being established as the single source of truth for core documentation. Other duplicated documents in the repository are pending consolidation in a later phase.
*   **`tests/` Directory:** The project structure anticipates a `tests/` directory which is currently `MISSING` from the audited scope.