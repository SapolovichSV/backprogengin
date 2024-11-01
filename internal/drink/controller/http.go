package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/labstack/echo"
)

type Drink struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

var ErrNotFound = fmt.Errorf("not found")

type storage interface {
	CreateDrink(context.Context, Drink) (Drink, error)
	UpdateDrink(context.Context, Drink) (Drink, error)
	DeleteDrink(ctx context.Context, name string) error
	DrinksByTags(ctx context.Context, tag []string) ([]Drink, error)
	AllDrinks(ctx context.Context, id int) ([]Drink, error)
	DrinkByName(ctx context.Context, name string) (Drink, error)
}

type httpHandler struct {
	st   storage
	echo *echo.Echo
	ctx  context.Context
}

func NewHTTPHandler(st storage, ctx context.Context) *httpHandler {
	echo := echo.New()
	return &httpHandler{
		st:   st,
		echo: echo,
		ctx:  ctx,
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
	router.GET("/drink", h.allDrinks)
	router.GET("/drink/name/:name", h.drinkByName)
}

func (h *httpHandler) createDrink(c echo.Context) error {
	var drink Drink
	if err := c.Bind(&drink); err != nil {
		return c.JSON(400, err.Error())
	}
	d, err := h.st.CreateDrink(h.ctx, drink)
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
	d, err := h.st.UpdateDrink(h.ctx, drink)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) deleteDrink(c echo.Context) error {
	name := c.Param("name")
	err := h.st.DeleteDrink(h.ctx, name)
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, "deleted")
}
func (h *httpHandler) drinksByTags(c echo.Context) error {
	tag := c.Param("tag")
	d, err := h.st.DrinksByTags(h.ctx, []string{tag})
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) allDrinks(c echo.Context) error {
	id, err := strconv.Atoi(
		c.Param("id"))
	if err != nil {
		return c.JSON(500, err.Error())
	}
	d, err := h.st.AllDrinks(h.ctx, id)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
func (h *httpHandler) drinkByName(c echo.Context) error {
	name := c.Param("name")
	d, err := h.st.DrinkByName(h.ctx, name)
	if err == ErrNotFound {
		return c.JSON(404, echo.ErrNotFound.Error())
	} else if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}
