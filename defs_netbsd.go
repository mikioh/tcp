// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>
*/
import "C"

const (
	sysFIONREAD  = C.FIONREAD
	sysFIONWRITE = C.FIONWRITE
	sysFIONSPACE = C.FIONSPACE
)
