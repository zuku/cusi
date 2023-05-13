package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.bug.st/serial"
)

const (
	COMMAND = "cusi"
	VERSION = "0.0.1"

	EXIT_OK              = 0
	EXIT_LIST_PORT_ERROR = 1
	EXIT_ARGUMENT_ERROR  = 2
	EXIT_OPEN_ERROR      = 3
	EXIT_SIGINT          = 0

	SERIAL_READ_TIMEOUT = "15s"
)

func main() {
	h := flag.Bool("h", false, "display help")
	v := flag.Bool("v", false, "display version")
	l := flag.Bool("l", false, "list serial ports")
	flag.Parse()

	if *h {
		showHelp(os.Stdout)
		os.Exit(EXIT_OK)
	}
	if *v {
		showVersion()
		os.Exit(EXIT_OK)
	}
	if *l {
		err := showPorts()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(EXIT_LIST_PORT_ERROR)
		}
		os.Exit(EXIT_OK)
	}

	name := flag.Arg(0)
	if name == "" {
		showHelp(os.Stderr)
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
	fmt.Printf("Connecting to %v ... ", name)
	port, err := open(name)
	if err != nil {
		fmt.Println("")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_OPEN_ERROR)
	}
	defer port.Close()
	fmt.Println("connected")
	reader := bufio.NewReader(os.Stdin)
	go func() {
		// Ctrl+C
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGINT)
		s := <-trap
		fmt.Fprintf(os.Stderr, "Receive signal: %v", s)
		fmt.Fprintln(os.Stderr)
		os.Exit(EXIT_SIGINT)
	}()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		command, _ := parseCommandLine(line)
		switch command {
		case "":
			// ignore
		case "help":
			showSubCommandHelp()
		case "exit":
			return
		default:
			fmt.Fprintf(os.Stderr, "%v: command not found", command)
			fmt.Fprintln(os.Stderr)
		}
	}
}

func showVersion() {
	fmt.Printf("%v %v", COMMAND, VERSION)
	fmt.Println()
}

func showHelp(out io.Writer) {
	flag.CommandLine.SetOutput(out)
	fmt.Fprintln(out, "Usage: cusi [OPTION] [PORT]")
	fmt.Fprintln(out)
	flag.PrintDefaults()
	fmt.Fprintln(out)
}

func showSubCommandHelp() {
	fmt.Println("help")
	fmt.Println("  display this help")

	fmt.Println("exit")
	fmt.Println("  exit interactive mode")
}

func showPorts() error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return err
	}
	for _, port := range ports {
		fmt.Println(port)
	}
	return nil
}

func open(p string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(p, mode)
	if err != nil {
		return nil, err
	}
	t, err := time.ParseDuration(SERIAL_READ_TIMEOUT)
	if err != nil {
		return nil, err
	}
	if err := port.SetReadTimeout(t); err != nil {
		return nil, err
	}
	return port, nil
}

func parseCommandLine(line string) (string, []string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", []string{}
	}
	if len(parts) == 1 {
		return parts[0], []string{}
	}
	return parts[0], parts[1:]
}
