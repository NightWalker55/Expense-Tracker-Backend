package handlers

import (
	"net/http"
	"time"

	"github.com/ayman/expense-tracker-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type ExpenseInput struct {
	Expense_Name string  `json:"expense_name"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"` // JSON will send date as string
}

func (r *Repository) CreateExpense(context *fiber.Ctx) error {
	var input ExpenseInput

	// 👇 Now parses date as string from JSON
	if err := context.BodyParser(&input); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// ✅ Parse string to time.Time
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid date format (expected YYYY-MM-DD)",
		})
	}

	expense := models.Expense{
		Expense_Name: input.Expense_Name,
		Amount:       input.Amount,
		Date:         parsedDate,
	}

	if err := r.DB.Create(&expense).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not create expense",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully created expense",
	})
}

func (r *Repository) DeleteExpense(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "ID cannot be empty",
		})
	}

	var expense models.Expense
	if err := r.DB.Where("id = ?", id).First(&expense).Error; err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Expense not found",
		})
	}

	if err := r.DB.Delete(&expense).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not delete expense",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Expense deleted successfully",
	})
}

func (r *Repository) GetDailyExpense(context *fiber.Ctx) error {
	dateStr := context.Query("date")
	if dateStr == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Date parameter is required",
		})
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	var expenses []models.Expense
	if err := r.DB.Where("date = ?", parsedDate).Find(&expenses).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not fetch expenses",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully fetched expenses",
		"data":    expenses,
	})
}

func (r *Repository) GetExpenseByID(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "ID cannot be empty",
		})
	}

	var expense models.Expense
	if err := r.DB.Where("id = ?", id).First(&expense).Error; err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Expense not found",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Expense fetched successfully",
		"data":    expense,
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/create_expense", r.CreateExpense)
	api.Delete("/delete_expense/:id", r.DeleteExpense)
	api.Get("/daily_expense", r.GetDailyExpense)
	api.Get("/daily_expense/:id", r.GetExpenseByID)
}
