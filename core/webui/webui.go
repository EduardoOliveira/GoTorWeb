package webui

import (
	"log"

	"github.com/EduardoOliveira/GoTorWeb/core/tormanager"
	"github.com/labstack/echo"
)

type WebUI struct {
	port int
	e    *echo.Echo
	tm   *tormanager.TorManager
}

func New(port int, tm *tormanager.TorManager) *WebUI {

	wu := &WebUI{
		port: port,
		tm:   tm,
	}
	wu.initServer()
	return wu
}

func (wu *WebUI) Start() {
	log.Println(wu.startServer())
}
