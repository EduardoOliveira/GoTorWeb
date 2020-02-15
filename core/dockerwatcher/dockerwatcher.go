package dockerwatcher

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/EduardoOliveira/GoTorWeb/core/lib"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Watcher struct {
	client   *client.Client
	filters  filters.Args
	watchers map[string][]func(*lib.Container)
}

func New(filters filters.Args) (w *Watcher, err error) {

	w = new(Watcher)
	w.watchers = make(map[string][]func(*lib.Container), 0)

	w.client, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	w.filters = filters

	return
}

func (w *Watcher) GetRunning() (containers []*lib.Container, err error) {
	containers = make([]*lib.Container, 0)
	cl, err := w.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return
	}
	for _, dc := range cl {
		log.Println(fmt.Sprintf("%+v    %s", dc, strings.TrimPrefix(dc.Names[0], "/")))
		c := &lib.Container{
			ID:          dc.ID,
			Name:        strings.TrimPrefix(dc.Names[0], "/"),
			Port:        dc.Labels["GWT-PORT"],
			PortForward: dc.Labels["GWT-PORT-FW"],
		}
		full, err := w.client.ContainerInspect(context.Background(), dc.ID)
		if err != nil {
			log.Println(c.Name, " fail to inspect container", err)
			continue
		}
		c.IPAddr = full.NetworkSettings.DefaultNetworkSettings.IPAddress
		containers = append(containers, c)
	}
	return
}

func (w *Watcher) AddWatcher(e string, f func(*lib.Container)) {
	_, ok := w.watchers[e]
	if !ok {
		w.watchers[e] = make([]func(*lib.Container), 0)
	}
	w.watchers[e] = append(w.watchers[e], f)
}

func (w *Watcher) Start() {
	go func() {
		event, errChan := w.client.Events(context.Background(), types.EventsOptions{Filters: w.filters})
		for {
			select {
			case e := <-event:
				c := &lib.Container{
					ID:          e.Actor.ID,
					Name:        e.Actor.Attributes["name"],
					Port:        e.Actor.Attributes["GWT-PORT"],
					PortForward: e.Actor.Attributes["GWT-PORT-FW"],
				}
				if c.Port == "" || c.PortForward == "" {
					log.Printf(e.Actor.Attributes["name"], " no GWT-PORT or GWT-PORT-FW discarding event")
					continue
				}

				if e.Action == "start" {
					full, err := w.client.ContainerInspect(context.Background(), e.Actor.ID)
					if err != nil {
						log.Println(e.Actor.Attributes["name"], " fail to inspect container", err)
						continue
					}
					c.IPAddr = full.NetworkSettings.DefaultNetworkSettings.IPAddress
				}

				for _, f := range w.watchers[e.Action] {
					go f(c)
				}
			case e := <-errChan:
				log.Println("error: ", e)
			}
		}
	}()
}
