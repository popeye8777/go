// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// Call fn with arg as its argument. Return what fn returns.
// fn is the raw pc value of the entry point of the desired function.
// Switches to the system stack, if not already there.
// Preserves the calling point as the location where a profiler traceback will begin.
//go:nosplit
func libcCall(fn, arg unsafe.Pointer) int32 {
	// Leave caller's PC/SP/G around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = getcallerpc()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = getcallersp()
	} else {
		// Make sure we don't reset libcallsp. This makes
		// libcCall reentrant; We remember the g/pc/sp for the
		// first call on an M, until that libcCall instance
		// returns.  Reentrance only matters for signals, as
		// libc never calls back into Go.  The tricky case is
		// where we call libcX from an M and record g/pc/sp.
		// Before that call returns, a signal arrives on the
		// same M and the signal handling code calls another
		// libc function.  We don't want that second libcCall
		// from within the handler to be recorded, and we
		// don't want that call's completion to zero
		// libcallsp.
		// We don't need to set libcall* while we're in a sighandler
		// (even if we're not currently in libc) because we block all
		// signals while we're handling a signal. That includes the
		// profile signal, which is the one that uses the libcall* info.
		mp = nil
	}
	res := asmcgocall(fn, arg)
	if mp != nil {
		mp.libcallsp = 0
	}
	return res
}

// The X versions of syscall expect the libc call to return a 64-bit result.
// Otherwise (the non-X version) expects a 32-bit result.
// This distinction is required because an error is indicated by returning -1,
// and we need to know whether to check 32 or 64 bits of the result.
// (Some libc functions that return 32 bits put junk in the upper 32 bits of AX.)

//go:linkname syscall_syscall syscall.syscall
//go:nosplit
//go:cgo_unsafe_args
func syscall_syscall(fn, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	entersyscall()
	libcCall(unsafe.Pointer(funcPC(syscall)), unsafe.Pointer(&fn))
	exitsyscall()
	return
}
func syscall()

//go:linkname syscall_syscall6 syscall.syscall6
//go:nosplit
//go:cgo_unsafe_args
func syscall_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	entersyscall()
	libcCall(unsafe.Pointer(funcPC(syscall6)), unsafe.Pointer(&fn))
	exitsyscall()
	return
}
func syscall6()

//go:linkname syscall_syscall6X syscall.syscall6X
//go:nosplit
//go:cgo_unsafe_args
func syscall_syscall6X(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	entersyscall()
	libcCall(unsafe.Pointer(funcPC(syscall6X)), unsafe.Pointer(&fn))
	exitsyscall()
	return
}
func syscall6X()

//go:linkname syscall_syscallPtr syscall.syscallPtr
//go:nosplit
//go:cgo_unsafe_args
func syscall_syscallPtr(fn, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	entersyscall()
	libcCall(unsafe.Pointer(funcPC(syscallPtr)), unsafe.Pointer(&fn))
	exitsyscall()
	return
}
func syscallPtr()

//go:linkname syscall_rawSyscall syscall.rawSyscall
//go:nosplit
//go:cgo_unsafe_args
func syscall_rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	libcCall(unsafe.Pointer(funcPC(syscall)), unsafe.Pointer(&fn))
	return
}

//go:linkname syscall_rawSyscall6 syscall.rawSyscall6
//go:nosplit
//go:cgo_unsafe_args
func syscall_rawSyscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	libcCall(unsafe.Pointer(funcPC(syscall6)), unsafe.Pointer(&fn))
	return
}

