package factory

import (
	"github.com/ancalabrese/tldr/pkg/cmdutil"
)

func Defaults() *cmdutil.Factory {
	f := &cmdutil.Factory{
		ExecutableName: "tldr",
	}

	return f
}
