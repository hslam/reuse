// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package reuse can reuse network ports and addresses.
package reuse

import (
	"net"
	"syscall"
)

// Control is called after creating the network connection
// but before actually dialing or binding it to the operating system.
func Control(network, address string, c syscall.RawConn) (err error) {
	if SOL_SOCKET == 0 {
		return nil
	}
	return c.Control(func(fd uintptr) {
		if err = syscall.SetsockoptInt(int(fd), SOL_SOCKET, SO_REUSEADDR, 1); err != nil {
			return
		}
		if err = syscall.SetsockoptInt(int(fd), SOL_SOCKET, SO_REUSEPORT, 1); err != nil {
			return
		}
	})
}

// ListenConfig contains reuse options for listening to an address.
var ListenConfig = net.ListenConfig{
	Control: Control,
}

// Dialer contains reuse options for connecting to an address.
var Dialer = net.Dialer{Control: Control}
