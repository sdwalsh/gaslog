package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sdwalsh/threadinator/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

// User Landing Page
func (env *Env) Dashboard(c echo.Context) error {
	cookie, err := c.Cookie("email")
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("%#v\n", err))
	}
	user, err := env.Db.GetUserByEmail(cookie.Value)
	result, err := env.Db.GetCharsByUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("%#v\n", err))
	}
	csrf, _ := c.Get("csrf").(string)
	data := struct {
		User   *models.User
		Result *[]models.Character
		Csrf   string
		Auth   bool
	}{
		user,
		result,
		csrf,
		true,
	}
	return c.Render(http.StatusOK, "dashboard.html", data)
}

func Signup(c echo.Context) error {
	fmt.Printf("%#v\n", c.Get("csrf"))
	return c.Render(http.StatusOK, "signup.html", c.Get("csrf"))
}

func (env *Env) CreateUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	//repeat := c.FormValue("repeatPassword")
	digest, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return c.Render(http.StatusTeapot, "signup.html", "digest wrong")
	}

	env.Db.InsertUser(name, digest, email)
	return c.Redirect(http.StatusSeeOther, "/")
}

// Using cookies for persistant login information for now
// will update to sessions before release
func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", c.Get("csrf"))
}

func (env *Env) CreateLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	result, err := env.Db.GetUserByEmail(email)
	if err != nil {
		return c.Render(http.StatusTeapot, "index.html", "error")
	}
	err = bcrypt.CompareHashAndPassword(result.Digest, []byte(password))
	if err != nil {
		return c.Render(http.StatusTeapot, "index.html", "error")
	}

	cookie := new(http.Cookie)
	cookie.Name = "email"
	cookie.Value = email
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, result.Id+"/dashboard")
}

// Function to remove cookie and logout user from website
// setting cookie's max age to -1 will delete cookie
func (env *Env) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "email"
	cookie.Value = ""
	cookie.MaxAge = -1

	return c.Redirect(http.StatusSeeOther, "/")
}
