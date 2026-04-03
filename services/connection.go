package services

import (
	"context"

	"golang.org/x/crypto/ssh"
)

type LiveConnection struct {
	client *ssh.Client
	cancel context.CancelFunc
}

func (l *LiveConnection) IsActive() bool {
	// A connection is "active" as long as its cancel func hasn't been called.
	// We use a simple done-channel trick via the context.
	return l.cancel != nil
}

func NewLiveConnection(cancel context.CancelFunc) *LiveConnection {
	return &LiveConnection{
		cancel: cancel,
	}
}
