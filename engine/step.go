package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Step executes a durable step with memoization.
func Step[T any](
	ctx *DurableContext,
	id string,
	fn func() (T, error),
) (T, error) {

	var zero T

	// ---- Crash simulation (if enabled) ----
	ctx.mu.Lock()
	ctx.stepCount++
	currentStep := ctx.stepCount
	crashAfter := ctx.crashAfter
	ctx.mu.Unlock()

	if crashAfter > 0 && currentStep > crashAfter {
		panic("simulated crash")
	}
	// --------------------------------------

	// 1. Generate logical sequence ID (thread-safe)
	ctx.mu.Lock()
	seq := ctx.seq
	ctx.seq++
	ctx.mu.Unlock()

	stepKey := fmt.Sprintf("%s-%d", id, seq)

	// 2. Check if step already executed
	var storedOutput string
	var status string

	err := ctx.DB.QueryRow(
		`SELECT status, output 
		 FROM steps 
		 WHERE workflow_id = ? AND step_key = ?`,
		ctx.WorkflowID,
		stepKey,
	).Scan(&status, &storedOutput)

	if err == nil && status == "COMPLETED" {
		var result T
		if err := json.Unmarshal([]byte(storedOutput), &result); err != nil {
			return zero, err
		}
		return result, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return zero, err
	}

	// 3. Mark step as IN_PROGRESS
	_, err = ctx.DB.Exec(
		`INSERT INTO steps (workflow_id, step_key, status, output)
		 VALUES (?, ?, ?, ?)`,
		ctx.WorkflowID,
		stepKey,
		"IN_PROGRESS",
		"",
	)
	if err != nil {
		return zero, err
	}

	// 4. Execute user function
	result, err := fn()
	if err != nil {
		return zero, err
	}

	// 5. Serialize output
	bytes, err := json.Marshal(result)
	if err != nil {
		return zero, err
	}

	// 6. Mark step as COMPLETED
	_, err = ctx.DB.Exec(
		`UPDATE steps
		 SET status = ?, output = ?
		 WHERE workflow_id = ? AND step_key = ?`,
		"COMPLETED",
		string(bytes),
		ctx.WorkflowID,
		stepKey,
	)
	if err != nil {
		return zero, err
	}

	return result, nil
}

// ------------------------------------------------------------
// BONUS: Automatic Step ID generation
// ------------------------------------------------------------

// AutoStep executes a durable step without requiring a manual step ID.
// Uniqueness and determinism are guaranteed by the internal sequence counter.
func AutoStep[T any](
	ctx *DurableContext,
	fn func() (T, error),
) (T, error) {

	// "__auto__" is only a label.
	// The real identity is (__auto__ + logical sequence number).
	return Step(ctx, "__auto__", fn)
}
