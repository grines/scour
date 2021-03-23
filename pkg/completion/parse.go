package completion

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

func ParseCMD(cmd string, count int, h help) []string {
	parts := strings.Split(cmd, " ")
	if len(parts) != count || strings.Contains(cmd, "help") {
		fmt.Println(h.helpText)
		fmt.Println(h.infoText)
		fmt.Println(h.autocomplete)
		return nil
	}
	return parts
}

func HelpText(helpdata string, info string, autocomplete string) help {
	green := color.FgGreen.Render
	blue := color.FgBlue.Render

	h := help{
		helpText:     "---\n" + green("[help] ") + helpdata,
		infoText:     blue("[info] ") + info,
		autocomplete: green("[autocomplete] ") + autocomplete + "\n",
	}
	return h
}
