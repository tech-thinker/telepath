package repository

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/tech-thinker/telepath/models"
	"golang.org/x/crypto/ssh"
)

type SocketRepo interface {
	StartConnection(host models.SSHConfig, localPort int, remoteHost string, remotePort int)
	StopConnection()
}

type socketRepo struct {
	running bool
}

func (repo *socketRepo) StopConnection() {
	repo.running = false
}

func (repo *socketRepo) StartConnection(host models.SSHConfig, localPort int, remoteHost string, remotePort int) {
	repo.running = true
	// Port forwarding details
	var client *ssh.Client
	var err error

	next := &host
	for next != nil {
		client, err = repo.createSSHClient(*next, client)
		if err != nil {
			log.Fatalf("failed to connect to %s: %v", next.Name, err)
		}
		defer client.Close()
		next = next.NextHop
	}

	// Set up local port forwarding
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		log.Fatalf("Failed to set up listener on local port %d: %v", localPort, err)
	}
	defer listener.Close()

	log.Printf("Forwarding local port %d to %s:%d via SSH", localPort, remoteHost, remotePort)

	for {
		if !repo.running {
			break
		}

		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		remoteConn, err := client.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost, remotePort))
		if err != nil {
			log.Printf("Failed to connect to remote host: %v", err)
			localConn.Close()
			continue
		}

		// Forward traffic between local and remote connections
		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(localConn, remoteConn)
		}()
		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(remoteConn, localConn)
		}()
	}
}

func (repo *socketRepo) createSSHClient(config models.SSHConfig, proxy *ssh.Client) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	// Add private key authentication if provided
	if config.KeyFile != "" {
		key, err := os.ReadFile(config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Configure the SSH client
	clientConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Establish the SSH connection
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if proxy != nil {
		// Connect through a proxy
		conn, err := proxy.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s via proxy: %w", address, err)
		}
		clientConn, chans, reqs, err := ssh.NewClientConn(conn, address, clientConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client connection: %w", err)
		}
		return ssh.NewClient(clientConn, chans, reqs), nil
	}

	// Direct connection
	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	return client, nil
}

func NewSocketRepo() SocketRepo {
	return &socketRepo{}
}
