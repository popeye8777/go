// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

//
// System call support for AMD64, OpenBSD
//

// Provide these function names via assembly so they are provided as ABI0,
// rather than ABIInternal.
//
// TODO(jsing): Can this be done using a //go: tag?

TEXT	·Syscall(SB),NOSPLIT,$0-56
	JMP	·syscallInternal(SB)

TEXT	·Syscall6(SB),NOSPLIT,$0-80
	JMP	·syscall6Internal(SB)

TEXT	·RawSyscall(SB),NOSPLIT,$0-56
	JMP	·rawSyscallInternal(SB)

TEXT	·RawSyscall6(SB),NOSPLIT,$0-80
	JMP	·rawSyscall6Internal(SB)

TEXT	·Syscall9(SB),NOSPLIT,$0-104
	JMP	·syscall9Internal(SB)
