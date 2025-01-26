package daemon

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/handler"
	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/utils"
)

type DaemonMgr interface {
	RunAsDaemon(ctx context.Context) error
	RunDaemonChild(ctx context.Context) error
	StopDaemon(ctx context.Context) error
	StatusDaemon(ctx context.Context) error
	SendCommandToDaemon(ctx context.Context, packet models.Packet) error
}

type daemonMgr struct {
	pidFilePath string
	socketPath  string
	handler     handler.Handler
}

// Check if the daemon is already running
func (ps *daemonMgr) IsDaemonRunning(ctx context.Context) bool {
	pidData, err := os.ReadFile(ps.pidFilePath)
	if err != nil {
		return false
	}

	pid := 0
	fmt.Sscanf(string(pidData), "%d", &pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// For windows, no need to check process signal
	if utils.IsWindows() {
		return true
	}

	// Check if the process is alive
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// Run as a background daemon
func (ps *daemonMgr) RunAsDaemon(ctx context.Context) error {
	// Check if the daemon is already running
	if ps.IsDaemonRunning(ctx) {
		return fmt.Errorf("daemon is already running")
	}

	pid, err := utils.FrokProcess(os.Args[0], "daemon", "start", "--daemon-child")
	if err != nil {
		// fmt.Println("Failed to start daemon process: ", err)
		return fmt.Errorf("failed to start daemon process: %v", err.Error())
	}

	fmt.Printf("Daemon started with PID: %d\n", pid)
	return nil
}

// Run the actual daemon process (child process)
func (ps *daemonMgr) RunDaemonChild(ctx context.Context) (err error) {
	// Write the PID file
	err = os.WriteFile(ps.pidFilePath, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("Failed to write PID file: %v", err)
		return err
	}
	defer os.Remove(ps.pidFilePath)

	// Set up the UNIX socket
	var listener net.Listener

	if !utils.IsWindows() {
		if _, err := os.Stat(ps.socketPath); err == nil {
			os.Remove(ps.socketPath)
		}
		listener, err = net.Listen("unix", ps.socketPath)
		if err != nil {
			log.Fatalf("Failed to create UNIX socket: %v", err)
			return err
		}
	} else {
		listener, err = net.Listen("tcp", constants.TCP_ADDR)
		if err != nil {
			log.Fatalf("Failed to create TCP socket: %v", err)
			return err
		}
	}

	defer listener.Close()

	log.Println("Daemon is running...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go ps.handleClient(ctx, conn)
	}
}

// Stop the daemon process
func (ps *daemonMgr) StopDaemon(ctx context.Context) error {
	// Check if the daemon is running
	if !ps.IsDaemonRunning(ctx) {
		fmt.Println("daemon is not running")
		return nil
	}

	// Read the PID from the file
	pidData, err := os.ReadFile(ps.pidFilePath)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidData)))
	if err != nil {
		return fmt.Errorf("invalid PID in PID file: %v", err)
	}

	// Kill the daemon process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find daemon process: %v", err)
	}

	err = process.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill daemon process: %v", err)
	}

	// Remove the PID file
	err = os.Remove(ps.pidFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove PID file: %v", err)
	}

	// Remove the PID file
	err = os.Remove(ps.socketPath)
	if err != nil {
		return fmt.Errorf("failed to remove socket file: %v", err)
	}

	fmt.Println("Daemon stopped successfully")
	return nil
}

// Status the daemon process
func (ps *daemonMgr) StatusDaemon(ctx context.Context) error {
	if ps.IsDaemonRunning(ctx) {
		pidData, err := os.ReadFile(ps.pidFilePath)
		if err != nil {
			return fmt.Errorf("failed to read PID file: %v", err)
		}

		pid, err := strconv.Atoi(strings.TrimSpace(string(pidData)))
		if err != nil {
			return fmt.Errorf("invalid PID in PID file: %v", err)
		}

		fmt.Println("Telepath daemon is running on PID: ", pid)
	} else {
		fmt.Println("Telepath daemon is stopped.")
	}
	return nil
}

