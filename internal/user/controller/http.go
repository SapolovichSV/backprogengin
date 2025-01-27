package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/labstack/echo/v4"
)

type storage interface {
	CreateUser(context.Context, entities.User) (entities.User, error)
	UserByID(context.Context, int) (entities.User, error)
	AddFav(ctx context.Context, drinkName string, userID int) (entities.User, error)
}
type authService interface {
	Auth(c echo.Context) (entities.User, error)
	Login(c echo.Context) (entities.User, error)
	Register(c echo.Context, user entities.User) error
}
type httpHandler struct {
	st   storage
	echo *echo.Echo
	ctx  context.Context
	auth authService
}

func New(st storage, auth authService, ctx context.Context) *httpHandler {
	e := echo.New()
	return &httpHandler{
		st:   st,
		echo: e,
		ctx:  ctx,
		auth: auth,
	}
}

func (h *httpHandler) AddRoutes(pathRoutesName string, router *echo.Router) {
	router.Add("POST", "/"+pathRoutesName+"/user", h.CreateUser)
	router.Add("GET", "/"+pathRoutesName+"/user/:id", h.UserByID)
	router.Add("PATCH", "/"+pathRoutesName+"/user/fav", h.AddFav)
	router.Add("GET", "/"+pathRoutesName+"/user/login", h.Login)
}

// CreateUser godoc
// @Summary Create a user
// @Decsription assigment to user cookie(jwt token) wih user info
// @Description field id will be ignored
// @Description id will be in response
// @Description Create a user,with his favourite drinks(optional),if such drinks non-existent: error,
// @Description otherwise return created user
// @Tags user
// @Accept json
// @Produce json
// @Param user body entities.User true "User object"
// @Success 201 {object} entities.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user [post]
func (h *httpHandler) CreateUser(c echo.Context) error {
	var user entities.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.st.CreateUser(h.ctx, user)
	fmt.Println(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := h.auth.Register(c, user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login
// @Description Login user,if user non-existent: error
// @Description otherwise return user info
// @Description assigment to user cookie(encoded jwt token) wih user info
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} entities.User
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /user/login [get]
func (h *httpHandler) Login(c echo.Context) error {
	user, err := h.auth.Login(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	user, err = h.st.UserByID(h.ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

// UserByID godoc
// @Summary Get user
// @Description Get user by ID(which contains in cookie: jwt token)
// @Tags user
// @Accept plain
// @Produce json
// @Success 200 {object} entities.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/{id} [get]
func (h *httpHandler) UserByID(c echo.Context) error {
	userInfo, err := h.auth.Auth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	user, err := h.st.UserByID(h.ctx, userInfo.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

// AddFav godoc
// @Summary Add a favourite drink to user
// @Description Add a favourite drink to user
// @Tags user
// @Accept json
// @Produce json
// @Param drinkname path string true "Drink name"
// @Success 202 {object} entities.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/fav [patch]
func (h *httpHandler) AddFav(c echo.Context) error {

	drinkName := c.Param("drinkname")
	fmt.Println(drinkName)
	if len(drinkName) == 0 {
		drinkName = c.QueryParam("drinkname")
	}
	userInfo, err := h.auth.Auth(c)
	fmt.Println(userInfo)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	user, err := h.st.AddFav(h.ctx, drinkName, userInfo.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusAccepted, user)
}
