package deathtax

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
)

// Session powershell process object
type Session struct {
	process *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser

	startupMux    sync.Mutex
	stdoutScanner *bufio.Scanner
	firstLine     string
	proxyStarted  bool
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

	s := &Session{
		process: process,
		stdin:   stdin,
		stdout:  stdout,

		startupMux:    sync.Mutex{},
		stdoutScanner: bufio.NewScanner(stdout),
		firstLine:     "",
		proxyStarted:  false,
	}

	go s.primeProcess()
	return s
}

// RunWebsocketProxy pipes input and output channels from/to the websocket
func (s *Session) RunWebsocketProxy(wsConn *websocket.Conn) {
	go s.inputPump(wsConn)
	s.outputPump(wsConn)
}

// IsReady returns true if the session is fully booted and available for use.
func (s *Session) IsReady() bool {
	return s.firstLine != ""
}

func (s *Session) primeProcess() {
	s.startupMux.Lock()
	for s.firstLine = ""; s.firstLine == ""; s.firstLine = s.stdoutScanner.Text() {
		s.startupMux.Unlock()

		// Take a small break
		time.Sleep(time.Millisecond * 250)

		// If the proxy started while sleeping, stop checking for output
		s.startupMux.Lock()
		if s.proxyStarted {
			s.startupMux.Unlock()
			return
		}
	}
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

	s.startupMux.Lock()
	s.proxyStarted = true

	if s.firstLine != "" {
		err = wsConn.WriteMessage(websocket.TextMessage, []byte(s.firstLine))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
	s.startupMux.Unlock()

	for s.stdoutScanner.Scan() {
		text = []byte(s.stdoutScanner.Text())

		// Whitespace or blank output means a blank line
		// was emitted from the script. Send a newline.
		if strings.TrimSpace(string(text)) == "" {
			text = []byte("\n")
		}

		err = wsConn.WriteMessage(websocket.TextMessage, text)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
