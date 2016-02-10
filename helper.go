// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "errors"

var (
	errOpNoSupport     = errors.New("operation not supported")
	errInvalidConnType = errors.New("invalid conn type")
	errInvalidArgument = errors.New("invalid argument")
)

func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}
