package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Tentar carregar .env
	_ = godotenv.Load("../../.env")

	// Pega URL do banco ou usa default
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/barber_nexo?sslmode=disable"
		fmt.Println("⚠️  DATABASE_URL não encontrada, usando padrão:", connStr)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Falha ao conectar no banco: %v", err)
	}
	fmt.Println("✅ Conectado ao banco de dados com sucesso.")

	// 1. Buscar Unidade Mangabeiras
	fmt.Println("\n🔍 Buscando unidade 'Mangabeiras'...")
	var unitID, unitName string
	err = db.QueryRow("SELECT id, name FROM units WHERE name ILIKE '%Mangabeira%' LIMIT 1").Scan(&unitID, &unitName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("❌ Unidade 'Mangabeiras' NÃO encontrada.")
			listarTodasUnidades(db)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("✅ Unidade Encontrada: %s (ID: %s)\n", unitName, unitID)

	// 2. Verificar Vínculos
	auditarTabela(db, "professionals", unitID)
	auditarTabela(db, "categorias", unitID) // Verifica schema se é 'categories' ou 'categorias'
	auditarTabela(db, "servicos", unitID)
	auditarTabela(db, "users", unitID) // Via user_units

	// 3. Verificar UserUnits
	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM user_units WHERE unit_id = $1", unitID).Scan(&userCount)
	if err != nil {
		fmt.Printf("❌ Erro ao contar user_units: %v\n", err)
	} else {
		fmt.Printf("📊 Usuários vinculados a esta unidade: %d\n", userCount)
	}

	// 4. Verificar vazamento de dados (NULL unit_id)
	fmt.Println("\n🛡️  Auditoria de Isolamento (Registros sem unit_id):")
	verificarIsolamento(db, "professionals")
	verificarIsolamento(db, "categorias")
	verificarIsolamento(db, "servicos")
	verificarIsolamento(db, "appointments")
	verificarIsolamento(db, "commands")
}

func listarTodasUnidades(db *sql.DB) {
	fmt.Println("📋 Unidades disponíveis:")
	rows, err := db.Query("SELECT id, name FROM units")
	if err != nil {
		log.Println("Erro ao listar unidades:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		fmt.Printf("- %s (ID: %s)\n", name, id)
	}
}

func auditarTabela(db *sql.DB, tableName, unitID string) {
	// Ajuste para tabelas que podem ter nomes diferentes ou user_units
	if tableName == "users" {
		return // Já tratado separadamente
	}

	// Tenta verificar se a tabela existe antes
	var exists bool
	db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&exists)
	if !exists {
		// Tenta singular/plural ou tradução básica
		if tableName == "categorias" { tableName = "categories" }
		if tableName == "servicos" { tableName = "services" }
		
		// Recheck
		db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&exists)
		if !exists {
			fmt.Printf("⚠️  Tabela '%s' não encontrada no banco.\n", tableName)
			return
		}
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE unit_id = $1", tableName)
	err := db.QueryRow(query, unitID).Scan(&count)
	if err != nil {
		fmt.Printf("❌ Erro ao auditar tabela %s: %v\n", tableName, err)
		return
	}
	
	status := "✅"
	if count == 0 {
		status = "⚠️ "
	}
	fmt.Printf("%s Registros em '%s' para esta unidade: %d\n", status, tableName, count)
}

func verificarIsolamento(db *sql.DB, tableName string) {
	// Normalização de nomes
	if tableName == "categorias" { tableName = "categories" }
	if tableName == "servicos" { tableName = "services" }

	// Verifica se tabela existe
	var exists bool
	db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&exists)
	if !exists {
		return
	}

	// Verifica se coluna unit_id existe
	var colExists bool
	db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = $1 AND column_name = 'unit_id')", tableName).Scan(&colExists)
	if !colExists {
		fmt.Printf("ℹ️  Tabela '%s' não possui coluna unit_id (global ou erro de schema).\n", tableName)
		return
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE unit_id IS NULL", tableName)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Printf("❌ Erro ao verificar isolamento em %s: %v\n", tableName, err)
		return
	}

	if count > 0 {
		fmt.Printf("🚨 ALERTA: %d registros em '%s' SEM unit_id (Possível vazamento global)\n", count, tableName)
	} else {
		fmt.Printf("Values OK: Tabela '%s' totalmente isolada (0 nulos).\n", tableName)
	}
}
