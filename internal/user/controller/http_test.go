package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	mockAuth "github.com/SapolovichSV/backprogeng/mocks/authmiddleware"
	mocks "github.com/SapolovichSV/backprogeng/mocks/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Пример тестов для httpHandler.Login
func Test_httpHandler_Login(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name         string
		mockSetup    func(*mocks.MockuserModel, *mockAuth.MockauthService)
		wantHTTPCode int
		wantErr      bool
	}{
		{
			name: "login_ok",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Login(gomock.Any()).
					Return(entities.User{ID: 100, Username: "TestName"}, nil)
				mu.EXPECT().UserByID(gomock.Any(), 100).
					Return(entities.User{ID: 100, Username: "TestName"}, nil)
			},
			wantHTTPCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "login_unauthorized",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Login(gomock.Any()).
					Return(entities.User{}, errors.New("unauthorized"))
			},
			wantHTTPCode: http.StatusUnauthorized,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockuserModel(ctrl)
			mockAuth := mockAuth.NewMockauthService(ctrl)

			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage, mockAuth)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/login", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &httpHandler{
				st:   mockStorage,
				echo: e,
				ctx:  context.Background(),
				auth: mockAuth,
			}
			_ = h.Login(c)
			require.Equal(t, tt.wantHTTPCode, rec.Code)
		})
	}
}

// Пример тестов для httpHandler.UserByID
func Test_httpHandler_UserByID(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name         string
		mockSetup    func(*mocks.MockuserModel, *mockAuth.MockauthService)
		wantHTTPCode int
		wantErr      bool
	}{
		{
			name: "ok",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Auth(gomock.Any()).
					Return(entities.User{ID: 50}, nil)
				mu.EXPECT().UserByID(gomock.Any(), 50).
					Return(entities.User{ID: 50, Username: "John"}, nil)
			},
			wantHTTPCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "no_token",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Auth(gomock.Any()).
					Return(entities.User{}, errors.New("no token"))
			},
			wantHTTPCode: http.StatusUnauthorized,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockuserModel(ctrl)
			mockAuth := mockAuth.NewMockauthService(ctrl)

			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage, mockAuth)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/user/50", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &httpHandler{
				st:   mockStorage,
				echo: e,
				ctx:  context.Background(),
				auth: mockAuth,
			}
			_ = h.UserByID(c)

			require.Equal(t, tt.wantHTTPCode, rec.Code)
		})
	}
}

// Пример тестов для httpHandler.AddFav
func Test_httpHandler_AddFav(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name         string
		mockSetup    func(*mocks.MockuserModel, *mockAuth.MockauthService)
		wantHTTPCode int
		wantErr      bool
	}{
		{
			name: "add_fav_ok",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Auth(gomock.Any()).
					Return(entities.User{ID: 10}, nil)
				mu.EXPECT().AddFav(gomock.Any(), "Coke", entities.User{ID: 10}.ID)
			},
			wantHTTPCode: http.StatusAccepted,
			wantErr:      false,
		},
		{
			name: "auth_error",
			mockSetup: func(mu *mocks.MockuserModel, ma *mockAuth.MockauthService) {
				ma.EXPECT().Auth(gomock.Any()).
					Return(entities.User{}, errors.New("unauthorized"))
			},
			wantHTTPCode: http.StatusUnauthorized,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockuserModel(ctrl)
			mockAuth := mockAuth.NewMockauthService(ctrl)

			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage, mockAuth)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/user/fav", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("drinkname")
			c.SetParamValues("Coke")
			h := &httpHandler{
				st:   mockStorage,
				echo: e,
				ctx:  context.Background(),
				auth: mockAuth,
			}
			_ = h.AddFav(c)

			require.Equal(t, tt.wantHTTPCode, rec.Code)
		})
	}
}
