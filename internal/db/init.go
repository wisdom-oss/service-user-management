package db

import (
	"context"
	"io/fs"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qustavo/dotsql"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"microservice/resources"
)

func init() {
	l := log.With().Str("package", "internal/db").Logger()
	l.Debug().Msg("connecting to the database")

	var err error
	Pool, err = pgxpool.New(context.Background(), "")
	if err != nil {
		l.Fatal().Err(err).Msg("could not connect to database")
	}
	err = Pool.Ping(context.Background())
	if err != nil {
		l.Fatal().Err(err).Msg("could not ping database")
	}
	l.Debug().Msg("connected to the database")

	l.Debug().Msg("loading prepared sql queries")
	files, err := fs.ReadDir(resources.QueryFiles, ".")
	if err != nil {
		l.Fatal().Err(err).Msg("could not load queries")
	}
	var instances []*dotsql.DotSql
	for _, queryFile := range files {
		fd, err := resources.QueryFiles.Open(queryFile.Name())
		if err != nil {
			l.Fatal().Err(err).Msg("could not open query file")
		}
		instance, err := dotsql.Load(fd)
		if err != nil {
			l.Fatal().Err(err).Msg("could not load query file")
		}
		instances = append(instances, instance)
	}
	Queries = dotsql.Merge(instances...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load prepared queries")
	}

	redisUri, isSet := os.LookupEnv("REDIS_URI")
	if !isSet {
		l.Fatal().Msg("REDIS_URI is not set")
	}
	redisOptions, err := redis.ParseURL(redisUri)
	if err != nil {
		l.Fatal().Err(err).Msg("could not parse redis URI")
	}

	Redis = redis.NewClient(redisOptions)
}
