package modules

import (
	"fmt"
	"go-scrap/config"
	"log"
	"net"
	"os/exec"
	"time"
)

type Tor struct {
	ControlAddress string
}

func NewTor(ctrlServer string) *Tor {
	return &Tor{
		ControlAddress: ctrlServer,
	}
}

func (t *Tor) Init() {
	cmd := exec.Command("tor")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Failed to start tor : ", err)
	}
	fmt.Println("Wait for tor to start")
	time.Sleep(10 * time.Second)
	fmt.Println("TOR started")
}

func (t *Tor) ChangeIP() {
	conn, _ := net.Dial("tcp", t.ControlAddress)
	fmt.Fprintf(conn, "AUTHENTICATE \"%s\"\r\n", config.Cfg.TORCONTROL_PASSWORD)
	fmt.Fprintf(conn, "SIGNAL NEWNYM\r\n")
	time.Sleep(1 * time.Second)
	fmt.Println("IP Changed")
}
