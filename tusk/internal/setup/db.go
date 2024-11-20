package setup

import (
	"database/sql"
	"tusk/internal/config"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
)

func SetupDb(config *config.Config) (*sql.DB, error) {
	dbConn, err := otelsql.Open("postgres", config.DBConfig.ConnectionURI)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
