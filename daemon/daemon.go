package daemon

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

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

	// Fork the process
	cmd := exec.Command(os.Args[0], "daemon", "start", "--daemon-child")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	fmt.Printf("Daemon started with PID: %d\n", cmd.Process.Pid)
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
	if _, err := os.Stat(ps.socketPath); err == nil {
		os.Remove(ps.socketPath)
	}
	listener, err := net.Listen("unix", ps.socketPath)
	if err != nil {
		log.Fatalf("Failed to create UNIX socket: %v", err)
		return err
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
		return fmt.Errorf("daemon is not running")
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

// Send commands to the daemon
func (ps *daemonMgr) SendCommandToDaemon(ctx context.Context, packet models.Packet) error {
	// Ensure the daemon is running
	if !ps.IsDaemonRunning(ctx) {
		return fmt.Errorf("daemon is not running")
	}

	conn, err := net.Dial("unix", ps.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon: %v", err)
	}
	defer conn.Close()

	data := packet.ToByte()

	// Send request to daemon
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send command: %v", err)
	}

	// Read the response from the daemon
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	result, err := utils.ToPacket(buf[:n])
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	fmt.Printf(string(result.Data))
	return nil
}

func (ps *daemonMgr) handleClient(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading command: %v", err)
		return
	}

	packet, err := utils.ToPacket(buf[:n])
	if err != nil {
		log.Printf("Error decoding packet: %v", err)
		return
	}

	result, err := ps.handler.Handle(ctx, packet)
	if err != nil {
		log.Printf("Error executing command: %v\n", err)
	}
	conn.Write(result.ToByte())
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
