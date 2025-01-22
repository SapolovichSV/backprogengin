package authmiddleware

import (
	"net/http"
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
		want    echo.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			got, err := a.Register(tt.args.c, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authMiddle.Register() = %v, want %v", got, tt.want)
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
		want    echo.Context
		want1   entities.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			got, got1, err := a.Auth(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authMiddle.Auth() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("authMiddle.Auth() got1 = %v, want %v", got1, tt.want1)
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
		want    echo.Context
		want1   entities.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authMiddle{
				secretKey: tt.fields.secretKey,
			}
			got, got1, err := a.Login(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("authMiddle.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authMiddle.Login() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("authMiddle.Login() got1 = %v, want %v", got1, tt.want1)
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
		// TODO: Add test cases.
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
