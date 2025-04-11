package main

import (
	"apps/ecosystem/core/migrations"
	"apps/ecosystem/tools/config"
	"apps/ecosystem/tools/db"
	"log"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewMariaDBStorage(mysqlCfg.Config{
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

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

// fix dirty database
// migrate -path PATH_TO_YOUR_MIGRATIONS -database YOUR_DATABASE_URL force VERSION
// migrate -path ./cmd/migrate/migrations/ -database mysql://david:GOmagpies!88@localhost/go_tutorial?sslmode=disable force 1
// migrate -path ./cmd/migrate/migrations/ -database mysql://test:test@localhost/go_tutorial?sslmode=disable force 1

// DROP DATABASE dave_erp;CREATE DATABASE dave_erp;GRANT ALL PRIVILEGES ON dave_erp.* TO 'david'@'localhost';FLUSH PRIVILEGES;USE dave_erp;