// The *_trampoline functions convert from the Go calling convention to the C calling convention
// and then call the underlying libc function.  They are defined in sys_openbsd_$ARCH.s.

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_init(attr *pthreadattr) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_attr_init_trampoline)), unsafe.Pointer(&attr))
}
func pthread_attr_init_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_getstacksize(attr *pthreadattr, size *uintptr) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_attr_getstacksize_trampoline)), unsafe.Pointer(&attr))
}
func pthread_attr_getstacksize_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_setdetachstate(attr *pthreadattr, state int) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_attr_setdetachstate_trampoline)), unsafe.Pointer(&attr))
}
func pthread_attr_setdetachstate_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_create(attr *pthreadattr, start uintptr, arg unsafe.Pointer) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_create_trampoline)), unsafe.Pointer(&attr))
}
func pthread_create_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_self() (t pthread) {
	libcCall(unsafe.Pointer(funcPC(pthread_self_trampoline)), unsafe.Pointer(&t))
	return
}
func pthread_self_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_kill(t pthread, sig uint32) {
	libcCall(unsafe.Pointer(funcPC(pthread_kill_trampoline)), unsafe.Pointer(&t))
}
func pthread_kill_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_mutex_init(m *pthreadmutex, attr *pthreadmutexattr) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_mutex_init_trampoline)), unsafe.Pointer(&m))
}
func pthread_mutex_init_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_mutex_lock(m *pthreadmutex) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_mutex_lock_trampoline)), unsafe.Pointer(&m))
}
func pthread_mutex_lock_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_mutex_unlock(m *pthreadmutex) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_mutex_unlock_trampoline)), unsafe.Pointer(&m))
}
func pthread_mutex_unlock_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_cond_init(c *pthreadcond, attr *pthreadcondattr) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_cond_init_trampoline)), unsafe.Pointer(&c))
}
func pthread_cond_init_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_cond_wait(c *pthreadcond, m *pthreadmutex) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_cond_wait_trampoline)), unsafe.Pointer(&c))
}
func pthread_cond_wait_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_cond_timedwait(c *pthreadcond, m *pthreadmutex, t *timespec) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_cond_timedwait_trampoline)), unsafe.Pointer(&c))
}
func pthread_cond_timedwait_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_cond_signal(c *pthreadcond) int32 {
	return libcCall(unsafe.Pointer(funcPC(pthread_cond_signal_trampoline)), unsafe.Pointer(&c))
}
func pthread_cond_signal_trampoline()

//go:nosplit
//go:cgo_unsafe_args
//
// This is exported via linkname to assembly in runtime/cgo.
//go:linkname exit
func exit(code int32) {
	libcCall(unsafe.Pointer(funcPC(exit_trampoline)), unsafe.Pointer(&code))
}
func exit_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func raiseproc(sig uint32) {
	libcCall(unsafe.Pointer(funcPC(raiseproc_trampoline)), unsafe.Pointer(&sig))
}
func raiseproc_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func raise(sig uint32) {
	libcCall(unsafe.Pointer(funcPC(raise_trampoline)), unsafe.Pointer(&sig))
}
func raise_trampoline()

func osyield() {
	libcCall(unsafe.Pointer(funcPC(sched_yield_trampoline)), unsafe.Pointer(nil))
}
func sched_yield_trampoline()

// mmap is used to do low-level memory allocation via mmap. Don't allow stack
// splits, since this function (used by sysAlloc) is called in a lot of low-level
// parts of the runtime and callers often assume it won't acquire any locks.
// go:nosplit
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (unsafe.Pointer, int) {
	args := struct {
		addr            unsafe.Pointer
		n               uintptr
		prot, flags, fd int32
		off             uint32
		ret1            unsafe.Pointer
		ret2            int
	}{addr, n, prot, flags, fd, off, nil, 0}
	libcCall(unsafe.Pointer(funcPC(mmap_trampoline)), unsafe.Pointer(&args))
	return args.ret1, args.ret2
}
func mmap_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func munmap(addr unsafe.Pointer, n uintptr) {
	libcCall(unsafe.Pointer(funcPC(munmap_trampoline)), unsafe.Pointer(&addr))
}
func munmap_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func madvise(addr unsafe.Pointer, n uintptr, flags int32) {
	libcCall(unsafe.Pointer(funcPC(madvise_trampoline)), unsafe.Pointer(&addr))
}
func madvise_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func open(name *byte, mode, perm int32) (ret int32) {
	return libcCall(unsafe.Pointer(funcPC(open_trampoline)), unsafe.Pointer(&name))
}
func open_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func closefd(fd int32) int32 {
	return libcCall(unsafe.Pointer(funcPC(close_trampoline)), unsafe.Pointer(&fd))
}
func close_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func read(fd int32, p unsafe.Pointer, n int32) int32 {
	return libcCall(unsafe.Pointer(funcPC(read_trampoline)), unsafe.Pointer(&fd))
}
func read_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func write1(fd uintptr, p unsafe.Pointer, n int32) int32 {
	return libcCall(unsafe.Pointer(funcPC(write_trampoline)), unsafe.Pointer(&fd))
}
func write_trampoline()

