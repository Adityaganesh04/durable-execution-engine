package engine

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DurableContext struct {
	WorkflowID string
	DB         *sql.DB

	mu         sync.Mutex
	seq        int64
	crashAfter int
	stepCount  int
}

// NewDurableContext creates or resumes a workflow context
func NewDurableContext(workflowID string, dbPath string) (*DurableContext, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable better concurrency behavior for SQLite
	_, err = db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		return nil, err
	}

	// Create steps table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS steps (
		workflow_id TEXT,
		step_key TEXT,
		status TEXT,
		output TEXT,
		PRIMARY KEY (workflow_id, step_key)
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}

	ctx := &DurableContext{
		WorkflowID: workflowID,
		DB:         db,
		seq:        0,
	}

	return ctx, nil
}

// EnableCrashSimulation configures the context to crash after N steps
func (ctx *DurableContext) EnableCrashSimulation(n int) {
	ctx.crashAfter = n
}
