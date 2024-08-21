package routes

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qustavo/dotsql"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"microservice/globals"
)

var contract *openapi3.T

func TestMain(m *testing.M) {
	// load the variables found in the .env file into the process environment
	godotenv.Load()

	globals.Db, _ = pgxpool.New(context.Background(), "")
	globals.SqlQueries, _ = dotsql.LoadFromFile("./resources/queries.sql")

	var err error
	contract, err = openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestBasicHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/", BasicHandler)

	request := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)
}
