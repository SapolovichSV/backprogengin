package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
	"github.com/labstack/echo/v4"
)

type storage interface {
	CreateDrink(context.Context, entities.Drink) (entities.Drink, error)
	UpdateDrink(context.Context, entities.Drink) (entities.Drink, error)
	DeleteDrink(ctx context.Context, name string) error
	DrinksByTags(ctx context.Context, tag []string) ([]entities.Drink, error)
	AllDrinks(ctx context.Context, id int) ([]entities.Drink, error)
	DrinkByName(ctx context.Context, name string) (entities.Drink, error)
}

var ErrNotFound = fmt.Errorf("not found")

type httpHandler struct {
	st   storage
	echo *echo.Echo
	ctx  context.Context
}

func New(st storage, ctx context.Context) *httpHandler {
	echo := echo.New()
	return &httpHandler{
		st:   st,
		echo: echo,
		ctx:  ctx,
	}
}

// BuildRouter is a method that creates a new router group in the echo instance
// Example usage:
// h.BuildRouter("/api")
// This will create a new router group in the echo instance with the prefix /api

// AddRoutes is a method that adds the routes to the router
// Example usage:
// h.AddRoutes(router)
// This will add the routes to the router
// The routes are:
// POST /{{pathRoutesName}}/drink
// and e.t.c
func (h *httpHandler) AddRoutes(pathRoutesName string, router *echo.Router) {
	router.Add("POST", "/"+pathRoutesName+"/drink", h.createDrink)
	//router.POST("/drink", h.createDrink)
	router.Add("PUT", "/"+pathRoutesName+"/drink", h.updateDrink)
	//router.PUT("/drink", h.updateDrink)
	router.Add("DELETE", "/"+pathRoutesName+"/drink/:name", h.deleteDrink)
	//router.DELETE("/drink/:name", h.deleteDrink)
	router.Add("GET", "/"+pathRoutesName+"/drink/tag/:tag", h.drinksByTags)
	//router.GET("/drink/tag/:tag", h.drinksByTags)
	router.Add("GET", "/"+pathRoutesName+"/drink/id/:id", h.allDrinks)
	//router.GET("/drink/id/:id", h.allDrinks)
	router.Add("GET", "/"+pathRoutesName+"/drink/name/:name", h.drinkByName)
	//router.GET("/drink/name/:name", h.drinkByName)
}

func (h *httpHandler) createDrink(c echo.Context) error {
	var drink entities.Drink
	if err := c.Bind(&drink); err != nil {
		return c.JSON(400, err.Error())
	}
	d, err := h.st.CreateDrink(h.ctx, drink)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(http.StatusCreated, d)
}
func (h *httpHandler) updateDrink(c echo.Context) error {
	var drink entities.Drink
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
	param := c.Param("id")
	id, err := strconv.Atoi(param)
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
