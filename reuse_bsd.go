// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package reuse

const (
	SOL_SOCKET   = 0xffff
	SO_REUSEADDR = 0x4
	SO_REUSEPORT = 0x200
)
