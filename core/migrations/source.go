package migrations

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigrationDriver() *MigrationDriver {
    return &MigrationDriver{Migrations: source.NewMigrations()}
}

type MigrationDriver struct {
	Migrations *source.Migrations
}

func (p *MigrationDriver) Open(url string) (source.Driver, error) {
	return p, nil
}

func (p *MigrationDriver) Close() error {
	return nil
}

func (p *MigrationDriver) First() (version uint, err error) {
	if first, ok := p.Migrations.First(); !ok {
        return 0, &os.PathError{Op: "first", Err: os.ErrNotExist}
    } else {
        return first, nil
    }
}

func (p *MigrationDriver) Prev(version uint) (prevVersion uint, err error) {
    if prev, ok := p.Migrations.Prev(version); !ok {
        return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Err: os.ErrNotExist}
    } else {
        return prev, nil
    }
}

func (p *MigrationDriver) Next(version uint) (nextVersion uint, err error) {
    if next, ok := p.Migrations.Next(version); !ok {
        return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Err: os.ErrNotExist}
    } else {
        return next, nil
    }
}

func (p *MigrationDriver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
    if m, ok := p.Migrations.Up(version); ok {
        return io.NopCloser(bytes.NewBufferString(m.Raw)), m.Identifier, nil
    }
    return nil, "", &os.PathError{Op: fmt.Sprintf("read up version %v", version), Err: os.ErrNotExist}
}

func (p *MigrationDriver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
    if m, ok := p.Migrations.Down(version); ok {
        return io.NopCloser(bytes.NewBufferString(m.Raw)), m.Identifier, nil
    }
    return nil, "", &os.PathError{Op: fmt.Sprintf("read down version %v", version), Err: os.ErrNotExist}
}
