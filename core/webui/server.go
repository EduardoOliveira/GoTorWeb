package webui

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (wu *WebUI) initServer() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${remote_ip} - - ${time_rfc3339} ${latency_human} "${method} ${uri}" ${status} ${bytes_out} "${refer}" "${user_agent}"` + "\n",
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
	return c.JSON(http.StatusOK, wu.tm.Containers)
}
