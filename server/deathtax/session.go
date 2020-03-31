package deathtax

import (
	"bufio"
	"io"
	"log"
	"os/exec"

	"github.com/fasthttp/websocket"
)

// Session powershell process object
type Session struct {
	process *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
}

// NewSession returns a new running DeathTax process
func NewSession() *Session {
	process := exec.Command("pwsh", "/usr/local/share/deathtax/DeathTax")

	// Input and Output pipes need to be created before the go
	// routines start. Otherwise data will be missed between process
	// start and go routine start.  There is an unknowable delay between
	// telling a go routine to start and when it actually starts.
	stdin, err := process.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := process.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting SubProcess")

	if err := process.Start(); err != nil {
		log.Fatal(err)
	}

	return &Session{
		process: process,
		stdin:   stdin,
		stdout:  stdout,
	}
}

// RunWebsocketProxy pipes input and output channels from/to the websocket
func (s *Session) RunWebsocketProxy(wsConn *websocket.Conn) {
	go s.inputPump(wsConn)
	s.outputPump(wsConn)
}

func (s *Session) inputPump(wsConn *websocket.Conn) {
	// stdin is created outside of the go routine but it should
	// be setup to close when the go routine exits
	defer s.stdin.Close()

	var err error
	var msg []byte

	for {
		_, msg, err = wsConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		_, err = io.WriteString(s.stdin, string(msg))
	}
}

func (s *Session) outputPump(wsConn *websocket.Conn) {
	// stdout is created outside of the go routine but it should
	// be setup to close when the go routine exits
	defer s.stdout.Close()

	var text []byte
	var err error
	stdScanner := bufio.NewScanner(s.stdout)

	for stdScanner.Scan() {
		text = []byte(stdScanner.Text())
		if err = wsConn.WriteMessage(websocket.TextMessage, text); err != nil {
			log.Println("write:", err)
			return
		}
	}
}