func pipe() (r, w int32, errno int32) {
	return pipe2(0)
}

func pipe2(flags int32) (r, w int32, errno int32) {
	var p [2]int32
	args := struct {
		p     unsafe.Pointer
		flags int32
	}{noescape(unsafe.Pointer(&p)), flags}
	errno = libcCall(unsafe.Pointer(funcPC(pipe2_trampoline)), unsafe.Pointer(&args))
	return p[0], p[1], errno
}
func pipe2_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func setitimer(mode int32, new, old *itimerval) {
	libcCall(unsafe.Pointer(funcPC(setitimer_trampoline)), unsafe.Pointer(&mode))
}
func setitimer_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func usleep(usec uint32) {
	libcCall(unsafe.Pointer(funcPC(usleep_trampoline)), unsafe.Pointer(&usec))
}
func usleep_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sysctl(mib *uint32, miblen uint32, out *byte, size *uintptr, dst *byte, ndst uintptr) int32 {
	return libcCall(unsafe.Pointer(funcPC(sysctl_trampoline)), unsafe.Pointer(&mib))
}
func sysctl_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func fcntl(fd, cmd, arg int32) int32 {
	return libcCall(unsafe.Pointer(funcPC(fcntl_trampoline)), unsafe.Pointer(&fd))
}
func fcntl_trampoline()

//go:nosplit
func nanotime1() int64 {
	var ts timespec
	args := struct {
		clock_id int32
		tp       unsafe.Pointer
	}{_CLOCK_MONOTONIC, unsafe.Pointer(&ts)}
	libcCall(unsafe.Pointer(funcPC(clock_gettime_trampoline)), unsafe.Pointer(&args))
	return ts.tv_sec*1e9 + ts.tv_nsec
}
func clock_gettime_trampoline()

//go:nosplit
func walltime1() (int64, int32) {
	var ts timespec
	args := struct {
		clock_id int32
		tp       unsafe.Pointer
	}{_CLOCK_REALTIME, unsafe.Pointer(&ts)}
	libcCall(unsafe.Pointer(funcPC(clock_gettime_trampoline)), unsafe.Pointer(&args))
	return ts.tv_sec, int32(ts.tv_nsec)
}

//go:nosplit
//go:cgo_unsafe_args
func kqueue() int32 {
	return libcCall(unsafe.Pointer(funcPC(kqueue_trampoline)), nil)
}
func kqueue_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func kevent(kq int32, ch *keventt, nch int32, ev *keventt, nev int32, ts *timespec) int32 {
	return libcCall(unsafe.Pointer(funcPC(kevent_trampoline)), unsafe.Pointer(&kq))
}
func kevent_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sigaction(sig uint32, new *sigactiont, old *sigactiont) {
	libcCall(unsafe.Pointer(funcPC(sigaction_trampoline)), unsafe.Pointer(&sig))
}
func sigaction_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sigprocmask(how uint32, new *sigset, old *sigset) {
	libcCall(unsafe.Pointer(funcPC(sigprocmask_trampoline)), unsafe.Pointer(&how))
}
func sigprocmask_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sigaltstack(new *stackt, old *stackt) {
	libcCall(unsafe.Pointer(funcPC(sigaltstack_trampoline)), unsafe.Pointer(&new))
}
func sigaltstack_trampoline()

