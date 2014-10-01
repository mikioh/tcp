// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.2

TEXT	·socketcall(SB),4,$0-36
	JMP	syscall·socketcall(SB)
