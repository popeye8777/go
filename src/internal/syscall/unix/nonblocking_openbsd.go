// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd

package unix

import (
	"syscall"
	_ "unsafe" // for go:linkname
)

func IsNonblock(fd int) (nonblocking bool, err error) {
	flag, err := fcntl(fd, syscall.F_GETFL, 0)
	if err != nil {
		return false, err
	}
	return flag&syscall.O_NONBLOCK != 0, nil
}

// Implemented in syscall/syscall_openbsd.go.
//go:linkname fcntl syscall.fcntl
func fcntl(fd int, cmd int, arg int) (int, error)
