package migrations

import (
	"testing"

	"github.com/golang-migrate/migrate/v4/source"
	st "github.com/golang-migrate/migrate/v4/source/testing"
)

func TestMigrationDriver(t *testing.T) {
	s := &MigrationDriver{}
	d, err := s.Open("")
	if err != nil {
		t.Fatal(err)
	}

	m := source.NewMigrations()
	m.Append(&source.Migration{Version: 1, Identifier: "v1_up", Direction: source.Up})
	m.Append(&source.Migration{Version: 1, Identifier: "v1_down", Direction: source.Down})
	m.Append(&source.Migration{Version: 3, Identifier: "v3_up", Direction: source.Up})
	m.Append(&source.Migration{Version: 4, Identifier: "v4_up", Direction: source.Up})
	m.Append(&source.Migration{Version: 4, Identifier: "v4_down", Direction: source.Down})
	m.Append(&source.Migration{Version: 5, Identifier: "v5_down", Direction: source.Down})
	m.Append(&source.Migration{Version: 7, Identifier: "v7_up", Direction: source.Up})
	m.Append(&source.Migration{Version: 7, Identifier: "v7_down", Direction: source.Down})

	d.(*MigrationDriver).Migrations = m

	st.Test(t, d)
}
