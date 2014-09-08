// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.4

#include "textflag.h"

TEXT	·socketcall(SB),NOSPLIT,$0-36
	JMP	syscall·socketcall(SB)

TEXT	·rawsocketcall(SB),NOSPLIT,$0-36
	JMP	syscall·rawsocketcall(SB)
