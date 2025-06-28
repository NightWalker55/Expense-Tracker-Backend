package handlers

import (
	"net/http"

	"github.com/ayman/expense-tracker-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateExpense(context *fiber.Ctx) error {
	var expense models.Expense
	if err := context.BodyParser(&expense); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Request failed",
		})
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
	var expenses []models.Expense
	if err := r.DB.Find(&expenses).Error; err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not get expenses",
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
