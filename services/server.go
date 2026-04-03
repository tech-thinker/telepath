package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"golang.org/x/crypto/ssh"
)

type Server interface {
	StartAll(ctx context.Context, wg *sync.WaitGroup)
	StopAll()
}

type srv struct {
	cfg             []models.Config
	liveConnections map[string]*LiveConnection
	mu              sync.Mutex
}

func NewServer(cfg []models.Config) Server {
	return &srv{
		cfg:             cfg,
		liveConnections: make(map[string]*LiveConnection),
	}
}

func (s *srv) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, liveConn := range s.liveConnections {
		liveConn.cancel()
	}
}

// StartAll launches every tunnel concurrently.
// It returns as soon as all goroutines have been spawned; the caller must
// wait on wg for them to finish.
func (s *srv) StartAll(ctx context.Context, wg *sync.WaitGroup) {
	for _, cfg := range s.cfg {
		cfg := cfg // capture loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.start(ctx, cfg)
		}()
	}
}

func (s *srv) start(ctx context.Context, cfg models.Config) {
	// Create a child context so this tunnel can be stopped individually.
	tCtx, cancel := context.WithCancel(ctx)

	s.mu.Lock()
	if liveConn, ok := s.liveConnections[cfg.Name]; ok && liveConn.IsActive() {
		s.mu.Unlock()
		cancel()
		log.Printf("[%s] tunnel already active, skipping", cfg.Name)
		return
	}
	liveConn := NewLiveConnection(cancel)
	s.liveConnections[cfg.Name] = liveConn
	s.mu.Unlock()

	defer func() {
		liveConn.cancel()
		s.mu.Lock()
		delete(s.liveConnections, cfg.Name)
		s.mu.Unlock()
		log.Printf("[%s] tunnel stopped", cfg.Name)
	}()

	if cfg.Server == nil {
		log.Printf("[%s] no server attached, skipping", cfg.Name)
		return
	}

	// Build SSH client chain (jump hosts → final host)
	client, err := s.buildSSHClient(cfg)
	if err != nil {
		log.Printf("[%s] SSH setup failed: %v", cfg.Name, err)
		return
	}
	defer client.Close()

	switch cfg.Type {
	case "R":
		s.startRemote(tCtx, cfg, client)
	default: // "L" or unset — local forwarding
		s.startLocal(tCtx, cfg, client)
	}
}

// buildSSHClient dials every hop in the jump chain and returns the final client.
func (s *srv) buildSSHClient(cfg models.Config) (*ssh.Client, error) {
	var client *ssh.Client
	var err error
	for _, hop := range s.makeHostChain(cfg.Server) {
		client, err = s.createSSHClient(hop, client)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", hop.Host, err)
		}
	}
	return client, nil
}

// startLocal implements -L forwarding:
//
//	Listen on localHost:localPort → dial remoteHost:remotePort through SSH.
func (s *srv) startLocal(ctx context.Context, cfg models.Config, client *ssh.Client) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.LocalHost, cfg.LocalPort))
	if err != nil {
		log.Printf("[%s] failed to listen on %s:%d: %v", cfg.Name, cfg.LocalHost, cfg.LocalPort, err)
		return
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	log.Printf("[%s] LOCAL  %s:%d → (SSH) → %s:%d",
		cfg.Name, cfg.LocalHost, cfg.LocalPort, cfg.RemoteHost, cfg.RemotePort)

	s.acceptLoop(ctx, cfg, listener, func(conn net.Conn) {
		// dial the remote target through the SSH tunnel
		s.handleConnection(ctx, cfg, client, conn)
	})
}

// startRemote implements -R forwarding:
//
//	SSH server listens on remoteHost:remotePort → dial localHost:localPort directly.
func (s *srv) startRemote(ctx context.Context, cfg models.Config, client *ssh.Client) {
	// Ask the SSH daemon to open a listener on the remote side.
	remoteAddr := fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort)
	listener, err := client.Listen("tcp", remoteAddr)
	if err != nil {
		log.Printf("[%s] failed to open remote listener on %s: %v", cfg.Name, remoteAddr, err)
		return
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	log.Printf("[%s] REMOTE (SSH) → %s:%d → %s:%d",
		cfg.Name, cfg.RemoteHost, cfg.RemotePort, cfg.LocalHost, cfg.LocalPort)

	s.acceptLoop(ctx, cfg, listener, func(remoteConn net.Conn) {
		// dial the local target directly (no SSH needed — traffic arrives here)
		s.handleRemoteConnection(ctx, cfg, remoteConn)
	})
}

