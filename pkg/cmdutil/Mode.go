package cmdutil

import (
	"fmt"
)

type Mode int

const (
	Tldr Mode = iota
	Interactive
)

var (
	modeString = [...]string{"TLDR", "INTERACTIVE"}

	modeToStringMap = map[Mode]string{
		Tldr:        modeString[Tldr],
		Interactive: modeString[Interactive],
	}

	stringToModeMap = map[string]Mode{
		modeString[Tldr]:        Tldr,
		modeString[Interactive]: Interactive,
	}
)

func WhichMode(m string) (Mode, error) {
	if mode, ok := stringToModeMap[m]; ok {
		return mode, nil
	}

	return -1, fmt.Errorf("unsupported mode %s", m)
}
