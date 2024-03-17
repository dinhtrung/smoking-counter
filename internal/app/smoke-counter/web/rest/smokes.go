package rest

import (
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SmokeRestAPI struct {
	svc services.SmokeService
}

func NewSmokeRestAPI(svc services.SmokeService) *SmokeRestAPI {
	return &SmokeRestAPI{
		svc: svc,
	}
}

// GetAll get a user
func (h *SmokeRestAPI) GetAll(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	sub, err := user.Claims.GetSubject()
	if err != nil {
		return err
	}
	return c.JSON(h.svc.GetAll(sub))
}

// Create a new Smoke Event for current user
func (h *SmokeRestAPI) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	sub, err := user.Claims.GetSubject()
	if err != nil {
		return err
	}
	res, err := h.svc.Create(sub, string(c.Body()))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
