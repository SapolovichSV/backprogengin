package authmiddleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/labstack/echo/v4"
)

func Test_authMiddle_Register(t *testing.T) {
	type fields struct {
		secretKey secretKey
	}
	type args struct {
		c    echo.Context
		user entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Valid registration",
			fields: fields{secretKey: secretKey{key: "testkey"}},
			args: func() args {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return args{c: c, user: entities.User{ID: 1, Username: "TestUser", Password: "123"}}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			if err := a.Register(tt.args.c, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authMiddle_Auth(t *testing.T) {
	type fields struct {
		secretKey secretKey
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.User
		wantErr bool
	}{
		{
			name:   "No token cookie",
			fields: fields{secretKey: secretKey{key: "testkey"}},
			args: func() args {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				return args{c: e.NewContext(req, rec)}
			}(),
			want:    entities.User{},
			wantErr: true,
		},
		{
			name:   "Invalid token cookie",
			fields: fields{secretKey: secretKey{key: "testkey"}},
			args: func() args {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				cookie := &http.Cookie{Name: "token", Value: "broken.token"}
				req.AddCookie(cookie)
				rec := httptest.NewRecorder()
				return args{c: e.NewContext(req, rec)}
			}(),
			want:    entities.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			got, err := a.Auth(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authMiddle.Auth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authMiddle_Login(t *testing.T) {
	type fields struct {
		secretKey secretKey
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.User
		wantErr bool
	}{
		{
			name:   "No cookie",
			fields: fields{secretKey: secretKey{key: "testkey"}},
			args: func() args {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				return args{c: e.NewContext(req, rec)}
			}(),
			want:    entities.User{},
			wantErr: true,
		},
		{
			name:   "Invalid token cookie",
			fields: fields{secretKey: secretKey{key: "testkey"}},
			args: func() args {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				cookie := &http.Cookie{Name: "token", Value: "broken.token"}
				req.AddCookie(cookie)
				rec := httptest.NewRecorder()
				return args{c: e.NewContext(req, rec)}
			}(),
			want:    entities.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			got, err := a.Login(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authMiddle.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClaims(t *testing.T) {
	type args struct {
		cookie *http.Cookie
		key    string
	}
	tests := []struct {
		name    string
		args    args
		want    jwtCustomClaims
		wantErr bool
	}{
		{
			name: "No cookie",
			args: args{
				cookie: nil,
				key:    "testkey",
			},
			wantErr: true,
		},
		{
			name: "Empty cookie",
			args: args{
				cookie: &http.Cookie{Name: "token", Value: ""},
				key:    "testkey",
			},
			wantErr: true,
		},
		{
			name: "Broken token",
			args: args{
				cookie: &http.Cookie{Name: "token", Value: "broken.token"},
				key:    "testkey",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getClaims(tt.args.cookie, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}