func (ps *daemonMgr) SendCommandToDaemon(ctx context.Context, packet models.Packet) error {
	// Ensure the daemon is running
	if !ps.IsDaemonRunning(ctx) {
		return fmt.Errorf("daemon is not running")
	}

	var conn net.Conn
	var err error

	if !utils.IsWindows() {
		conn, err = net.Dial("unix", ps.socketPath)
		if err != nil {
			return fmt.Errorf("failed to connect to daemon: %v", err)
		}

	} else {
		conn, err = net.Dial("tcp", constants.TCP_ADDR)
		if err != nil {
			return fmt.Errorf("failed to connect to daemon: %v", err)
		}
	}
	defer conn.Close()

	data := packet.ToByte()
	packetSize := uint32(len(data))

	// Send the size header (4 bytes)
	sizeHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeHeader, packetSize)
	_, err = conn.Write(sizeHeader)
	if err != nil {
		return fmt.Errorf("failed to send size header: %v", err)
	}

	// Send the data in chunks
	const chunkSize = 4096
	for start := 0; start < len(data); start += chunkSize {
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		_, err = conn.Write(data[start:end])
		if err != nil {
			return fmt.Errorf("failed to send data chunk: %v", err)
		}
	}

	// Read the response size header (4 bytes)
	responseSizeHeader := make([]byte, 4)
	_, err = io.ReadFull(conn, responseSizeHeader)
	if err != nil {
		return fmt.Errorf("failed to read response size header: %v", err)
	}
	responseSize := binary.BigEndian.Uint32(responseSizeHeader)

	// Read the full response data in chunks
	fullResponse := make([]byte, responseSize)
	_, err = io.ReadFull(conn, fullResponse)
	if err != nil {
		return fmt.Errorf("failed to read full response: %v", err)
	}

	// Decode the response
	result, err := utils.ToPacket(fullResponse)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	fmt.Println(string(result.Data))
	return nil
}

func (ps *daemonMgr) handleClient(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Use a buffered reader for handling fragmented data
	reader := bufio.NewReader(conn)

	// Read the size header (4 bytes)
	header := make([]byte, 4)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		log.Printf("Error reading size header: %v", err)
		return
	}

	packetSize := binary.BigEndian.Uint32(header)
	log.Printf("Expected packet size: %d", packetSize)

	// Read the full packet based on the size
	packetData := make([]byte, packetSize)
	_, err = io.ReadFull(reader, packetData)
	if err != nil {
		log.Printf("Error reading packet data: %v", err)
		return
	}

	// Decode the packet
	packet, err := utils.ToPacket(packetData)
	if err != nil {
		log.Printf("Error decoding packet: %v", err)
		return
	}

	// Process the packet
	result, err := ps.handler.Handle(ctx, packet)
	if err != nil {
		log.Printf("Error executing command: %v", err)
		return
	}

	// Prepare the response
	responseData := result.ToByte()
	responseSize := uint32(len(responseData))

	// Send the size header (4 bytes)
	sizeHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeHeader, responseSize)
	_, err = conn.Write(sizeHeader)
	if err != nil {
		log.Printf("Error sending response size header: %v", err)
		return
	}

	// Send the response in chunks
	const chunkSize = 4096
	for start := 0; start < len(responseData); start += chunkSize {
		end := start + chunkSize
		if end > len(responseData) {
			end = len(responseData)
		}
		_, err = conn.Write(responseData[start:end])
		if err != nil {
			log.Printf("Error sending response data: %v", err)
			return
		}
	}
}

func NewDaemonMgr(
	pidFilePath string,
	socketPath string,
	handler handler.Handler,
) DaemonMgr {
	return &daemonMgr{
		pidFilePath: pidFilePath,
		socketPath:  socketPath,
		handler:     handler,
	}
}
