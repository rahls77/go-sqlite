package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rahls77/go-sqlite/book"
	"github.com/rahls77/go-sqlite/database"
)

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func setUpRoutes(app *fiber.App) {
	app.Get("/", helloWorld)

	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
}

func initDatabase() {
	var err error

	database.DBConn, err = gorm.Open("sqlite3", "books.db")

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection opened to Database")
	database.DBConn.AutoMigrate(&book.Book{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()

	initDatabase()

	setUpRoutes(app)

	app.Listen(3000)

	defer database.DBConn.Close()

}
