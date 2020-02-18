package webui

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (wu *WebUI) initServer() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${remote_ip} - - ${time_rfc3339} ${latency_human} "${method} ${uri}" ${status} ${bytes_out} "${refer}" "${user_agent}"` + "\n",
	}))
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(wu.username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(wu.password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Gzip())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.GET("/", wu.serveUI)
	wu.e = e
}

func (wu *WebUI) startServer() error {
	return wu.e.Start(fmt.Sprintf(":%d", wu.port))
}

func (wu *WebUI) serveUI(c echo.Context) error {

	var tpl bytes.Buffer
	if err := wu.template.Execute(&tpl, wu.tm.Containers); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.HTML(http.StatusOK, tpl.String())
}
