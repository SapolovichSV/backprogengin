package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
	"github.com/SapolovichSV/backprogeng/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_httpHandler_createDrink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		reqBody      entities.Drink
		respBody     entities.Drink
		hasRespBody  bool
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodPost,
			path:         "/drink",
			reqBody:      entities.Drink{Name: "test01"},
			respBody:     entities.Drink{Name: "test01"},
			hasRespBody:  true,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "test02",
			method:       http.MethodPost,
			path:         "/drink",
			reqBody:      entities.Drink{Name: "test02", Tags: []string{"spicy", "non-alcohol"}},
			respBody:     entities.Drink{Name: "test02", Tags: []string{"spicy", "non-alcohol"}},
			hasRespBody:  true,
			expectedCode: http.StatusCreated,
		},
	}

	mockStorage.EXPECT().CreateDrink(gomock.Any(), ts[0].reqBody).Return(ts[0].respBody, nil)
	mockStorage.EXPECT().CreateDrink(gomock.Any(), ts[1].reqBody).Return(ts[1].respBody, nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		reqBody, _ := json.Marshal(v.reqBody)
		req := httptest.NewRequest(v.method, v.path, strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)

		if assert.NoError(t, h.createDrink(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
			if v.hasRespBody {
				resData, _ := json.Marshal(&v.respBody)
				assert.JSONEq(t, string(resData), rec.Body.String())
			}
		}
	}
}

func Test_httpHandler_updateDrink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		reqBody      entities.Drink
		respBody     entities.Drink
		hasRespBody  bool
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodPut,
			path:         "/drink",
			reqBody:      entities.Drink{Name: "test01"},
			respBody:     entities.Drink{Name: "test01"},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "test02",
			method:       http.MethodPut,
			path:         "/drink",
			reqBody:      entities.Drink{Name: "test02", Tags: []string{"spicy", "non-alcohol"}},
			respBody:     entities.Drink{Name: "test02", Tags: []string{"spicy", "non-alcohol"}},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
		},
	}

	mockStorage.EXPECT().UpdateDrink(gomock.Any(), ts[0].reqBody).Return(ts[0].respBody, nil)
	mockStorage.EXPECT().UpdateDrink(gomock.Any(), ts[1].reqBody).Return(ts[1].respBody, nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		reqBody, _ := json.Marshal(v.reqBody)
		req := httptest.NewRequest(v.method, v.path, strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)

		if assert.NoError(t, h.updateDrink(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
			if v.hasRespBody {
				resData, _ := json.Marshal(&v.respBody)
				assert.JSONEq(t, string(resData), rec.Body.String())
			}
		}
	}
}

func Test_httpHandler_deleteDrink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodDelete,
			path:         "/drink/:name",
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "name",
					value: "test01",
				},
			},
		},
		{
			name:         "test02",
			method:       http.MethodDelete,
			path:         "/drink/:name",
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "name",
					value: "test02",
				},
			},
		},
	}

	mockStorage.EXPECT().DeleteDrink(gomock.Any(), "test01").Return(nil)
	mockStorage.EXPECT().DeleteDrink(gomock.Any(), "test02").Return(nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		req := httptest.NewRequest(v.method, v.path, nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)
		for _, v := range v.params {
			c.SetParamNames(v.name)
			c.SetParamValues(v.value)
		}

		if assert.NoError(t, h.deleteDrink(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
		}
	}
}

func Test_httpHandler_drinksByTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		respBody     []entities.Drink
		hasRespBody  bool
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodGet,
			path:         "/drink/tag/:tag",
			respBody:     []entities.Drink{{Name: "test01"}},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "tag",
					value: "spicy",
				},
			},
		},
		{
			name:         "test02",
			method:       http.MethodGet,
			path:         "/drink/tag/:tag",
			respBody:     []entities.Drink{{Name: "test02", Tags: []string{"spicy", "non-alcohol"}}},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "tag",
					value: "non-alcohol",
				},
			},
		},
	}

	mockStorage.EXPECT().DrinksByTags(gomock.Any(), []string{"spicy"}).Return(ts[0].respBody, nil)
	mockStorage.EXPECT().DrinksByTags(gomock.Any(), []string{"non-alcohol"}).Return(ts[1].respBody, nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		req := httptest.NewRequest(v.method, v.path, nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)
		for _, v := range v.params {
			c.SetParamNames(v.name)
			c.SetParamValues(v.value)
		}

		if assert.NoError(t, h.drinksByTags(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
			if v.hasRespBody {
				resData, _ := json.Marshal(&v.respBody)
				assert.JSONEq(t, string(resData), rec.Body.String())
			}
		}
	}
}

func Test_httpHandler_allDrinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		respBody     []entities.Drink
		hasRespBody  bool
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodGet,
			path:         "/drink/id/:id",
			respBody:     []entities.Drink{{Name: "test01"}},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "id",
					value: "1",
				},
			},
		},
		{
			name:         "test02",
			method:       http.MethodGet,
			path:         "/drink/id/:id",
			respBody:     []entities.Drink{{Name: "test02", Tags: []string{"spicy", "non-alcohol"}}},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "id",
					value: "2",
				},
			},
		},
	}

	mockStorage.EXPECT().AllDrinks(gomock.Any(), 1).Return(ts[0].respBody, nil)
	mockStorage.EXPECT().AllDrinks(gomock.Any(), 2).Return(ts[1].respBody, nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		req := httptest.NewRequest(v.method, v.path, nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)
		for _, v := range v.params {
			c.SetParamNames(v.name)
			c.SetParamValues(v.value)
		}

		if assert.NoError(t, h.allDrinks(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
			if v.hasRespBody {
				resData, _ := json.Marshal(&v.respBody)
				assert.JSONEq(t, string(resData), rec.Body.String())
			}
		}
	}
}

func Test_httpHandler_drinkByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockDrinkModel(ctrl)
	type params struct {
		name  string
		value string
	}
	type TestCase struct {
		name         string
		method       string
		path         string
		respBody     entities.Drink
		hasRespBody  bool
		expectedCode int
		params       []params
	}
	ts := []TestCase{
		{
			name:         "test01",
			method:       http.MethodGet,
			path:         "/drinks/:name",
			respBody:     entities.Drink{},
			hasRespBody:  false,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "name",
					value: "test01",
				},
			},
		},
		{
			name:   "test02",
			method: http.MethodGet,
			path:   "/drinks/:name",
			respBody: entities.Drink{
				Name: "test02",
				Tags: []string{"spicy", "non-alcohol"},
			},
			hasRespBody:  true,
			expectedCode: http.StatusOK,
			params: []params{
				{
					name:  "name",
					value: "test02",
				},
			},
		},
	}

	mockStorage.EXPECT().DrinkByName(gomock.Any(), "test01").Return(ts[0].respBody, nil)
	mockStorage.EXPECT().DrinkByName(gomock.Any(), "test02").Return(ts[1].respBody, nil)

	h := &httpHandler{mockStorage, nil, nil, nil}

	for _, v := range ts {

		e := echo.New()
		req := httptest.NewRequest(v.method, v.path, nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(v.path)
		for _, v := range v.params {
			c.SetParamNames(v.name)
			c.SetParamValues(v.value)
		}

		if assert.NoError(t, h.drinkByName(c)) {
			assert.Equal(t, v.expectedCode, rec.Code)
			if v.hasRespBody {

				resData, _ := json.Marshal(&v.respBody)
				assert.JSONEq(t, string(resData), rec.Body.String())
			}
		}
	}
}
