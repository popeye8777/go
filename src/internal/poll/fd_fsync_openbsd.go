// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poll

import (
	"syscall"
	_ "unsafe" // for go:linkname
)

// Fsync wraps syscall.Fsync.
func (fd *FD) Fsync() error {
	if err := fd.incref(); err != nil {
		return err
	}
	defer fd.decref()
	return syscall.Fsync(fd.Sysfd)
}

// Implemented in syscall/syscall_openbsd.go.
//go:linkname fcntl syscall.fcntl
func fcntl(fd int, cmd int, arg int) (int, error)