// Not used on OpenBSD, but must be defined.
func exitThread(wait *uint32) {
}

//go:nosplit
func closeonexec(fd int32) {
	fcntl(fd, _F_SETFD, _FD_CLOEXEC)
}

//go:nosplit
func setNonblock(fd int32) {
	flags := fcntl(fd, _F_GETFL, 0)
	fcntl(fd, _F_SETFL, flags|_O_NONBLOCK)
}

// Tell the linker that the libc_* functions are to be found
// in a system library, with the libc_ prefix missing.

// TODO(jsing): This needs to be libpthread.so/libc.so without the
// version... however the seenlib in cmd/link/internal/ld does not
// account for libc.so being the same as libc.so.96.0, hence we get
// duplicate NEEDED entries.

//go:cgo_import_dynamic libc_pthread_attr_init pthread_attr_init "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_attr_getstacksize pthread_attr_getstacksize "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_attr_setdetachstate pthread_attr_setdetachstate "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_create pthread_create "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_sigmask pthread_sigmask "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_self pthread_self "libpthread.so.26.1"
//go:cgo_import_dynamic libc_pthread_kill pthread_kill "libpthread.so.26.1"

//go:cgo_import_dynamic libc_pthread_mutex_init pthread_mutex_init "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_mutex_lock pthread_mutex_lock "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_mutex_unlock pthread_mutex_unlock "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_cond_init pthread_cond_init "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_cond_wait pthread_cond_wait "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_cond_timedwait pthread_cond_timedwait "libc.so.96.0"
//go:cgo_import_dynamic libc_pthread_cond_signal pthread_cond_signal "libc.so.96.0"

//go:cgo_import_dynamic libc_errno __errno "libc.so.96.0"
//go:cgo_import_dynamic libc_exit exit "libc.so.96.0"
//go:cgo_import_dynamic libc_raise raise "libc.so.96.0"
//go:cgo_import_dynamic libc_sched_yield sched_yield "libc.so.96.0"

//go:cgo_import_dynamic libc_mmap mmap "libc.so.96.0"
//go:cgo_import_dynamic libc_munmap munmap "libc.so.96.0"
//go:cgo_import_dynamic libc_madvise madvise "libc.so.96.0"

//go:cgo_import_dynamic libc_open open "libc.so.96.0"
//go:cgo_import_dynamic libc_close close "libc.so.96.0"
//go:cgo_import_dynamic libc_read read "libc.so.96.0"
//go:cgo_import_dynamic libc_write write "libc.so.96.0"
//go:cgo_import_dynamic libc_pipe2 pipe2 "libc.so.96.0"

//go:cgo_import_dynamic libc_clock_gettime clock_gettime "libc.so.96.0"
//go:cgo_import_dynamic libc_setitimer setitimer "libc.so.96.0"
//go:cgo_import_dynamic libc_usleep usleep "libc.so.96.0"
//go:cgo_import_dynamic libc_sysctl sysctl "libc.so.96.0"
//go:cgo_import_dynamic libc_fcntl fcntl "libc.so.96.0"
//go:cgo_import_dynamic libc_getpid getpid "libc.so.96.0"
//go:cgo_import_dynamic libc_kill kill "libc.so.96.0"
//go:cgo_import_dynamic libc_kqueue kqueue "libc.so.96.0"
//go:cgo_import_dynamic libc_kevent kevent "libc.so.96.0"

//go:cgo_import_dynamic libc_sigaction sigaction "libc.so.96.0"
//go:cgo_import_dynamic libc_sigaltstack sigaltstack "libc.so.96.0"

//go:cgo_import_dynamic _ _ "libpthread.so.26.1"
//go:cgo_import_dynamic _ _ "libc.so.96.0"
