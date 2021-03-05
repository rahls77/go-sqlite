package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupApp() *fiber.App {
	app := fiber.New()
	setUpRoutes(app)
	return app
}

func removeDbFile(filename string) {
	e := os.Remove(filename)
	if e != nil {
		errors.New("Couldnt remove file")
	}
}

func Test_get_fail(t *testing.T) {
	app := setupApp()
	initDatabase("storeCreditTest.db")

	req := httptest.NewRequest("GET", "/api/v1/storeCredit/1", nil)

	res, _ := app.Test(req)

	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusNotFound, res.StatusCode)
	}
}

func Test_get_pass(t *testing.T) {
	app := setupApp()
	db := initDatabase("storeCreditTest.db")

	db.Exec("INSERT into store_credits (customer_id, credit, shop_domain) values(1,1,'example.com')")
	req := httptest.NewRequest("GET", "/api/v1/storeCredit/1", nil)

	res, _ := app.Test(req)

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}

	removeDbFile("storeCreditTest.db")
}

func Test_post_pass(t *testing.T) {
	app := setupApp()
	initDatabase("storeCreditTest.db")

	// db.Exec("INSERT into store_credits (customer_id, credit, shop_domain) values(1,1,'example.com')")
	values := map[string]interface{}{"credit": 1, "shopDomain": "xyz"}
	jsonValue, _ := json.Marshal(values)

	req := httptest.NewRequest("POST", "/api/v1/storeCredit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
	removeDbFile("storeCreditTest.db")

}
