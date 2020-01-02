// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix_test

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestSysctlClockinfo(t *testing.T) {
	ci, err := unix.SysctlClockinfo("kern.clockrate")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tick = %v, tickadj = %v, hz = %v, profhz = %v, stathz = %v",
		ci.Tick, ci.Tickadj, ci.Hz, ci.Profhz, ci.Stathz)
}

func TestSysctlUvmexp(t *testing.T) {
	uvm, err := unix.SysctlUvmexp("vm.uvmexp")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("free = %v", uvm.Free)
}
