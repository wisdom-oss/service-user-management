package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// This file contains the connection to the database which is automatically
// initialized on import/app startup

// Pool is automatically initialized at the app startup using the init
// function in the internal package
var Pool *pgxpool.Pool
