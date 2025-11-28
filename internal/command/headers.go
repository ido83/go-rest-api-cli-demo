package command

import (
	"fmt"
	"strings"
)

// HeaderFlag implements flag.Value for repeated --header flags.
type HeaderFlag map[string]string

func (h *HeaderFlag) String() string {
	if *h == nil {
		return ""
	}
	parts := make([]string, 0, len(*h))
	for k, v := range *h {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "; ")
}

func (h *HeaderFlag) Set(value string) error {
	if *h == nil {
		*h = make(map[string]string)
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid header, expected 'Key: Value'")
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if key == "" {
		return fmt.Errorf("header key cannot be empty")
	}
	(*h)[key] = val
	return nil
}
