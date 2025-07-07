package texts

import (
	"fmt"
	"strings"
)

func Format(template string, data map[string]string) string {
	args := make([]string, 0, len(data)*2)
	for o, n := range data {
		args = append(args, fmt.Sprintf("{%s}", o), n)
	}

	return strings.NewReplacer(args...).Replace(template)
}

func FormatFloat(num float64) string {
	return fmt.Sprintf("%.2f", num)
}
