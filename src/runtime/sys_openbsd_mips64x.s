// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd
// +build mips64 mips64le

//
// System calls and other sys.stuff for mips64, OpenBSD
// /usr/src/sys/kern/syscalls.master for syscall numbers.
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

#define CLOCK_REALTIME	$0
#define	CLOCK_MONOTONIC	$3

// Exit the entire program (like C exit)
TEXT runtime·exit(SB),NOSPLIT|NOFRAME,$0
	MOVW	code+0(FP), R4		// arg 1 - status
	MOVV	$1, R2			// sys_exit
	SYSCALL
	BEQ	R7, 3(PC)
	MOVV	$0, R2			// crash on syscall failure
	MOVV	R2, (R2)
	RET

// func exitThread(wait *uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0
	MOVV	wait+0(FP), R4		// arg 1 - notdead
	MOVV	$302, R2		// sys___threxit
	SYSCALL
	MOVV	$0, R2			// crash on syscall failure
	MOVV	R2, (R2)
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT|NOFRAME,$0
	MOVV	name+0(FP), R4		// arg 1 - path
	MOVW	mode+8(FP), R5		// arg 2 - mode
	MOVW	perm+12(FP), R6		// arg 3 - perm
	MOVV	$5, R2			// sys_open
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+16(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT|NOFRAME,$0
	MOVW	fd+0(FP), R4		// arg 1 - fd
	MOVV	$6, R2			// sys_close
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+8(FP)
	RET

TEXT runtime·read(SB),NOSPLIT|NOFRAME,$0
	MOVW	fd+0(FP), R4		// arg 1 - fd
	MOVV	p+8(FP), R5		// arg 2 - buf
	MOVW	n+16(FP), R6		// arg 3 - nbyte
	MOVV	$3, R2			// sys_read
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+24(FP)
	RET

// func pipe() (r, w int32, errno int32)
TEXT runtime·pipe(SB),NOSPLIT|NOFRAME,$0-12
	MOVV	R29, R4
	MOVW	$0, R5
	MOVV	$101, R2		// sys_pipe2
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, errno+8(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT|NOFRAME,$0-20
	ADD	$8, R29, R4
	MOVW	flags+0(FP), R5
	MOVV	$101, R2		// sys_pipe2
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, errno+16(FP)
	RET

TEXT runtime·write1(SB),NOSPLIT|NOFRAME,$0
	MOVV	fd+0(FP), R4		// arg 1 - fd
	MOVV	p+8(FP), R5		// arg 2 - buf
	MOVW	n+16(FP), R6		// arg 3 - nbyte
	MOVV	$4, R2			// sys_write
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+24(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$24-4
	MOVWU	usec+0(FP), R2
	MOVV	R2, R5
	MOVW	$1000000, R4
	DIVVU	R4, R2
	MOVV	R2, 8(R29)		// tv_sec
	MUL	R2, R4
	SUB	R4, R5
	MOVW	$1000, R4
	MUL	R4, R5
	MOVV	R5, 16(R29)		// tv_nsec

	ADD	$8, R29, R4		// arg 1 - rqtp
	MOVV	$0, R5			// arg 2 - rmtp
	MOVV	$91, R2			// sys_nanosleep
	SYSCALL
	RET

TEXT runtime·getthrid(SB),NOSPLIT,$0-4
	MOVV	$299, R2		// sys_getthrid
	SYSCALL
	MOVW	R2, ret+0(FP)
	RET

TEXT runtime·thrkill(SB),NOSPLIT,$0-16
	MOVW	tid+0(FP), R4		// arg 1 - tid
	MOVV	sig+8(FP), R5		// arg 2 - signum
	MOVW	$0, R6			// arg 3 - tcb
	MOVV	$119, R2		// sys_thrkill
	SYSCALL
	RET

TEXT runtime·raiseproc(SB),NOSPLIT,$0
	MOVV	$20, R4			// sys_getpid
	SYSCALL
	MOVV	R2, R4			// arg 1 - pid
	MOVW	sig+0(FP), R5		// arg 2 - signum
	MOVV	$122, R2		// sys_kill
	SYSCALL
	RET

TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVV	addr+0(FP), R4		// arg 1 - addr
	MOVV	n+8(FP), R5		// arg 2 - len
	MOVW	prot+16(FP), R6		// arg 3 - prot
	MOVW	flags+20(FP), R7	// arg 4 - flags
	MOVW	fd+24(FP), R8		// arg 5 - fd
	MOVW	$0, R9			// arg 6 - pad
	MOVW	off+28(FP), R10		// arg 7 - offset
	MOVV	$197, R2		// sys_mmap
	SYSCALL
	MOVV	$0, R4
	BEQ	R7, 3(PC)
	MOVV	R2, R4			// if error, move to R4
	MOVV	$0, R2
	MOVV	R2, p+32(FP)
	MOVV	R4, err+40(FP)
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVV	addr+0(FP), R4		// arg 1 - addr
	MOVV	n+8(FP), R5		// arg 2 - len
	MOVV	$73, R2			// sys_munmap
	SYSCALL
	BEQ	R7, 3(PC)
	MOVV	$0, R2			// crash on syscall failure
	MOVV	R2, (R2)
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0
	MOVV	addr+0(FP), R4		// arg 1 - addr
	MOVV	n+8(FP), R5		// arg 2 - len
	MOVW	flags+16(FP), R6	// arg 2 - flags
	MOVV	$75, R2			// sys_madvise
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+24(FP)
	RET

TEXT runtime·setitimer(SB),NOSPLIT,$0
	MOVW	mode+0(FP), R4		// arg 1 - mode
	MOVV	new+8(FP), R5		// arg 2 - new value
	MOVV	old+16(FP), R6		// arg 3 - old value
	MOVV	$69, R2			// sys_setitimer
	SYSCALL
	RET

// func walltime1() (sec int64, nsec int32)
TEXT runtime·walltime1(SB), NOSPLIT, $32
	MOVW	CLOCK_REALTIME, R4	// arg 1 - clock_id
	MOVV	$8(R29), R5		// arg 2 - tp
	MOVV	$87, R2			// sys_clock_gettime
	SYSCALL

	MOVV	8(R29), R4		// sec
	MOVV	16(R29), R5		// nsec
	MOVV	R4, sec+0(FP)
	MOVW	R5, nsec+8(FP)

	RET

// int64 nanotime1(void) so really
// void nanotime1(int64 *nsec)
TEXT runtime·nanotime1(SB),NOSPLIT,$32
	MOVW	CLOCK_MONOTONIC, R4	// arg 1 - clock_id
	MOVV	$8(R29), R5		// arg 2 - tp
	MOVV	$87, R2			// sys_clock_gettime
	SYSCALL

	MOVW	8(R29), R3		// sec
	MOVW	16(R29), R5		// nsec

	MOVV	$1000000000, R4
	MUL	R4, R3
	ADD	R5, R3
	MOVV	R3, ret+0(FP)
	RET

TEXT runtime·sigaction(SB),NOSPLIT,$0
	MOVW	sig+0(FP), R4		// arg 1 - signum
	MOVV	new+8(FP), R5		// arg 2 - new sigaction
	MOVV	old+16(FP), R6		// arg 3 - old sigaction
	MOVV	$46, R2			// sys_sigaction
	SYSCALL
	BEQ	R7, 3(PC)
	MOVV	$3, R2			// crash on syscall failure
	MOVV	R2, (R2)
	RET

TEXT runtime·obsdsigprocmask(SB),NOSPLIT,$0
	MOVW	how+0(FP), R4		// arg 1 - mode
	MOVW	new+4(FP), R5		// arg 2 - new
	MOVV	$48, R2			// sys_sigprocmask
	SYSCALL
	BEQ	R7, 3(PC)
	MOVV	$3, R2			// crash on syscall failure
	MOVV	R2, (R2)
	MOVW	R2, ret+8(FP)
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVW	sig+8(FP), R4
	MOVV	info+16(FP), R5
	MOVV	ctx+24(FP), R6
	MOVV	fn+0(FP), R7
	CALL	(R7)			// Alignment for ELF ABI?
	RET

// TODO(jsing): Fix/sync for mips64.
TEXT runtime·sigtramp(SB),NOSPLIT,$192
	// Save callee-save registers in the case of signal forwarding.
	// Please refer to https://golang.org/issue/31827 .
	//MOVV	R19, 8*4(R29)
	//MOVV	R20, 8*5(R29)
	//MOVV	R21, 8*6(R29)
	//MOVV	R22, 8*7(R29)
	//MOVV	R23, 8*8(R29)
	//MOVV	R24, 8*9(R29)
	//MOVV	R25, 8*10(R29)
	//MOVV	R26, 8*11(R29)
	//MOVV	R27, 8*12(R29)
	//MOVV	g, 8*13(R29)
	//MOVV	R29, 8*14(R29)
	//FMOVV	F8, 8*15(R29)
	//FMOVV	F9, 8*16(R29)
	//FMOVV	F10, 8*17(R29)
	//FMOVV	F11, 8*18(R29)
	//FMOVV	F12, 8*19(R29)
	//FMOVV	F13, 8*20(R29)
	//FMOVV	F14, 8*21(R29)
	//FMOVV	F15, 8*22(R29)

	// If called from an external code context, g will not be set.
	// Save R0, since runtime·load_g will clobber it.
	MOVW	R2, 8(R29)		// signum
	MOVB	runtime·iscgo(SB), R0
	BEQ	R2, 2(PC)
	CALL	runtime·load_g(SB)

	MOVV	R1, 16(R29)
	MOVV	R2, 24(R29)
	CALL	runtime·sigtrampgo(SB)

	// Restore callee-save registers.
	//MOVV	8*4(R29), R19
	//MOVV	8*5(R29), R20
	//MOVV	8*6(R29), R21
	//MOVV	8*7(R29), R22
	//MOVV	8*8(R29), R23
	//MOVV	8*9(R29), R24
	//MOVV	8*10(R29), R25
	//MOVV	8*11(R29), R26
	//MOVV	8*12(R29), R27
	//MOVV	8*13(R29), g
	//MOVV	8*14(R29), R29
	//FMOVV	8*15(R29), F8
	//FMOVV	8*16(R29), F9
	//FMOVV	8*17(R29), F10
	//FMOVV	8*18(R29), F11
	//FMOVV	8*19(R29), F12
	//FMOVV	8*20(R29), F13
	//FMOVV	8*21(R29), F14
	//FMOVV	8*22(R29), F15

	RET

// int32 tfork(void *param, uintptr psize, M *mp, G *gp, void (*fn)(void));
TEXT runtime·tfork(SB),NOSPLIT,$0

	// Copy mp, gp and fn off parent stack for use by child.
	MOVV	mm+16(FP), R7
	MOVV	gg+24(FP), R8
	MOVV	fn+32(FP), R9

	MOVV	param+0(FP), R4		// arg 1 - param
	MOVV	psize+8(FP), R5		// arg 2 - psize
	MOVV	$8, R2			// sys___tfork
	SYSCALL

	// Return if syscall failed.
	BEQ	R7, 4(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+40(FP)
	RET

	// In parent, return.
	BEQ	R2, 3(PC)
	MOVW	R2, ret+40(FP)
	RET

	// Initialise m, g.
	MOVV	R8, g
	MOVV	R7, g_m(g)

	// Call fn.
	CALL	(R9)

	// fn should never return.
	MOVV	$2, R8			// crash if reached
	MOVV	R8, (R8)
	RET

TEXT runtime·sigaltstack(SB),NOSPLIT,$0
	MOVV	new+0(FP), R4		// arg 1 - new sigaltstack
	MOVV	old+8(FP), R5		// arg 2 - old sigaltstack
	MOVV	$288, R2		// sys_sigaltstack
	SYSCALL
	BEQ	R7, 3(PC)
	MOVV	$0, R8			// crash on syscall failure
	MOVV	R8, (R8)
	RET

TEXT runtime·osyield(SB),NOSPLIT,$0
	MOVV	$298, R2		// sys_sched_yield
	SYSCALL
	RET

TEXT runtime·thrsleep(SB),NOSPLIT,$0
	MOVV	ident+0(FP), R4		// arg 1 - ident
	MOVW	clock_id+8(FP), R5	// arg 2 - clock_id
	MOVV	tsp+16(FP), R6		// arg 3 - tsp
	MOVV	lock+24(FP), R7		// arg 4 - lock
	MOVV	abort+32(FP), R8	// arg 5 - abort
	MOVV	$94, R2			// sys___thrsleep
	SYSCALL
	MOVW	R2, ret+40(FP)
	RET

TEXT runtime·thrwakeup(SB),NOSPLIT,$0
	MOVV	ident+0(FP), R4		// arg 1 - ident
	MOVW	n+8(FP), R5		// arg 2 - n
	MOVV	$301, R2		// sys___thrwakeup
	SYSCALL
	MOVW	R2, ret+16(FP)
	RET

TEXT runtime·sysctl(SB),NOSPLIT,$0
	MOVV	mib+0(FP), R4		// arg 1 - mib
	MOVW	miblen+8(FP), R5	// arg 2 - miblen
	MOVV	out+16(FP), R6		// arg 3 - out
	MOVV	size+24(FP), R7		// arg 4 - size
	MOVV	dst+32(FP), R8		// arg 5 - dest
	MOVV	ndst+40(FP), R9		// arg 6 - newlen
	MOVV	$202, R2		// sys___sysctl
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+48(FP)
	RET

// int32 runtime·kqueue(void);
TEXT runtime·kqueue(SB),NOSPLIT,$0
	MOVV	$269, R2		// sys_kqueue
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+0(FP)
	RET

// int32 runtime·kevent(int kq, Kevent *changelist, int nchanges, Kevent *eventlist, int nevents, Timespec *timeout);
TEXT runtime·kevent(SB),NOSPLIT,$0
	MOVW	kq+0(FP), R4		// arg 1 - kq
	MOVV	ch+8(FP), R5		// arg 2 - changelist
	MOVW	nch+16(FP), R6		// arg 3 - nchanges
	MOVV	ev+24(FP), R7		// arg 4 - eventlist
	MOVW	nev+32(FP), R8		// arg 5 - nevents
	MOVV	ts+40(FP), R9		// arg 6 - timeout
	MOVV	$72, R2			// sys_kevent
	SYSCALL
	BEQ	R7, 2(PC)
	SUBVU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+48(FP)
	RET

// func closeonexec(fd int32)
TEXT runtime·closeonexec(SB),NOSPLIT,$0
	MOVW	fd+0(FP), R4		// arg 1 - fd
	MOVV	$2, R5			// arg 2 - cmd (F_SETFD)
	MOVV	$1, R6			// arg 3 - arg (FD_CLOEXEC)
	MOVV	$92, R2			// sys_fcntl
	SYSCALL
	RET

// func runtime·setNonblock(int32 fd)
TEXT runtime·setNonblock(SB),NOSPLIT|NOFRAME,$0-4
	MOVW	fd+0(FP), R4		// arg 1 - fd
	MOVV	$3, R5			// arg 2 - cmd (F_GETFL)
	MOVV	$0, R6			// arg 3
	MOVV	$92, R2			// sys_fcntl
	SYSCALL
	MOVV	$0x800, R6		// O_NONBLOCK
	OR	R2, R6			// arg 3 - flags
	MOVW	fd+0(FP), R4		// arg 1 - fd
	MOVV	$4, R5			// arg 2 - cmd (F_SETFL)
	MOVV	$92, R2			// sys_fcntl
	SYSCALL
	RET
