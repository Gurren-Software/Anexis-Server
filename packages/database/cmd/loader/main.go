package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-go-sdk/atlasexec"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/Gurren-Software/Anexis-Server/packages/database/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.File{},
		&models.Link{},
		&models.MigrationJob{},
		&models.BackupJob{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}

// This is used by Atlas for schema inspection
var _ atlasexec.MigrateApplyParams
