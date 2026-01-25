# Native Durable Execution Engine  

---

## Overview

The project focuses on building a **production-grade, fault-tolerant backend system** using a **native durable execution pattern**.  
The engine ensures correctness and durability even in the presence of crashes, restarts, and concurrent execution.

---

## Repository Structure

```
durable-execution-engine/
├── engine/                # Core durable execution engine (Go)
├── docker-compose.yml     # Containerized execution setup
├── Prompts.txt            # AI usage and architectural decision log
└── README.md              # Project documentation
```

---

### Tech Stack
- **Go:** 1.25  
- **Database:** SQLite (WAL mode enabled)  
- **Containerization:** Docker & Docker Compose  

---

## Key Features

- **Fault Tolerance**
  - Implements a **Check → Act → Persist** pattern
  - Ensures completed steps are **never re-executed** after a crash

- **Concurrency**
  - Supports parallel step execution using `errgroup`
  - SQLite **WAL mode** enabled for safe concurrent reads and writes

- **Resilience**
  - Handles the **Zombie Step Problem**
  - Maintains explicit `IN_PROGRESS` and `COMPLETED` states in persistent storage

- **Bonus Feature**
  - **AutoStep ID generation**
  - Internal logical sequence counter removes need for manual step identifiers

---

## How to Run (Crash & Resume Demonstration)

The engine is **fully containerized**, allowing durability to be demonstrated without local environment setup issues.

---

### Start the Workflow (Normal Run)

```bash
cd durable-execution-engine
docker compose up --build
```

This runs the onboarding workflow **from start to finish**.

---

### Simulate a Crash (Durability Test)

This command intentionally crashes the workflow after a specific number of steps.  
The SQLite database is mounted to the host machine to preserve state across crashes.

#### Windows (PowerShell)
```powershell
docker run --rm -v ${PWD}\engine.db:/root/engine.db durable-execution-engine-durable-engine --crash-after=2
```

#### Linux / macOS / Git Bash
```bash
docker run --rm -v $(pwd)/engine.db:/root/engine.db durable-execution-engine-durable-engine --crash-after=2
```

**Expected Result:**
- The process panics and exits midway
- An `engine.db` file appears on the host machine

---

### Resume Execution

Restart the engine using the default command:

#### Windows (PowerShell)
```powershell
docker run --rm -v ${PWD}\engine.db:/root/engine.db durable-execution-engine-durable-engine
```

#### Linux / macOS / Git Bash
```bash
docker run --rm -v $(pwd)/engine.db:/root/engine.db durable-execution-engine-durable-engine
```

**Expected Result:**
- Previously completed steps are **automatically skipped**
- The workflow resumes and completes successfully
- This conclusively proves **durable execution**

---

## AI Usage Declaration

AI tools (**ChatGPT** and **Cursor**) were used to:

- Accelerate boilerplate generation
- Explore architectural patterns
- Validate edge cases and failure scenarios

All **core logic**, including:
- SQLite durability guarantees  
- Concurrency handling  
- Crash recovery semantics  
- Docker persistence  

was **manually reviewed, refined, and validated** through repeated testing.

A full and transparent record of prompts and manual interventions is provided in **`Prompts.txt`**.

---

## Final Notes

- The system guarantees **no re-execution of completed steps**
- Crash recovery is deterministic and repeatable
- Architecture is extensible for additional workflow steps
- Suitable for production-grade backend orchestration

---
