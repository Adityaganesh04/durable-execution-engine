Overview

This repository contains the submission for the Zeotap Software Engineer Intern – Assignment 1.
The project focuses on building a production-grade, fault-tolerant backend system using a native durable execution pattern.

Repository Structure

durable-execution-engine/
A Go-based native durable execution engine using SQLite and Docker
Status: Completed (Crash and Resume capability verified)

Prompts.txt
A detailed log of AI tools used and all manual architectural decisions made

Assignment 1: Native Durable Execution Engine

Tech Stack: Go 1.25, SQLite, Docker

Key Features:
-Fault Tolerance
Implements a durable Check–Act–Persist pattern to ensure completed steps are never re-executed after a crash.

-Concurrency
Supports parallel step execution using errgroup with SQLite WAL mode enabled for safe concurrent readers and writers.

-Resilience
Handles the Zombie Step problem by maintaining explicit IN_PROGRESS and COMPLETED states in persistent storage.

-Bonus Feature
Automatic step ID generation (AutoStep) using an internal logical sequence counter, eliminating the need for manual step identifiers.

How to Run (Crash and Resume Demonstration)

The engine is fully containerized so durability can be demonstrated without local environment setup issues.

1. Start the Workflow (Normal Run)
cd durable-execution-engine
docker compose up --build


This runs the onboarding workflow from start to finish.

2. Simulate a Crash (Durability Test)

This command intentionally crashes the workflow after a specific number of steps.
The SQLite database is mounted to the host machine to preserve state across crashes.

Windows (PowerShell):

docker run --rm -v ${PWD}:/root/ durable-execution-engine-durable-engine --crash-after=2


Linux / macOS / Git Bash:

docker run --rm -v $(pwd):/root/ durable-execution-engine-durable-engine --crash-after=2


Expected Result:
The process panics and exits midway. An engine.db file will be visible on the host machine.

3. Resume Execution

Restart the engine using the default command:

docker compose up


Expected Result:
Previously completed steps are automatically skipped, and the workflow resumes and completes successfully, proving durability.

AI Usage Declaration

As permitted by the assignment instructions, AI tools (ChatGPT and Cursor) were used to accelerate boilerplate generation and explore architectural patterns.
All core logic, including SQLite durability guarantees, concurrency handling, crash recovery, and Docker persistence, was manually reviewed, refined, and validated through repeated testing.

A full and transparent record of prompts and manual interventions is provided in Prompts.txt.