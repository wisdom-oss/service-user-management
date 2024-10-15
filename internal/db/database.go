package db

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This file contains the connection to the database which is automatically
// initialized on import/app startup

// Pool is automatically initialized at the app startup using the init
// function in the internal package
var Pool *pgxpool.Pool

// StateDB contains all states currently not used to retreive a access token
var StateDB *badger.DB
