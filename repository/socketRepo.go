package repository

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/tech-thinker/telepath/config"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"golang.org/x/crypto/ssh"
)

type SocketRepo interface {
	Start(tunnel models.Tunnel)
	Stop(tunnel models.Tunnel)
}

type socketRepo struct {
	cfg             config.Configuration
	liveConnections map[string]*liveConnection
}

type liveConnection struct {
	client *ssh.Client
	active bool
}

func NewLiveConnection() *liveConnection {
	return &liveConnection{}
}

func (l *liveConnection) Activate() {
	l.active = true
}

func (l *liveConnection) IsActive() bool {
	return l.active
}

func (l *liveConnection) Deactivate() {
	l.active = false
}

func (repo *socketRepo) Stop(tunnel models.Tunnel) {
	liveConn, ok := repo.liveConnections[tunnel.Name]
	if ok {
		if liveConn.IsActive() {
			liveConn.Deactivate()
			return
		}
	}
}

func (repo *socketRepo) Start(tunnel models.Tunnel) {
	// Take live connection if exists or create
	liveConn, ok := repo.liveConnections[tunnel.Name]
	if ok {
		if liveConn.IsActive() {
			// Connection already exists
			return
		}
	} else {
		// creating nre connection
		liveConn = NewLiveConnection()
		repo.liveConnections[tunnel.Name] = liveConn
		liveConn.Activate()
		defer liveConn.Deactivate()
	}

	// Port forwarding details
	// var client *ssh.Client
	var err error

	for _, h := range tunnel.HostChain {
		host, ok := repo.cfg.Config().Hosts[h]
		if !ok {
			continue
		}

		liveConn.client, err = repo.createSSHClient(host, liveConn.client)
		if err != nil {
			log.Fatalf("failed to connect to %s: %v", host.Name, err)
		}
		defer liveConn.client.Close()
	}

	// Set up local port forwarding
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", tunnel.LocalPort))
	if err != nil {
		log.Fatalf("Failed to set up listener on local port %d: %v", tunnel.LocalPort, err)
	}
	defer listener.Close()

	log.Printf("Forwarding local port %d to %s:%d via SSH", tunnel.LocalPort, tunnel.RemoteHost, tunnel.RemotePort)

	for {
		if !liveConn.IsActive() {
			break
		}

		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		remoteConn, err := liveConn.client.Dial("tcp", fmt.Sprintf("%s:%d", tunnel.RemoteHost, tunnel.RemotePort))
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

func (repo *socketRepo) createSSHClient(config models.HostConfig, proxy *ssh.Client) (*ssh.Client, error) {
	// read credientials
	cred, ok := repo.cfg.Config().Credientials[config.CredientialName]
	if !ok {
		return proxy, fmt.Errorf(`No crediential found for host: %s`, config.Name)
	}

	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if cred.Type == constants.CREDIENTIAL_PASS {
		authMethods = append(authMethods, ssh.Password(cred.Password))
	}

	// Add private key authentication if provided
	if cred.Type == constants.CREDIENTIAL_KEY {
		key, err := os.ReadFile(cred.KeyFile)
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

func NewSocketRepo(
	cfg config.Configuration,
) SocketRepo {
	return &socketRepo{
		cfg:             cfg,
		liveConnections: make(map[string]*liveConnection),
	}
}
