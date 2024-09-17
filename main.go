package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Name     string `validate:"required,min=3,max=32" error:"Name must be between 3 and 32 characters long"`
	Age      int    `validate:"required" error:"age is required"`
	Password string `validate:"required,min=8" error:"Password must be at least 8 characters long"`
}

func (m *User) Validate(validate *validator.Validate) ([]ValidationError, error) {
	return ValidateFunc[User](*m, validate)
}

type Response struct {
	Status int               `json:"status"`
	Errors []ValidationError `json:"errors"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		newUser := User{
			Name:     "J",
			Age:      30,
			Password: "p",
		}

		validate := validator.New()

		validationErrors, err := newUser.Validate(validate)
		if err != nil {
			response := Response{
				Status: 400,
				Errors: validationErrors,
			}

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":8002")
}
