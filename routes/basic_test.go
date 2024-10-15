package routes

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	validator "openapi.tanna.dev/go/validator/openapi3"

	_ "microservice/internal/db"
)

var contract *openapi3.T

func TestMain(m *testing.M) {
	var err error
	contract, err = openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestBasicHandler(t *testing.T) {
	router := gin.New()
	router.GET("/", BasicHandler)

	request := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)
}
