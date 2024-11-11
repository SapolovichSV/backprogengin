package controller

import (
	"context"
	"net/http"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/labstack/echo"
)

type storage interface {
	CreateUser(context.Context, entities.User) (entities.User, error)
	UserByID(context.Context, int) (entities.User, error)
}
type httpHandler struct {
	st   storage
	echo *echo.Echo
	ctx  context.Context
}

func New(st storage, ctx context.Context) *httpHandler {
	e := echo.New()
	return &httpHandler{
		st:   st,
		echo: e,
		ctx:  ctx,
	}
}
func (h *httpHandler) AddRoutes(pathRoutesName string, router *echo.Router) {
	router.Add("POST", "/"+pathRoutesName+"/user", h.CreateUser)
	router.Add("GET", "/"+pathRoutesName+"/user/:id", h.UserByID)
}
func (h *httpHandler) CreateUser(c echo.Context) error {
	var user entities.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	user, err := h.st.CreateUser(h.ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}
func (h *httpHandler) UserByID(c echo.Context) error {
	var id int
	if err := c.Bind(&id); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	user, err := h.st.UserByID(h.ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}
