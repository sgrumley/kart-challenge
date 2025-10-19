package seeding

import (
	"embed"
	"io/fs"
	"log"

	"github.com/jmoiron/sqlx"
)

//go:embed scripts/*.sql
var sql embed.FS

func InsertAll(db *sqlx.DB) {
	var sqlFiles []string

	files, err := fs.ReadDir(sql, "scripts")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := sql.ReadFile("scripts/" + file.Name())
		if err != nil {
			log.Fatal("error: ", file.Name(), ": ", err)
		}
		sqlFiles = append(sqlFiles, string(data))
	}

	for _, query := range sqlFiles {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal("error: ", query, ": \n", err)
		}
	}
}
