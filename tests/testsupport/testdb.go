package testsupport

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/jackc/pgx/v5/pgconn"
	sqlm "github.com/rubenv/sql-migrate"
)

type TestDB struct {
	DB      *bun.DB
	SQL     *sql.DB
	root    string
	dbName  string
	dsnUsed string
}

func repoRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
		}
		dir = parent
	}
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func sanitizeDBName(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := strings.Trim(b.String(), "_")
	if len(out) > 63 {
		out = out[:63]
	}
	return out
}

func deriveDBName(base, testName string) string {
	if forced := os.Getenv("TEST_DB_NAME"); forced != "" {
		return sanitizeDBName(forced)
	}
	return sanitizeDBName(fmt.Sprintf("%s_%s", base, testName))
}

func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	root := repoRoot()
	_ = godotenv.Load(filepath.Join(root, ".env.testing"))

	baseDSN := getenvDefault("PG_DSN",
		getenvDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/oficina_test?sslmode=disable"),
	)

	baseName := "temp"
	if u0, err := url.Parse(baseDSN); err == nil {
		if n := strings.TrimPrefix(u0.Path, "/"); n != "" {
			baseName = n
		}
	}

	dbName := deriveDBName(baseName, t.Name())

	u, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}
	admin := *u
	admin.Path = "/postgres"
	u.Path = "/" + dbName
	dsn := u.String()
	adminDSN := admin.String()

	dropDatabaseIfExists(t, adminDSN, dbName)
	createDatabase(t, adminDSN, dbName)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	db := bun.NewDB(sqlDB, pgdialect.New())

	if testing.Verbose() || os.Getenv("DEBUG_TESTDB") == "1" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		fmt.Printf("[testdb] db=%s dsn=%s\n", dbName, dsn)
	}

	migrationsDir := getenvDefault("MIGRATIONS_DIR", "migrations")
	migrationsTable := getenvDefault("MIGRATIONS_TABLE", "migrations")
	schemaPath := filepath.Join(root, migrationsDir)
	if _, err := os.Stat(schemaPath); err != nil {
		t.Fatalf("migrations path not found: %s (root=%s) err=%v", migrationsDir, root, err)
	}
	sqlm.SetTable(migrationsTable)
	if n, err := sqlm.Exec(sqlDB, "postgres", &sqlm.FileMigrationSource{Dir: schemaPath}, sqlm.Up); err != nil {
		t.Fatalf("sql-migrate up (schema em %s): %v", migrationsDir, err)
	} else if os.Getenv("DEBUG_TESTDB") == "1" {
		fmt.Printf("[testdb] applied schema migrations: %d\n", n)
	}

	testSeedersDir := getenvDefault("TEST_SEEDERS_DIR", "migrations/test_seeders")
	seedsTable := getenvDefault("SEEDS_TABLE", "seeds_migrations")
	seedsPath := filepath.Join(root, testSeedersDir)
	if st, err := os.Stat(seedsPath); err == nil && st.IsDir() {
		sqlm.SetTable(seedsTable)
		if n, err := sqlm.Exec(sqlDB, "postgres", &sqlm.FileMigrationSource{Dir: seedsPath}, sqlm.Up); err != nil {
			t.Fatalf("sql-migrate up (test seeds %s): %v", testSeedersDir, err)
		} else if os.Getenv("DEBUG_TESTDB") == "1" {
			fmt.Printf("[testdb] applied test seeds: %d\n", n)
		}
	} else {
		t.Fatalf("tests seeds path not found: %s (root=%s)", testSeedersDir, root)
	}

	tdb := &TestDB{DB: db, SQL: sqlDB, root: root, dbName: dbName, dsnUsed: dsn}

	t.Cleanup(func() {
		tdb.Close()
		dropDatabaseIfExists(t, adminDSN, dbName)
	})

	return tdb
}

func (tdb *TestDB) Close() { _ = tdb.DB.Close(); _ = tdb.SQL.Close() }

func dropDatabaseIfExists(t *testing.T, adminDSN, dbName string) {
	adminDB, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatalf("sql.Open(admin): %v", err)
	}
	defer adminDB.Close()

	if _, err := adminDB.Exec(`DROP DATABASE IF EXISTS "` + strings.ReplaceAll(dbName, `"`, `""`) + `" WITH (FORCE)`); err == nil {
		return
	}

	_, _ = adminDB.Exec(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1 AND pid <> pg_backend_pid()
	`, dbName)
	if _, err := adminDB.Exec(`DROP DATABASE IF EXISTS "` + strings.ReplaceAll(dbName, `"`, `""`) + `"`); err != nil {
		var pgErr *pgconn.PgError
		if !(errors.As(err, &pgErr) && pgErr.Code == "3D000") {
			t.Fatalf("DROP DATABASE %s: %v", dbName, err)
		}
	}
}

func createDatabase(t *testing.T, adminDSN, dbName string) {
	adminDB, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatalf("sql.Open(admin): %v", err)
	}
	defer adminDB.Close()

	stmt := `CREATE DATABASE "` + strings.ReplaceAll(dbName, `"`, `""`) + `"`
	if _, err := adminDB.Exec(stmt); err != nil {
		var pgErr *pgconn.PgError
		if !(errors.As(err, &pgErr) && pgErr.Code == "42P04") {
			t.Fatalf("CREATE DATABASE %s: %v", dbName, err)
		}
	}
}