// acceptLoop runs the common Accept → dispatch pattern for both tunnel types.
func (s *srv) acceptLoop(ctx context.Context, cfg models.Config, listener net.Listener, handle func(net.Conn)) {
	var connWg sync.WaitGroup
	defer connWg.Wait()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return // context cancelled — clean exit
			}
			log.Printf("[%s] accept error: %v", cfg.Name, err)
			return
		}
		log.Printf("[%s] accepted connection from %s", cfg.Name, conn.RemoteAddr())

		connWg.Add(1)
		go func(c net.Conn) {
			defer connWg.Done()
			handle(c)
		}(conn)
	}
}

// handleConnection is used by LOCAL forwarding: dials the remote target
// through the SSH tunnel and proxies traffic.
func (s *srv) handleConnection(ctx context.Context, cfg models.Config, client *ssh.Client, local net.Conn) {
	defer local.Close()

	type dialResult struct {
		conn net.Conn
		err  error
	}
	resultCh := make(chan dialResult, 1)

	go func() {
		conn, err := client.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort))
		resultCh <- dialResult{conn: conn, err: err}
	}()

	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()

	var remote net.Conn
	select {
	case res := <-resultCh:
		if res.err != nil {
			log.Printf("[%s] failed to dial remote %s:%d: %v", cfg.Name, cfg.RemoteHost, cfg.RemotePort, res.err)
			return
		}
		remote = res.conn
	case <-timer.C:
		log.Printf("[%s] dial to remote %s:%d timed out", cfg.Name, cfg.RemoteHost, cfg.RemotePort)
		return
	case <-ctx.Done():
		return
	}
	defer remote.Close()

	pipe(local, remote)
}

// handleRemoteConnection is used by REMOTE forwarding: dials the local target
// directly and proxies traffic from the remote-initiated connection.
func (s *srv) handleRemoteConnection(ctx context.Context, cfg models.Config, remote net.Conn) {
	defer remote.Close()

	type dialResult struct {
		conn net.Conn
		err  error
	}
	resultCh := make(chan dialResult, 1)

	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.LocalHost, cfg.LocalPort))
		resultCh <- dialResult{conn: conn, err: err}
	}()

	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()

	var local net.Conn
	select {
	case res := <-resultCh:
		if res.err != nil {
			log.Printf("[%s] failed to dial local %s:%d: %v", cfg.Name, cfg.LocalHost, cfg.LocalPort, res.err)
			return
		}
		local = res.conn
	case <-timer.C:
		log.Printf("[%s] dial to local %s:%d timed out", cfg.Name, cfg.LocalHost, cfg.LocalPort)
		return
	case <-ctx.Done():
		return
	}
	defer local.Close()

	pipe(remote, local)
}

// pipe copies traffic bidirectionally between two connections.
// It returns as soon as one side closes.
func pipe(a, b net.Conn) {
	done := make(chan struct{}, 2)
	copy := func(dst, src net.Conn) {
		io.Copy(dst, src) //nolint:errcheck
		done <- struct{}{}
	}
	go copy(b, a)
	go copy(a, b)
	<-done
}

func (s *srv) makeHostChain(server *models.Server) []*models.Server {
	chain := []*models.Server{}
	for server != nil {
		chain = append(chain, server)
		server = server.Jump
	}
	// reverse: jump host first, final host last
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain
}

func (s *srv) createSSHClient(server *models.Server, proxy *ssh.Client) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	if server.AuthType == constants.CREDIENTIAL_PASS {
		authMethods = append(authMethods, ssh.Password(server.Password))
	}

	if server.AuthType == constants.CREDIENTIAL_KEY {
		key, err := os.ReadFile(server.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}
		var signer ssh.Signer
		if len(server.Passphrase) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(server.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	clientConfig := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", server.Host, server.Port)
	if proxy != nil {
		conn, err := proxy.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s via proxy: %w", address, err)
		}
		clientConn, chans, reqs, err := ssh.NewClientConn(conn, address, clientConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create SSH connection to %s: %w", address, err)
		}
		client := ssh.NewClient(clientConn, chans, reqs)
		startKeepalive(client)
		return client, nil
	}

	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	startKeepalive(client)
	return client, nil
}

func startKeepalive(client *ssh.Client) {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			_, _, err := client.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				client.Close()
				return
			}
		}
	}()
}
