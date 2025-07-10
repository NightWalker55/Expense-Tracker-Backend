package handlers

import (
	"net/http"
	"time"

	"strconv"

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
	Date         string  `json:"date"`
}

func (r *Repository) CreateExpense(context *fiber.Ctx) error {
	var input ExpenseInput

	if err := context.BodyParser(&input); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

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

func (r *Repository) GetAllExpense(c *fiber.Ctx) error {
	var expenses []models.Expense

	if err := r.DB.Find(&expenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch expenses",
		})
	}

	return c.Status(fiber.StatusOK).JSON(expenses)
}

func (r *Repository) UpdateExpense(context *fiber.Ctx) error {
	idParam := context.Params("id")
	if idParam == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "ID cannot be empty",
		})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid ID format",
		})
	}

	var input ExpenseInput
	if err := context.BodyParser(&input); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	if input.Expense_Name == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Expense name is required",
		})
	}
	if input.Amount <= 0 {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Amount must be greater than 0",
		})
	}

	var expense models.Expense
	if err := r.DB.First(&expense, uint(id)).Error; err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Expense not found",
		})
	}

	expense.Expense_Name = input.Expense_Name
	expense.Amount = input.Amount

	if input.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Invalid date format (expected YYYY-MM-DD)",
			})
		}
		if parsedDate.After(time.Now()) {
			return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Date cannot be in the future",
			})
		}
		expense.Date = parsedDate
	}

	if err := r.DB.Save(&expense).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not update expense",
			"error":   err.Error(),
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/create_expense", r.CreateExpense)
	api.Delete("/delete_expense/:id", r.DeleteExpense)
	api.Get("/daily_expense", r.GetDailyExpense)
	api.Get("/daily_expense/:id", r.GetExpenseByID)
	api.Patch("/update_expense/:id", r.UpdateExpense)
	api.Get("/all_daily_expense", r.GetAllExpense)
}
