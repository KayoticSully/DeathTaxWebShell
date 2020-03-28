package deathtax

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// DeathTax powershell process object
type DeathTax struct {
	Input  chan string
	Output chan string

	process *exec.Cmd
}

// // Input returns a channel to send strings to the process stdin
// func (dt *DeathTax) Input() chan<- string {
// 	return dt.input
// }

// // Output returns a channel to get strings from the process stdout
// func (dt *DeathTax) Output() <-chan string {
// 	return dt.input
// }

// New returns a new running DeathTax process
func New() *DeathTax {
	dt := &DeathTax{
		Input:   make(chan string, 1),
		Output:  make(chan string, 1),
		process: exec.Command("./test.sh"),
	}

	log.Println("Created Instance")

	stdin, err := dt.process.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go inputPump(stdin, dt.Input)

	stdout, err := dt.process.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	go outputPump(stdout, dt.Output)

	log.Println("Starting SubProcess")
	if err := dt.process.Start(); err != nil {
		log.Fatal(err)
	}

	return dt
}

func inputPump(stdin io.WriteCloser, input <-chan string) {
	log.Println("Starting Input Pump")
	defer stdin.Close()

	for msg := range input {
		io.WriteString(stdin, msg)
	}
}

func outputPump(stdout io.ReadCloser, output chan<- string) {
	log.Println("Starting Output Pump")
	defer stdout.Close()

	stdScanner := bufio.NewScanner(stdout)

	for stdScanner.Scan() {
		text := stdScanner.Text()
		log.Printf("Sending String '%s' To Output Channel\n", text)
		output <- text
		log.Println("SENT!")
	}
}
