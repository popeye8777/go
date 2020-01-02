// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd,go1.14

package ld

import (
	"syscall"
	"unsafe"
)

func (out *OutBuf) Mmap(filesize uint64) error {
	err := out.f.Truncate(int64(filesize))
	if err != nil {
		Exitf("resize output file failed: %v", err)
	}
	out.buf, err = syscall.Mmap(int(out.f.Fd()), 0, int(filesize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FILE)
	return err
}

func (out *OutBuf) Munmap() {
	err := out.Msync()
	if err != nil {
		Exitf("msync output file failed: %v", err)
	}
	syscall.Munmap(out.buf)
	out.buf = nil
	_, err = out.f.Seek(out.off, 0)
	if err != nil {
		Exitf("seek output file failed: %v", err)
	}
}

//go:linkname msync syscall.msync
func msync(addr uintptr, length uintptr, flags int32) (err error)

func (out *OutBuf) Msync() error {
	err := msync(uintptr(unsafe.Pointer(&out.buf[0])), uintptr(len(out.buf)), syscall.MS_SYNC)
	if err != nil {
		return err
	}
	return nil
}
