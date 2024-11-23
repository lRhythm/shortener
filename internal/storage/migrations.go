package storage

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
)

const up = "up"

var migrations = map[string][]string{
	up: {
		//    ↓ - `storage/`.
		// ↓    - `internal/`.
		"../../db/migrations/00001_urls_table.up.sql",
		"../../db/migrations/00002_urls_table_original_url_unique_index.up.sql",
		"../../db/migrations/00003_urls_table_user_id_column.up.sql",
		"../../db/migrations/00004_urls_table_is_deleted_column.up.sql",
	},
	// If needed: `down: {"..."}`.
}

func migrate(db *sqlx.DB) error {
	_, fp, _, _ := runtime.Caller(0)
	dp := filepath.Dir(fp)
	for _, migration := range migrations[up] {
		f, e := os.Open(filepath.Join(dp, migration))
		if e != nil {
			return e
		}
		b, e := io.ReadAll(f)
		_ = f.Close()
		if e != nil {
			return e
		}
		_, e = db.Exec(string(b))
		if e != nil {
			return e
		}
	}
	return nil
}
