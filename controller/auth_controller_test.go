package controller_test

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/repository"
	"auth-service/service"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Setup() *fiber.App {
	app := fiber.New()
	configApp := config.New(`./../.env.test`)
	db := config.NewPostgresDB(configApp)
	redisConfig := config.NewRedis(configApp)
	userRepository := repository.NewUserRepository(db)
	redisService := service.NewRedisConfig(redisConfig, 30)
	authService := service.NewAuthService(userRepository, redisService)
	authController := controller.NewAuthController(authService)
	authController.Route(app)
	return app
}

var app = Setup()

func TestLoginEmptyEmail(t *testing.T) {
	requestBody := strings.NewReader(`{
		"email":"",
		"password":"password"
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "BAD_REQUEST", parse["status"])
	assert.Equal(t, float64(400), parse["response_code"])
	data := response["data"].([]interface{})
	for _, val := range data {
		value := val.(map[string]interface{})
		assert.Equal(t, "email", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestLoginEmptyPassword(t *testing.T) {
	requestBody := strings.NewReader(`{
		"email":"adit@mail.com",
		"password":""
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "BAD_REQUEST", parse["status"])
	assert.Equal(t, float64(400), parse["response_code"])
	data := response["data"].([]interface{})
	for _, val := range data {
		value := val.(map[string]interface{})
		assert.Equal(t, "password", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestWrongEmailInput(t *testing.T) {
	requestBody := strings.NewReader(`{
		"email":"adit@",
		"password":"121242352435"
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "BAD_REQUEST", parse["status"])
	assert.Equal(t, float64(400), parse["response_code"])
	data := response["data"].([]interface{})
	for _, val := range data {
		value := val.(map[string]interface{})
		assert.Equal(t, "email", value["field"])
		assert.Equal(t, "format email salah", value["message"])
	}
}

func TestWrongEmailOrPassword(t *testing.T) {
	requestBody := strings.NewReader(`{
		"email":"adit@mail.com",
		"password":"121242352435"
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "BAD_REQUEST", parse["status"])
	assert.Equal(t, float64(400), parse["response_code"])
	assert.Equal(t, parse["data"], "")
	assert.Equal(t, response["message"], "Email Or Password Incorect")
}

func TestSuccessLogin(t *testing.T) {
	requestBody := strings.NewReader(`{
		"email":"adit@mail.com",
		"password":"password"
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", parse["status"])
	assert.Equal(t, float64(200), parse["response_code"])
	assert.Equal(t, response["message"], "Success")

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["access_token"])
	assert.NotNil(t, data["refresh_token"])
}
