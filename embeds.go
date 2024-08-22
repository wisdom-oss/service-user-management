package main

import "embed"

//go:embed resources/migrations/*.sql
var DatabaseMigrations embed.FS
