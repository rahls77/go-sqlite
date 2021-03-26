package storeCredit

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rahls77/go-sqlite/database"
)

type StoreCredit struct {
	gorm.Model
	CustomerId int    `json:"customerId"`
	Credit     int    `json:"credit"`
	ShopDomain string `json:"shopDomain"`
}

type WebHook struct {
	OrderStatusUrl string `json:"order_status_url"`
	Customer       struct {
		CustomerId int `json:"id"`
	} `json:"customer"`
}

func GetStoreCredit(c *fiber.Ctx) error {
	customerId, err := strconv.Atoi(c.Params("id"))
	fmt.Println(err)

	db := database.DBConn
	var storeCredit StoreCredit
	db.Find(&storeCredit, StoreCredit{CustomerId: customerId})
	fmt.Println("Called store credit get request")
	fmt.Println(storeCredit)
	if storeCredit.CustomerId == 0 {
		return c.Status(404).SendString("Not found")
	}
	return c.Status(200).JSON(storeCredit)
}

func NewStoreCredit(c *fiber.Ctx) error {
	db := database.DBConn
	storeCredit := new(StoreCredit)
	fmt.Println(c.BodyParser(storeCredit))
	if err := c.BodyParser(storeCredit); err != nil {
		fmt.Println(err)
		return c.Status(503).SendString(error.Error(err))
	}
	db.Create(&storeCredit)
	return c.JSON(storeCredit)
}

func NewWebHook(c *fiber.Ctx) error {
	orderStatus := WebHook{}
	textBytes := []byte(c.Body())
	err := json.Unmarshal(textBytes, &orderStatus)
	if err != nil {
		fmt.Println(err)
	}
	orderStatus.OrderStatusUrl = strings.Split(orderStatus.OrderStatusUrl, "/")[2]

	var storeCredit StoreCredit

	db := database.DBConn
	db.Where(StoreCredit{CustomerId: orderStatus.Customer.CustomerId}).Attrs(StoreCredit{ShopDomain: orderStatus.OrderStatusUrl}).FirstOrCreate(&storeCredit)
	newCredit := storeCredit.Credit + 1
	db.Model(&storeCredit).Where(StoreCredit{CustomerId: orderStatus.Customer.CustomerId}).Updates(StoreCredit{Credit: newCredit})
	s, _ := json.MarshalIndent(storeCredit, "", "\t");
	fmt.Println(string(s))
	fmt.Println("Called Webhook")
	return c.JSON(storeCredit)
}
