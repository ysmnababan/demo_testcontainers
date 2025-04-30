package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type handler struct {
	UserService IUserService
}

func NewUserHandler(db *gorm.DB) *handler {
	return &handler{
		UserService: NewUserService(db),
	}
}

func (h *handler) GetUsers(c echo.Context) (err error) {
	out, err := h.UserService.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK,
		map[string]interface{}{
			"data": out,
		})
}
