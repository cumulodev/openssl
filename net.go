// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openssl

import (
	"errors"
	"net"
	"time"
)

type listener struct {
	net.Listener
	ctx *Ctx
}

func (l *listener) Accept() (c net.Conn, err error) {
	c, err = l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	ssl_c, err := Server(c, l.ctx)
	if err != nil {
		c.Close()
		return nil, err
	}
	return ssl_c, nil
}

// NewListener wraps an existing net.Listener such that all accepted
// connections are wrapped as OpenSSL server connections using the provided
// context ctx.
func NewListener(inner net.Listener, ctx *Ctx) net.Listener {
	return &listener{
		Listener: inner,
		ctx:      ctx}
}

// Listen is a wrapper around net.Listen that wraps incoming connections with
// an OpenSSL server connection using the provided context ctx.
func Listen(network, laddr string, ctx *Ctx) (net.Listener, error) {
	if ctx == nil {
		return nil, errors.New("no ssl context provided")
	}
	l, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	return NewListener(l, ctx), nil
}

type DialFlags int

const (
	InsecureSkipHostVerification DialFlags = 1 << iota
	DisableSNI
)

// Dial will connect to network/address and then wrap the corresponding
// underlying connection with an OpenSSL client connection using context ctx.
// If flags includes InsecureSkipHostVerification, the server certificate's
// hostname will not be checked to match the hostname in addr. Otherwise, flags
// should be 0.
//
// Dial probably won't work for you unless you set a verify location or add
// some certs to the certificate store of the client context you're using.
// This library is not nice enough to use the system certificate store by
// default for you yet.
func Dial(network, addr string, ctx *Ctx, flags DialFlags) (*Conn, error) {
	return DialWithDialer(network, addr, ctx, flags, &net.Dialer{})
}

// DialWithDialer connects to the given network address using dialer.Dial and
// then initiates a TLS handshake, returning the resulting TLS connection. Any
// timeout or deadline given in the dialer apply to connection and TLS
// handshake as a whole.
func DialWithDialer(network, addr string, ctx *Ctx, flags DialFlags, dialer *net.Dialer) (*Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx, err = NewCtx()
		if err != nil {
			return nil, err
		}
		// TODO: use operating system default certificate chain?
	}

	c, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	// We want the Timeout and Deadline values from dialer to cover the
	// whole process: TCP connection and TLS handshake.
	timeout := dialer.Timeout
	if !dialer.Deadline.IsZero() {
		deadlineTimeout := dialer.Deadline.Sub(time.Now())
		if timeout == 0 || deadlineTimeout < timeout {
			timeout = deadlineTimeout
		}
	}

	if timeout > 0 {
		c.SetDeadline(time.Now().Add(timeout))
	}

	conn, err := Client(c, ctx)
	if err != nil {
		c.Close()
		return nil, err
	}
	if flags&DisableSNI == 0 {
		err = conn.SetTlsExtHostName(host)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	err = conn.Handshake()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
