// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package reuse

const (
	SOL_SOCKET   = 0
	SO_REUSEADDR = 0
	SO_REUSEPORT = 0
)
