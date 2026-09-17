# RX DISPATCH — Repository Map

**DOCUMENT STATUS:** `CONSOLIDATED BASELINE`
**PURPOSE:** A map of the declared repository structure.

---

## 1. DECLARED DIRECTORY STRUCTURE

This map reflects the repository structure as declared in the project's documentation. The physical existence and content of each item are subject to phased audits.

```text
RXDistpach/
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   ├── audit/
│   ├── delivery/
│   ├── gateway/
│   ├── image/
│   ├── reader/
│   ├── result/
│   ├── security/
│   ├── storage/
│   └── study/
├── internal/
│   ├── contracts/
│   ├── models/
│   ├── shared/
│   └── transport/
├── docs/
├── scripts/
├── bin/
├── reports/
└── .pids/
```

## 2. NOTES ON STRUCTURE

*   **`docs/` vs. Root:** This `docs/` directory is being established as the single source of truth for core documentation. Other duplicated documents in the repository are pending consolidation in a later phase.
*   **`tests/` Directory:** The project structure anticipates a `tests/` directory which is currently `MISSING` from the audited scope.