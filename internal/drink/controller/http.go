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

// createDrink godoc
//
//		@Summary Creates a drink
//		@Description Creates a drink with the specified name and tags
//		@Tags drink
//		@Accept json
//		@Produce json
//		@Success 201 {object} entities.Drink
//	 	@Failure 500 {string} string
//		@Param drink body entities.Drink true "Drink what we add with optional tags,if tags not: set tags will be empty, name is required,"
//		@Router /drink [post]
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

// updateDrink godoc
//
//		@Summary Updates drink tags
//		@Description Updates drink tags with the specified name(old tags will be deleted)
//		@Tags drink
//		@Accept json
//		@Produce json
//		@Success 200 {object} entities.Drink
//		@Failure 500 {string} string
//		@Param drink body entities.Drink true "Drink what we update with optional tags,if tags not: set tags will be empty, name is required,"
//	 @Router /drink [put]
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

// deleteDrink godoc
//
//		@Summary Deletes a drink
//		@Description Deletes a drink with the specified name,other fields will be ignored
//		@Tags drink
//	 @Accept plain
//		@Produce json
//		@Success 200 {string} deleted
//		@Failure 404 {string} string
//		@Failure 500 {string} string
//		@Param name	path string	true "Name of the drink to delete"
//		@Router /drink/{name} [delete]
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

// drinksByTags godoc
// @Summary Get drinks by tags
// @Description Get drinks by tags
//
//	@Tags drink
//
// @Accept plain
// @Produce json
// @Success 200 {array} entities.Drink
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Param tag path string true "tasty sweet spicy"
// @Router /drink/tag/{tag} [get]
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

// allDrinks godoc
// @Summary Get all drinks
// @Description Get all drinks with offset = id
//
//	@Tags drink
//
// @Accept plain
// @Produce json
// @Success 200 {array} entities.Drink
// @Failure 500 {string} string "Internal server error"
// @Param id path int true "id"
// @Router /drink/id/{id} [get]
func (h *httpHandler) allDrinks(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, err.Error())
	}
	d, err := h.st.AllDrinks(h.ctx, id)
	if err != nil {
		fmt.Println(err.Error() + "at storage")
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, d)
}

// drinkByName godoc
// @Summary Get drink by name
// @Description Get drink by name
//
//	@Tags drink
//
// @Accept plain
// @Produce json
// @Success 200 {object} entities.Drink
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Param name path string true "Name of the drink"
// @Router /drink/name/{name} [get]
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
