package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// This file contains the connection to the database which is automatically
// initialized on import/app startup

// Pool is automatically initialized at the app startup using the init
// function in the internal package
var Pool *pgxpool.Pool

// Redis holds the connection to the redis server used to distribute avavilable
// states between service instances
var Redis *redis.Client
