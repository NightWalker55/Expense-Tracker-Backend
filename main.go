package main

import (
	"log"
	"os"

	"github.com/ayman/expense-tracker-backend/db"
	"github.com/ayman/expense-tracker-backend/handlers"
	"github.com/ayman/expense-tracker-backend/models"
	"github.com/gofiber/fiber/v2"
	
)

type Expense struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	//Date       string    `json:"date"`
	//Created_At time.Time `json:"created_at"`
}

func main() {
	/*err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}*/

	app := fiber.New()

	config := &db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := db.NewConnection(config)

	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateExpense(db)

	if err != nil {
		log.Fatal("Could not migrate DB")
	}

	r := handlers.Repository{
		DB: db,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Expense Tracker API!")
	})

	r.SetupRoutes((app))

	log.Println("Server running on port 8080")
	log.Fatal(app.Listen(":8080"))

}
