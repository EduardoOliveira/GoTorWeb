package core

import (
	"log"

	"github.com/EduardoOliveira/GoTorWeb/core/dockerwatcher"
	"github.com/EduardoOliveira/GoTorWeb/core/tormanager"
	"github.com/EduardoOliveira/GoTorWeb/core/webui"
	"github.com/docker/docker/api/types/filters"
)

func Run() {
	filters := filters.NewArgs()
	filters.Add("label", "GTW=1")
	w, err := dockerwatcher.New(filters)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("running")

	containers, err := w.GetRunning()

	if err != nil {
		log.Println("unhable to check running containers", err)
	}

	tm := tormanager.NewWithLocalPort(containers, 80)
	go webui.New(80, tm).Start()

	w.AddWatcher("start", tm.HandleCreation)
	w.AddWatcher("destroy", tm.HandleDeletion)

	w.Start()

	<-make(chan bool, 0)
}
