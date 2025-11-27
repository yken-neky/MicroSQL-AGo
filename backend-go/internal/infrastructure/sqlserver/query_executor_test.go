package sqlserver

import (
	"testing"
)

func TestValidateQuery_forbidden_keywords_word_boundary(t *testing.T) {
	e := NewSQLServerQueryExecutor()

	// query containing the word BACKUP should be rejected
	if err := e.ValidateQuery("BACKUP DATABASE mydb"); err == nil {
		t.Fatalf("expected backup keyword to be forbidden")
	}

	// query containing 'backupset' should NOT be rejected after word-boundaries
	if err := e.ValidateQuery("SELECT * FROM msdb.dbo.backupset"); err != nil {
		t.Fatalf("unexpected validation error for backupset: %v", err)
	}
}
