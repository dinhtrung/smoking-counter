package rest

import (
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/services"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/authjwt/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
	user, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	return c.JSON(h.svc.GetAll(user))
}

// Create a new Smoke Event for current user
func (h *SmokeRestAPI) Create(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	res, err := h.svc.Create(user, string(c.Body()))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete a new Smoke Event for current user
func (h *SmokeRestAPI) Delete(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(user, string(c.Body())); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
