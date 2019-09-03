package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

const (
	JWTURL          = "http://%s/consumers/%s/jwt"
	DefaultKong     = "127.0.0.1:8001"
	DefaultConsumer = "my-user"
	DefaultDbURI    = "postgres://server:@127.0.0.1:5432/server?sslmode=disable"
)

type Response struct {
	Data struct {
		AccountID   string      `json:"account_id"`
		Credentials Credentials `json:"credentials"`
	} `json:"data"`
}

func NewResponse(u *User, token string) Response {
	response := Response{}
	response.Data.Credentials = Credentials{JWT: token}
	response.Data.AccountID = u.AccountID
	return response
}

func getLoginData(c echo.Context) (Login, error) {
	type input struct {
		Data Login `json:"data"`
	}
	i := input{}
	if err := c.Bind(&i); err != nil {
		return Login{}, err
	}
	if i.Data.Email == "" || i.Data.Password == "" {
		return Login{}, echo.NewHTTPError(http.StatusBadRequest, "email and password should be provided")
	}
	return i.Data, nil
}

func getDB() *sql.DB {
	dbURI := os.Getenv("PG_URI")
	if dbURI == "" {
		dbURI = DefaultDbURI
	}
	uri, err := pq.ParseURL(dbURI)
	if err != nil {
		panic(err)
	}
	pdb, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	return pdb
}

func main() {
	e := echo.New()

	repo := NewRepository(getDB())
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/login", func(c echo.Context) error {
		login, err := getLoginData(c)
		if err != nil {
			return err
		}
		user, err := repo.Get(login.Email)
		if err != nil {
			if _, ok := err.(ErrNotFound); ok {
				return echo.NewHTTPError(http.StatusNotFound, "user with such email not found")
			}
			return err
		}
		if !user.CheckPassword(login.Password) {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong email or password")
		}
		token, err := CraftToken(user)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, NewResponse(user, token))
	})
	e.POST("/auth", func(c echo.Context) error {
		login, err := getLoginData(c)
		if err != nil {
			return err
		}
		ks, err := GetSecret()
		if err != nil {
			return err
		}
		user := &User{
			Email:    login.Email,
			Password: login.Password,
			Key:      ks.Key,
			Secret:   ks.Secret,
		}
		user, err = repo.Add(user)
		if err != nil {
			return err
		}
		token, err := CraftToken(user)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, NewResponse(user, token))
	})
	e.Logger.Fatal(e.Start(":80"))
}

func GetSecret() (*KeySecret, error) {
	kong := DefaultKong
	if v, ok := os.LookupEnv("KONG"); ok {
		kong = v
	}
	consumer := DefaultConsumer
	if v, ok := os.LookupEnv("CONSUMER"); ok {
		consumer = v
	}
	jwtUrl := fmt.Sprintf(JWTURL, kong, consumer)
	r, err := http.PostForm(jwtUrl, url.Values{})
	if err != nil {
		return nil, err
	}
	ks := &KeySecret{}
	return ks, json.NewDecoder(r.Body).Decode(ks)
}

func CraftToken(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["iss"] = user.Key
	t, err := token.SignedString([]byte(user.Secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
