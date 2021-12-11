// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"testing"
)

func TestGetTextPieces(t *testing.T) {
	for _, tt := range []string{
		`Options:\x1b[2m\1b[37mMisc\1b]0m
  [[2m[37mMisc[0m]
      --config=[Locations of config files]  [0m[90mload config files from where you specified[3m[90m (default [Locations of config files]=)[0m
  -q, --quiet                               [0m[90mNo more screen output.[3m[90m [env: QUITE] (default=true)[0m
  -v, --verbose                             [0m[90mShow this help screen[3m[90m [env: VERBOSE] (default=false)[0m
[2m[37m
`,
	} {
		_ = getTextPiece(tt, 0, 1000)
	}
}
