package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/labstack/echo/v4"
)

type storage interface {
	CreateUser(context.Context, entities.User) (entities.User, error)
	UserByID(context.Context, int) (entities.User, error)
	AddFav(ctx context.Context, drinkName string, userID int) (entities.User, error)
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
	router.Add("PATCH", "/"+pathRoutesName+"/user/fav", h.AddFav)
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
	id, err := getParamId("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	user, err := h.st.UserByID(h.ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

func getParamId(paramName string, c echo.Context) (int, error) {
	param := c.Param(paramName)
	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (h *httpHandler) AddFav(c echo.Context) error {
	type inp struct {
		DrinkName string `json:"drinkname"`
		ID        int    `json:"id"`
	}
	var input inp
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	user, err := h.st.AddFav(h.ctx, input.DrinkName, input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusAccepted, user)
}
