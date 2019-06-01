/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
	"fmt"
	"os"
	"strings"
)

const ESC = 27

var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func clearLines(lineCount int) {
	_, _ = fmt.Fprint(os.Stdout, strings.Repeat(clear, lineCount))
}
