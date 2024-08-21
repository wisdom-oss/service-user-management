package globals

import (
	_ "embed"

	"github.com/qustavo/dotsql"
)

// This file contains globally shared variables (e.g., service name, sql queries)

// ServiceName contains the service's identifying name.
const ServiceName = "template" //TODO: CHANGE THIS VALUE

// SqlQueries contains the prepared sql queries from the resources folder
var SqlQueries *dotsql.DotSql
