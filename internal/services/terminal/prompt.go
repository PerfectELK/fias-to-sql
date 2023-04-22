package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func YesNoPrompt(label string) bool {
	choices := "Y/N (yes/no)"

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		} else if s == "n" || s == "no" {
			return false
		} else {
			return false
		}
	}
}

func InputPrompt(label string) string {
	var str string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		str, _ = r.ReadString('\n')
		if str != "" {
			break
		}
	}
	return strings.TrimSpace(str)
}
