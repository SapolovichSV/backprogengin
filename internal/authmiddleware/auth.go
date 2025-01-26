package authmiddleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/SapolovichSV/backprogeng/internal/errlib"
	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// copypast of the same code from internal/user/controller/http.go
// because need to mock
// TODO: refactor to use only one copy of this code
type authService interface {
	Auth(c echo.Context) (entities.User, error)
	Login(c echo.Context) (entities.User, error)
	Register(c echo.Context, user entities.User) error
}
type secretKey struct {
	key string
}

func parseSecretKey() secretKey {
	key := os.Getenv("SECRET")
	if key == "" {
		key = "simple"
	}
	return secretKey{
		key: key,
	}
}

type jwtCustomClaims struct {
	id       int
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}
type authMiddle struct {
	secretKey secretKey
}

func New() *authMiddle {
	return &authMiddle{
		secretKey: parseSecretKey(),
	}
}
func (a *authMiddle) Register(c echo.Context, user entities.User) error {

	if err := c.Bind(&user); err != nil {
		return errlib.WrapErr(err, "failing to get user data from http request")
	}
	claims := jwtCustomClaims{
		user.ID,
		user.Username,
		user.Password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.secretKey.key))
	if err != nil {
		return errlib.WrapErr(err, "failing to sign token")
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: t,
	}
	c.SetCookie(&cookie)
	return nil
}
func (a *authMiddle) Auth(c echo.Context) (entities.User, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return entities.User{}, errlib.WrapErr(err, "failing to get cookie with user Info")
	}
	claims, err := getClaims(cookie, a.secretKey.key)
	if err != nil {
		return entities.User{}, errlib.WrapErr(err, "failing to get claims from token")
	}
	user := entities.User{
		ID:       claims.id,
		Username: claims.Username,
		Password: claims.Password,
	}
	return user, nil
}
func (a *authMiddle) Login(c echo.Context) (entities.User, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return entities.User{}, errlib.WrapErr(err, "failing to get cookie with user Info")
	}
	claims, err := getClaims(cookie, a.secretKey.key)
	if err != nil {
		return entities.User{}, errlib.WrapErr(err, "failing to get claims from token")
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.secretKey.key))
	if err != nil {
		return entities.User{}, errlib.WrapErr(err, "failing to sign token")
	}
	cookie = &http.Cookie{
		Name:  "token",
		Value: t,
	}
	c.SetCookie(cookie)
	user := entities.User{
		ID:       claims.id,
		Username: claims.Username,
		Password: claims.Password,
	}
	return user, nil
}

func getClaims(cookie *http.Cookie, key string) (jwtCustomClaims, error) {
	if cookie == nil {
		return jwtCustomClaims{}, errlib.WrapErr(errors.New("no token cookie"), "no token cookie")
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return jwtCustomClaims{}, errlib.WrapErr(err, "failing to parse token")
	}
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return jwtCustomClaims{}, errlib.WrapErr(errors.New("failing to get claims from token"), "token broken")
	}
	return *claims, nil
}
