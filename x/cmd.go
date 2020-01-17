package x

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"text/template"

	"github.com/MasteryConnect/pipe/message"
)

// Cmd will execute a system command
type Cmd struct {
	Name string
	Args []string

	NoStdin bool // don't send the msg.String() to stdin of the command
}

// Command wraps the os/exec Command func
func Command(name string, arg ...string) *Cmd {
	return &Cmd{Name: name, Args: arg}
	//return &Cmd{Cmd: exec.Command(name, arg...)}
}

// P is the producer function for the pipe/line
func (e Cmd) P(out chan<- interface{}, errs chan<- error) {
	var errwg, readwg sync.WaitGroup

	// stdout
	c := exec.Command(e.Name, e.Args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		errs <- err
		return
	}
	reader := bufio.NewReader(stdout)

	// stderr
	stderr, err := c.StderrPipe()
	if err != nil {
		errs <- err
		return
	}
	errScanner := bufio.NewScanner(stderr)

	// start the command
	if err := c.Start(); err != nil {
		errs <- err
		return
	}

	// read all the lines and send down stream
	errwg.Add(1)
	go func() {
		defer errwg.Done()
		for errScanner.Scan() {
			errs <- fmt.Errorf(errScanner.Text()) // Println will add back the final '\n'
		}
	}()

	// read all the lines and send down stream
	readwg.Add(1)
	go func() {
		defer readwg.Done()
		readAndSend(reader, out, errs)
	}()

	// wait for read to finish before calling c.Wait()
	readwg.Wait()

	// wait for close
	if err := c.Wait(); err != nil {
		errs <- err
	}

	// wait for the errors to all be processed
	errwg.Wait()
}

// TStream is the transform function for the pipe/line
// It will run the shell command once, and keep piping the messages
// in to stdin of the command.
func (e Cmd) TStream(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	var errwg, readwg sync.WaitGroup
	c := exec.Command(e.Name, e.Args...)
	// stdin
	stdin, err := c.StdinPipe()
	if err != nil {
		errs <- err
		return
	}

	// stdout
	stdout, err := c.StdoutPipe()
	if err != nil {
		errs <- err
		return
	}
	reader := bufio.NewReader(stdout)

	// stderr
	stderr, err := c.StderrPipe()
	if err != nil {
		errs <- err
		return
	}
	errScanner := bufio.NewScanner(stderr)

	// start the command
	if err := c.Start(); err != nil {
		errs <- err
		return
	}

	// read all the lines and send down stream
	errwg.Add(1)
	go func() {
		defer errwg.Done()
		for errScanner.Scan() {
			errs <- fmt.Errorf(errScanner.Text()) // Println will add back the final '\n'
		}
	}()

	// read all the lines and send down stream
	readwg.Add(1)
	go func() {
		defer readwg.Done()
		readAndSend(reader, out, errs)
	}()

	// read all the messages and write to stdin
	for msg := range in {
		fmt.Fprintf(stdin, "%s\n", msg.(fmt.Stringer).String())
	}

	// let the stdin drain before closing
	// There could be more do to here after some research instead of sleep.
	// Since this only happens once in a streaming case, it shouldn't be that bad.
	time.Sleep(100 * time.Millisecond)

	stdin.Close()

	// wait for the reads before calling c.Wait()
	readwg.Wait()

	// wait for close
	if err := c.Wait(); err != nil {
		errs <- err
	}

	// wait for the errors to finish
	errwg.Wait()
}

// T is the transform function for the pipe/line.
// It will run the shell command for each message.
func (e Cmd) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	var errwg, readwg sync.WaitGroup
	var cmdTmpl *template.Template
	argTmpls := []*template.Template{}
	cmdTmpl = template.Must(template.New("x.cmd.name").Parse(e.Name))
	for i, arg := range e.Args {
		argTmpls = append(argTmpls, template.Must(
			template.New(fmt.Sprintf("x.cmd.args.%d", i)).Parse(arg),
		))
	}

	for msg := range in {
		name := e.Name
		args := append([]string{}, e.Args...)

		if cmd, ok := msg.(Commander); ok { // use the message as the command if it is a Commander
			name = cmd.Command()
			args = cmd.Arguments()
		} else { // otherwise, use the static name and arg templates
			var b strings.Builder

			// Name
			err := cmdTmpl.Execute(&b, msg)
			if err != nil {
				errs <- err
			} else {
				name = b.String()
			}

			// Args
			for i, argTmpl := range argTmpls {
				b.Reset()
				err := argTmpl.Execute(&b, msg)
				if err != nil {
					errs <- err
				} else {
					args[i] = b.String()
				}
			}
		}

		c := exec.Command(name, args...)
		// stdin
		stdin, err := c.StdinPipe()
		if err != nil {
			errs <- err
			return
		}

		// stdout
		stdout, err := c.StdoutPipe()
		if err != nil {
			errs <- err
			return
		}
		reader := bufio.NewReader(stdout)

		// stderr
		stderr, err := c.StderrPipe()
		if err != nil {
			errs <- err
			return
		}
		errScanner := bufio.NewScanner(stderr)

		// start the command
		if err := c.Start(); err != nil {
			errs <- err
			return
		}

		// read all the errs and send them on
		errwg.Add(1)
		go func() {
			defer errwg.Done()
			for errScanner.Scan() {
				errs <- fmt.Errorf(errScanner.Text()) // Println will add back the final '\n'
			}
		}()

		// read all the lines and send down stream
		readwg.Add(1)
		go func() {
			defer readwg.Done()
			readAndSend(reader, out, errs)
		}()

		// read all the messages and write to stdin
		if !e.NoStdin {
			fmt.Fprintf(stdin, "%s\n", message.String(msg))
		}
		stdin.Close()

		// wait for the reads to finish before calling c.Wait()
		readwg.Wait()

		// wait for close
		if err := c.Wait(); err != nil {
			errs <- err
		}
	}

	// wait for the errors to finish
	errwg.Wait()
}

// Commander defines what it takes to be a command
type Commander interface {
	Command() string
	Arguments() []string
}

func readAndSend(reader *bufio.Reader, out chan<- interface{}, errs chan<- error) {
	msg := []byte("")
	line, prefix, err := reader.ReadLine()
	for line != nil && err == nil {
		msg = append(msg, line...)
		if prefix {
			line, prefix, err = reader.ReadLine()
			continue
		}
		out <- bytes.NewBuffer(msg)
		msg = []byte("")
		line, prefix, err = reader.ReadLine()
	}
	if err != nil {
		if err != io.EOF {
			errs <- err
		}
	}
}
