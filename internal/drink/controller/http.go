package controller

import (
	"fmt"

	"github.com/labstack/echo"
)

type Drink struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

var ErrNotFound = fmt.Errorf("not found")

type storage interface {
	CreateDrink(Drink) (Drink, error)
	UpdateDrink(Drink) (Drink, error)
	DeleteDrink(name string) error
	DrinksByTags(tag []string) ([]Drink, error)
	AllDrinks() ([]Drink, error)
	DrinkByName(name string) (Drink, error)
}
type httpHandler struct {
	st   storage
	echo *echo.Echo
}

func NewHTTPHandler(st storage) *httpHandler {
	echo := echo.New()
	return &httpHandler{
		st:   st,
		echo: echo,
	}
}
func (h *httpHandler) Start(port string) error {
	return h.echo.Start(fmt.Sprintf(":%s", port))
}
func (h *httpHandler) Stop() error {
	return h.echo.Close()
}

// BuildRouter is a method that creates a new router group in the echo instance
// Example usage:
// h.BuildRouter("/api")
// This will create a new router group in the echo instance with the prefix /api
func (h *httpHandler) BuildRouter(group string) *echo.Group {
	router := h.echo.Group(group)
	return router
}

// AddRoutes is a method that adds the routes to the router
// Example usage:
func (h *httpHandler) AddRoutes(router *echo.Group) {
	router.POST("/drink", h.createDrink)
	router.PUT("/drink", h.updateDrink)
	router.DELETE("/drink/:name", h.deleteDrink)
	router.GET("/drink/tag/:tag", h.drinksByTags)
	router.GET("/drink", h.AllDrinks)
	router.GET("/drink/name/:name", h.DrinkByName)
}

func (h *httpHandler) createDrink(c echo.Context) error {
	var drink Drink
	if err := c.Bind(&drink); err != nil {
		return c.JSON(400, err.Error())
	}
	d, err := h.st.CreateDrink(drink)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) updateDrink(c echo.Context) error {
	var drink Drink
	if err := c.Bind(&drink); err != nil {
		return c.JSON(400, err.Error())
	}
	d, err := h.st.UpdateDrink(drink)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) deleteDrink(c echo.Context) error {
	name := c.Param("name")
	err := h.st.DeleteDrink(name)
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, "deleted")
}
func (h *httpHandler) drinksByTags(c echo.Context) error {
	tag := c.Param("tag")
	d, err := h.st.DrinksByTags([]string{tag})
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) AllDrinks(c echo.Context) error {
	d, err := h.st.AllDrinks()
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) DrinkByName(c echo.Context) error {
	name := c.Param("name")
	d, err := h.st.DrinkByName(name)
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
