package main

import (
	"apps/ecosystem/cmd/api"
	"apps/ecosystem/core/migrations"
	"apps/ecosystem/tools/config"
	"apps/ecosystem/tools/db"
	"database/sql"
	"log"

	mySqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source"
)

func main() {
	// init db connection
	db, err := db.NewMariaDBStorage(mySqlCfg.Config{
		User: config.Env.DBUser,
		Passwd: config.Env.DBPassword,
		Addr: config.Env.DBAddress,
		DBName: config.Env.DBName,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
		MultiStatements: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	// todo is below necessary?
	defer db.Close()

	initStorage(db)
	runMigrations(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}

func runMigrations(db *sql.DB) {
	// Initialize the migration driver
	dbDriver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	source.Register("systemMigrations", &migrations.MigrationDriver{})
	sourceDriver := migrations.NewMigrationDriver()

	sourceDriver.Migrations, err = migrations.GetMigrations()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithInstance("systemMigrations", sourceDriver, "mariadb", dbDriver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Apply core migrations - should we check if we need to first?
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
