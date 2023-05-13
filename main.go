package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
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

	SERIAL_READ_TIMEOUT = "500ms"

	COMMAND_LIST_DIR = 0x03
	COMMAND_DOWNLOAD = 0x05
	COMMAND_UPLOAD   = 0x06
	COMMAND_REMOVE   = 0x07

	BASE_DIR = "/flash"

	MAX_CHUNK_SIZE = 256 // 2**8
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
		command, args := parseCommandLine(line)
		switch command {
		case "":
			// ignore
		case "ls":
			if err := listDir(port, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		case "put":
			if err := upload(port, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
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

	fmt.Println("ls [PATH]")
	fmt.Println("  list directory")

	fmt.Println("put src dst")
	fmt.Println("  upload local file to device")

	fmt.Println("exit")
	fmt.Println("  exit application")
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

func writeAndRead(port serial.Port, b []byte) ([]byte, error) {
	b = appendCrc(b)
	if _, err := port.Write(b); err != nil {
		return nil, err
	}
	buff := make([]byte, 256)
	result := make([]byte, 0, 512)
	for {
		n, err := port.Read(buff)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}
		result = append(result, buff[:n]...)
	}
	if !verifyReceivedContainer(result) {
		return nil, fmt.Errorf("invalid response")
	}
	data, err := extractReceivedData(result)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func createCommand(command int8, data string) []byte {
	buff := make([]byte, 1, 64)
	buff[0] = byte(command)
	buff = append(buff, data...)
	return buff
}

func verifyReceivedContainer(data []byte) bool {
	head := [...]byte{0xaa, 0xab, 0xaa}
	foot := [...]byte{0xab, 0xcc, 0xab}
	return reflect.DeepEqual(data[0:3], head[:]) && reflect.DeepEqual(data[len(data)-3:], foot[:])
}

func extractReceivedData(data []byte) ([]byte, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("invalid data received")
	}
	content := data[5 : len(data)-5]
	if data[4] != byte(0x00) {
		return nil, fmt.Errorf("error: %v", string(content))
	}
	return content, nil
}

func crc16(data []byte) uint16 {
	// Modbus
	crc := uint16(0xffff)
	for i := 0; i < len(data); i++ {
		crc ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			f := crc & 1
			crc >>= 1
			if f > 0 {
				crc ^= 0xa001
			}
		}
	}
	return crc
}

func appendCrc(data []byte) []byte {
	crc := crc16(data)
	return append(data, byte(crc>>8), byte(crc&0x00ff))
}

func normalizePath(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("absolute path is not permitted: %v", path)
	}
	ret := filepath.Clean(BASE_DIR + "/" + path)
	if !strings.HasPrefix(ret, BASE_DIR) {
		return "", fmt.Errorf("forbidden path: %v", path)
	}
	return ret, nil
}

func listDir(port serial.Port, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}
	var path string
	if len(args) == 0 {
		path = BASE_DIR
	} else {
		var err error
		path, err = normalizePath(args[0])
		if err != nil {
			return nil
		}
	}
	result, err := writeAndRead(port, createCommand(COMMAND_LIST_DIR, path))
	if err != nil {
		return err
	}
	for _, e := range strings.Split(string(result), ",") {
		if len(e) > 0 {
			fmt.Println(e)
		}
	}
	return nil
}

func upload(port serial.Port, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("src and dst path arguments are required")
	}
	src := args[0]
	dst, err := normalizePath(args[1])
	if err != nil {
		return err
	}
	fp, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open src file: %w", err)
	}
	defer fp.Close()
	info, err := fp.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat src file: %w", err)
	}
	fmt.Println("Uploading...")
	buff := make([]byte, MAX_CHUNK_SIZE)
	first := true
	uploaded := 0
	for {
		n, err := fp.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read src file: %w", err)
		}
		if n == 0 {
			break
		}
		data := make([]byte, 0, 300)
		data = append(data, createCommand(COMMAND_UPLOAD, dst)...)
		data = append(data, byte(0x00))
		if !first {
			data = append(data, byte(0x00))
		} else {
			data = append(data, byte(0x01))
			first = false
		}
		data = append(data, buff[:n]...)
		result, err := writeAndRead(port, data)
		if err != nil {
			return fmt.Errorf("failed to upload: %w", err)
		}
		if !strings.Contains(string(result), "done") {
			return fmt.Errorf("unexpected result")
		}
		uploaded += n
		fmt.Printf("\r%d / %d bytes", uploaded, info.Size())
	}
	fmt.Println()
	return nil
}
