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

// CreateUser godoc
// @Summary Create a user
// @Description field id will be ignored
// @Description id will be in response
// @Description Create a user,with his favourite drinks(optional),if such drinks non-existent: error,
// @Description otherwise return created user
// @Tags user
// @Accept json
// @Produce json
// @Param user body entities.User true "User object"
// @Success 201 {object} entities.User
// @Router /user [post]
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

// UserByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Accept plain
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/{id} [get]
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

// AddFav godoc
// @Summary Add a favourite drink to user
// @Description Add a favourite drink to user
// @Tags user
// @Accept json
// @Produce json
// @Param drinkname path string true "Drink name"
// @Param id path int true "User ID"
// @Success 202 {object} entities.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/fav [patch]
func (h *httpHandler) AddFav(c echo.Context) error {

	drinkName := c.Param("drinkname")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	user, err := h.st.AddFav(h.ctx, drinkName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusAccepted, user)
}
