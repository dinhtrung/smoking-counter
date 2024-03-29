package rest

import (
	"fmt"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/authjwt/utils"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/services"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/shared"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserResource interface {
	GetAllUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetAuthorities(c *fiber.Ctx) error
}

type DefaultUserResource struct {
	Svc  services.UserService
	Repo services.UserRepository
}

func NewDefaultUserResource(svc services.UserService, repo services.UserRepository) UserResource {
	return &DefaultUserResource{
		Svc:  svc,
		Repo: repo,
	}
}

func (r *DefaultUserResource) GetAllUser(c *fiber.Ctx) error {
	rows, err := r.Repo.FindAll()
	if err != nil {
		return err
	}
	cnt, err := r.Repo.Count()
	if err != nil {
		return err
	}
	c.Set("X-Total-Count", fmt.Sprint(cnt))
	return c.JSON(rows)
}

func (r *DefaultUserResource) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id parameter")
	}
	user, err := r.Repo.FindById(id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (r *DefaultUserResource) CreateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	if r.Repo.ExistsByLogin(user.Login) {
		return fiber.ErrConflict
	}
	if strings.TrimSpace(user.Password) == "" {
		user.Password = utils.RandomString(8)
		c.Response().Header.Add("X-Password", user.Password)
	}
	hash, err := r.Svc.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	hasRoleUser := false
	for _, r := range user.Authorities {
		if r == "ROLE_USER" {
			hasRoleUser = true
			break
		}
	}
	if !hasRoleUser {
		user.Authorities = append(user.Authorities, "ROLE_USER")
	}

	user.CreatedDate = time.Now().Format(time.RFC3339)
	if currentLogin, err := utils.GetCurrentUserLogin(c); err == nil {
		user.CreatedBy = currentLogin
	}
	if err := r.Repo.Save(&user); err != nil {
		return err
	}
	return c.JSON(user.UserDTO)
}

func (r *DefaultUserResource) UpdateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	exists, err := r.Repo.FindByLogin(user.Login)
	if err != nil {
		return fiber.ErrNotFound
	}
	if strings.TrimSpace(user.Password) == "" {
		user.Password = utils.RandomString(8)
		c.Response().Header.Add("X-Password", user.Password)
	}
	hash, err := r.Svc.HashPassword(user.Password)
	if err != nil {
		return err
	}

	exists.Password = hash
	exists.UserDTO = user.UserDTO

	user.LastModifiedDate = time.Now().Format(time.RFC3339)
	if currentLogin, err := utils.GetCurrentUserLogin(c); err == nil {
		user.LastModifiedBy = currentLogin
	}

	if err := r.Repo.Save(&exists); err != nil {
		return err
	}
	return c.JSON(user.UserDTO)
}

func (r *DefaultUserResource) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}
	if err := r.Repo.DeleteById(id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (r *DefaultUserResource) GetAuthorities(c *fiber.Ctx) error {
	authorities, err := r.Repo.FindAllAuthorities()
	if err != nil {
		return err
	}
	return c.JSON(authorities)
}
