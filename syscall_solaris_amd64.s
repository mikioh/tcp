// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

TEXT ·rtioctl(SB),4,$0
	JMP	runtime·syscall_ioctl(SB)

TEXT ·rtsysvicall6(SB),4,$0
	JMP	runtime·syscall_sysvicall6(SB)
