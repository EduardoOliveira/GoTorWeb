package tormanager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/EduardoOliveira/GoTorWeb/core/lib"
)

type TorManager struct {
	Containers  map[string]*lib.Container
	LocalPort   int
	LocalAddess string
	torCmd      *exec.Cmd
}

func New(runningContainers []*lib.Container) (tm *TorManager) {
	tm = new(TorManager)
	tm.Containers = make(map[string]*lib.Container, 0)
	for _, c := range runningContainers {
		tm.Containers[c.ID] = c
	}
	go tm.watchRuningAddresses()
	return tm
}

func NewWithLocalPort(runningContainers []*lib.Container, port int) (tm *TorManager) {
	tm = New(runningContainers)
	tm.LocalPort = port

	err := tm.genRC()
	if err != nil {
		log.Println(err)
		return
	}
	go tm.startTor()
	return tm
}

func (tm *TorManager) AddLocalPort(port int) {
	tm.LocalPort = port
}

func (tm *TorManager) HandleCreation(c *lib.Container) {
	tm.Containers[c.ID] = c
	err := tm.genRC()
	if err != nil {
		log.Println(err)
		return
	}
	go tm.reloadTor()
}

func (tm *TorManager) HandleDeletion(c *lib.Container) {
	delete(tm.Containers, c.ID)
	err := tm.genRC()
	if err != nil {
		log.Println(err)
		return
	}
	go tm.reloadTor()
}

func (tm *TorManager) reloadTor() {
	if tm.torCmd != nil {
		if err := tm.torCmd.Process.Kill(); err != nil {
			log.Printf("could not kill cmd: %v", err)
		}
	}
	tm.startTor()
}

func (tm *TorManager) startTor() {
	tm.torCmd = exec.Command("/usr/bin/tor", "-f", "/opt/go-torrc")
	tm.torCmd.Stdout = os.Stdout
	tm.torCmd.Stderr = os.Stderr
	if err := tm.torCmd.Run(); err != nil {
		log.Printf("could not run cmd: %v", err)
	}
}

func (tm *TorManager) watchRuningAddresses() {
	for {
		<-time.NewTicker(2 * time.Second).C
		for _, c := range tm.Containers {
			if c.Address == "" {
				file, err := os.Open(fmt.Sprintf("/config/%s/hostname", c.Name))
				if err != nil {
					continue
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					c.Address = scanner.Text()
				}
				fmt.Println("Hosting:", c.Address, c.PortForward)
				continue
			}
		}
		//ahhhhhh
		if tm.LocalPort != 0 && tm.LocalAddess == "" {
			file, err := os.Open("/config/local/hostname")
			if err != nil {
				continue
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				tm.LocalAddess = scanner.Text()
			}
			fmt.Println("------------------------------------")
			fmt.Println("Local:", tm.LocalAddess, tm.LocalPort)
			fmt.Println("------------------------------------")
		}
	}
}
