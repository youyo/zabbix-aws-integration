package aws_integration

import (
	"bytes"
	"regexp"
	"strings"
)

type DiscoveryData []DiscoveryItem
type DiscoveryItem map[string]string

var macroIllegalPattern = regexp.MustCompile(`[^A-Z0-9_]+`)

func (c DiscoveryData) Json() string {
	b := bytes.Buffer{}

	b.WriteString("{\n\t\"data\":[")

	for i, item := range c {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString("\n\t{")

		firstMacro := true
		for macro, val := range item {
			if firstMacro {
				firstMacro = false
			} else {
				b.WriteString(",")
			}

			b.WriteString("\n\t\t\"{#")
			b.WriteString(macroName(macro))
			b.WriteString("}\":\"")
			b.WriteString(jsonEscape(val))
			b.WriteString("\"")
		}

		b.WriteString("}")
	}

	b.WriteString("]}")

	return b.String()
}

func jsonEscape(a string) string {
	return strings.Replace(a, "\"", "\\\"", -1)
}

func macroName(name string) string {
	name = strings.ToUpper(name)
	name = strings.Replace(name, " ", "_", -1)
	name = macroIllegalPattern.ReplaceAllString(name, "")
	return name
}
