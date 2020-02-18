package webui

import (
	"html/template"
	"log"
	"os"

	"github.com/EduardoOliveira/GoTorWeb/core/tormanager"
	"github.com/labstack/echo"
)

type WebUI struct {
	port     int
	username string
	password string
	template *template.Template
	e        *echo.Echo
	tm       *tormanager.TorManager
}

func New(port int, tm *tormanager.TorManager) *WebUI {

	wu := &WebUI{
		port: port,
		tm:   tm,
	}

	wu.username = os.Getenv("GTW-USERNAME")
	wu.password = os.Getenv("GTW-PASSWORD")

	if wu.username == "" || wu.password == "" {
		log.Fatal("web ui: username or password not defined")
	}

	t, err := template.ParseFiles("./templates/webui.tmpl")
	if err != nil {
		log.Panic("failed to parse web ui templte")
	}

	wu.template = t

	wu.initServer()
	return wu
}

func (wu *WebUI) Start() {
	log.Println(wu.startServer())
}
