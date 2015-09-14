/* helpers.go
 *
 * Copyright (c) 2015, Michael Guzelevich <mguzelevich@gmail.com>
 * Copyright (c) 2010, Kyle Lemons <kyle@kylelemons.net> (creator).
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */

// Package log4go provides level-based and highly configurable logging.
package log4go

import (
	"fmt"
	"os"
	"runtime"
)

func traceCheck() {
	pc, _, lineno, ok := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, ">>> [%s][%s][%s]\n", runtime.FuncForPC(pc).Name(), lineno, ok)
}
