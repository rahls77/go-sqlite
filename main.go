package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rahls77/go-sqlite/database"
	"github.com/rahls77/go-sqlite/storeCredit"
)

func setUpRoutes(app *fiber.App) {

	app.Get("/api/v1/storeCredit/:id", storeCredit.GetStoreCredit)
	app.Post("/api/v1/storeCredit", storeCredit.NewStoreCredit)
	app.Post("/api/v1/webhook", storeCredit.NewWebHook)
}

func initDatabase(filename string) *gorm.DB {
	var err error

	database.DBConn, err = gorm.Open("sqlite3", filename)

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection opened to Database")
	database.DBConn.AutoMigrate(&storeCredit.StoreCredit{})
	fmt.Println("Database Migrated")
	return database.DBConn
}

func main() {
	app := fiber.New()

		// Default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
			AllowOrigins: "https://rahulsharma-bootcamp.myshopify.com",
			AllowHeaders:  "Origin, Content-Type, Accept",
	}))

	initDatabase("storeCredit.db")

	setUpRoutes(app)

	app.Listen(":3000")

	defer database.DBConn.Close()

}
