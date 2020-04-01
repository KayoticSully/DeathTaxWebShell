package deathtax

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/fasthttp/websocket"
)

// Session powershell process object
type Session struct {
	process *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser

	stdoutReadLock sync.Mutex
	stdoutScanner  *bufio.Scanner
	firstLine      string
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

	scanner := bufio.NewScanner(stdout)
	scanner.Split(scanLinesWithInput)

	s := &Session{
		process: process,
		stdin:   stdin,
		stdout:  stdout,

		stdoutReadLock: sync.Mutex{},
		stdoutScanner:  scanner,
		firstLine:      "",
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
	s.stdoutReadLock.Lock()
	defer s.stdoutReadLock.Unlock()

	for s.stdoutScanner.Scan() {
		s.firstLine = s.stdoutScanner.Text()

		if s.firstLine != "" {
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

	s.stdoutReadLock.Lock()
	defer s.stdoutReadLock.Unlock()

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(s.firstLine))
	if err != nil {
		log.Println("write:", err)
		return
	}

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

func scanLinesWithInput(data []byte, atEOF bool) (advance int, token []byte, err error) {
	trimmedData := bytes.TrimSpace(data)
	if i := bytes.LastIndexByte(data, ':'); i == len(trimmedData) {
		// We have a request for input
		return len(data), data, nil
	}

	return bufio.ScanLines(data, atEOF)
}
