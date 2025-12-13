package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Pré-passo: função de trigger compatível com tabelas que usam updated_at e/ou atualizado_em.
	if _, err := conn.Exec(ctx, `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS trigger
		LANGUAGE plpgsql
		AS $$
		BEGIN
			IF to_jsonb(NEW) ? 'updated_at' THEN
				NEW.updated_at = NOW();
			END IF;
			IF to_jsonb(NEW) ? 'atualizado_em' THEN
				NEW.atualizado_em = NOW();
			END IF;
			RETURN NEW;
		END;
		$$;
	`); err != nil {
		log.Fatalf("Unable to ensure update_updated_at_column trigger function: %v", err)
	}

	// Controle de versões do NEXO (não conflita com schema_migrations de outras ferramentas).
	if _, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS nexo_schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`); err != nil {
		log.Fatalf("Unable to ensure nexo_schema_migrations table: %v", err)
	}

	// Default: aplicar a partir do bloco de multi-unidade (035). Pode sobrescrever via MIGRATE_FROM.
	migrateFrom := 35
	if v := strings.TrimSpace(os.Getenv("MIGRATE_FROM")); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid MIGRATE_FROM=%q: %v", v, err)
		}
		migrateFrom = parsed
	}

	migrationsDir := "migrations"
	paths, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		log.Fatalf("Error listing migrations: %v", err)
	}
	if len(paths) == 0 {
		log.Printf("No .up.sql migrations found in %s", migrationsDir)
		return
	}

	sort.Strings(paths)
	prefixRe := regexp.MustCompile(`^(\d+)_`)

	appliedCount := 0
	skippedCount := 0
	for _, filePath := range paths {
		migration := filepath.Base(filePath)
		m := prefixRe.FindStringSubmatch(migration)
		if len(m) < 2 {
			continue
		}
		prefix, err := strconv.Atoi(m[1])
		if err != nil {
			log.Printf("Skipping %s: invalid numeric prefix: %v", migration, err)
			continue
		}
		if prefix < migrateFrom {
			skippedCount++
			continue
		}

		var alreadyApplied bool
		if err := conn.QueryRow(ctx,
			"SELECT EXISTS (SELECT 1 FROM nexo_schema_migrations WHERE version = $1)",
			migration,
		).Scan(&alreadyApplied); err != nil {
			log.Printf("Error checking migration %s: %v", migration, err)
			continue
		}
		if alreadyApplied {
			skippedCount++
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", migration, err)
			continue
		}

		fmt.Printf("Applying %s...\n", migration)
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Printf("  ❌ Error beginning tx for %s: %v", migration, err)
			continue
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			if isProbablyAlreadyApplied(err) {
				fmt.Printf("  ⚠️  %s - parece já aplicada (erro de duplicidade), marcando como aplicada\n", migration)
				_, insErr := conn.Exec(ctx,
					"INSERT INTO nexo_schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING",
					migration,
				)
				if insErr != nil {
					log.Printf("  ❌ Error recording %s: %v", migration, insErr)
					continue
				}
				appliedCount++
				continue
			}
			log.Printf("  ❌ Error applying %s: %v", migration, err)
			continue
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO nexo_schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING",
			migration,
		); err != nil {
			_ = tx.Rollback(ctx)
			log.Printf("  ❌ Error recording %s: %v", migration, err)
			continue
		}
		if err := tx.Commit(ctx); err != nil {
			log.Printf("  ❌ Error committing %s: %v", migration, err)
			continue
		}

		fmt.Printf("  ✅ %s applied successfully\n", migration)
		appliedCount++
	}

	fmt.Printf("\n✅ Migration process completed! applied=%d skipped=%d\n", appliedCount, skippedCount)

	/*
		CORRUPTED_OLD_MIGRATOR (mantido apenas para referência; não executa)

		// Controle simples de versões (evita re-aplicar migrations)
		// Observação: alguns ambientes já possuem schema_migrations(version BIGINT).
		_, err = conn.Exec(ctx, `

		// =============================================================================
		// Pré-passo: compatibilidade do trigger update_updated_at_column
		// Alguns schemas usam updated_at, outros usam atualizado_em.
		// =============================================================================
		_, err = conn.Exec(ctx, `
			package main

			import (
				"context"
				"fmt"
				"log"
				"os"
				"path/filepath"
				"regexp"
				"sort"
				"strconv"
				"strings"

				"github.com/jackc/pgx/v5"
				"github.com/joho/godotenv"
			)

			func main() {
				_ = godotenv.Load()

				dbURL := os.Getenv("DATABASE_URL")
				if dbURL == "" {
					log.Fatal("DATABASE_URL not set")
				}

				ctx := context.Background()
				conn, err := pgx.Connect(ctx, dbURL)
				if err != nil {
					log.Fatalf("Unable to connect to database: %v", err)
				}
				defer conn.Close(ctx)

				// Pré-passo: função de trigger compatível com ambos schemas.
				if _, err := conn.Exec(ctx, `
					CREATE OR REPLACE FUNCTION update_updated_at_column()
					RETURNS trigger
					LANGUAGE plpgsql
					AS $$
					BEGIN
						IF to_jsonb(NEW) ? 'updated_at' THEN
							NEW.updated_at = NOW();
						END IF;
						IF to_jsonb(NEW) ? 'atualizado_em' THEN
							NEW.atualizado_em = NOW();
						END IF;
						RETURN NEW;
					END;
					$$;
				`); err != nil {
					log.Fatalf("Unable to ensure update_updated_at_column trigger function: %v", err)
				}

				// Controle de versões do NEXO (não conflita com schema_migrations de outras ferramentas).
				if _, err := conn.Exec(ctx, `
					CREATE TABLE IF NOT EXISTS nexo_schema_migrations (
						version TEXT PRIMARY KEY,
						applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
					);
				`); err != nil {
					log.Fatalf("Unable to ensure nexo_schema_migrations table: %v", err)
				}

				migrateFrom := 35
				if v := strings.TrimSpace(os.Getenv("MIGRATE_FROM")); v != "" {
					parsed, err := strconv.Atoi(v)
					if err != nil {
						log.Fatalf("Invalid MIGRATE_FROM=%q: %v", v, err)
					}
					migrateFrom = parsed
				}

				migrationsDir := "migrations"
				paths, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
				if err != nil {
					log.Fatalf("Error listing migrations: %v", err)
				}
				if len(paths) == 0 {
					log.Printf("No .up.sql migrations found in %s", migrationsDir)
					return
				}

				sort.Strings(paths)
				prefixRe := regexp.MustCompile(`^(\d+)_`)

				appliedCount := 0
				skippedCount := 0

				for _, filePath := range paths {
					migration := filepath.Base(filePath)

					m := prefixRe.FindStringSubmatch(migration)
					if len(m) < 2 {
						continue
					}
					prefix, err := strconv.Atoi(m[1])
					if err != nil {
						log.Printf("Skipping %s: invalid numeric prefix: %v", migration, err)
						continue
					}
					if prefix < migrateFrom {
						skippedCount++
						continue
					}

					var alreadyApplied bool
					if err := conn.QueryRow(ctx,
						"SELECT EXISTS (SELECT 1 FROM nexo_schema_migrations WHERE version = $1)",
						migration,
					).Scan(&alreadyApplied); err != nil {
						log.Printf("Error checking migration %s: %v", migration, err)
						continue
					}
					if alreadyApplied {
						skippedCount++
						continue
					}

					content, err := os.ReadFile(filePath)
					if err != nil {
						log.Printf("Error reading %s: %v", migration, err)
						continue
					}

					fmt.Printf("Applying %s...\n", migration)
					tx, err := conn.Begin(ctx)
					if err != nil {
						log.Printf("  ❌ Error beginning tx for %s: %v", migration, err)
						continue
					}

					if _, err := tx.Exec(ctx, string(content)); err != nil {
						_ = tx.Rollback(ctx)
						if isProbablyAlreadyApplied(err) {
							fmt.Printf("  ⚠️  %s - parece já aplicada (erro de duplicidade), marcando como aplicada\n", migration)
							_, insErr := conn.Exec(ctx,
								"INSERT INTO nexo_schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING",
								migration,
							)
							if insErr != nil {
								log.Printf("  ❌ Error recording %s: %v", migration, insErr)
								continue
							}
							appliedCount++
							continue
						}
						log.Printf("  ❌ Error applying %s: %v", migration, err)
						continue
					}

					if _, err := tx.Exec(ctx,
						"INSERT INTO nexo_schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING",
						migration,
					); err != nil {
						_ = tx.Rollback(ctx)
						log.Printf("  ❌ Error recording %s: %v", migration, err)
						continue
					}
					if err := tx.Commit(ctx); err != nil {
						log.Printf("  ❌ Error committing %s: %v", migration, err)
						continue
					}

					fmt.Printf("  ✅ %s applied successfully\n", migration)
					appliedCount++
				}

				fmt.Printf("\n✅ Migration process completed! applied=%d skipped=%d\n", appliedCount, skippedCount)
			}

			func isProbablyAlreadyApplied(err error) bool {
				msg := err.Error()
				return strings.Contains(msg, "already exists") ||
					strings.Contains(msg, "duplicate") ||
					strings.Contains(msg, "already applied")
			}

	*/

}

func isProbablyAlreadyApplied(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "already applied")
}
