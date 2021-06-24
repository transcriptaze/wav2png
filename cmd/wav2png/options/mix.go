package options

import (
	"fmt"
	"regexp"
	"strings"
)

type Mix string

func (m Mix) String() string {
	return fmt.Sprintf("%v", string(m))
}

func (m *Mix) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile(`^(1|2|1\+2)$`).FindStringSubmatch(ss)

	if match != nil && len(match) > 1 {
		*m = Mix(match[1])
		return nil
	}

	return nil
}

func (m Mix) Channels() []int {
	switch m {
	case "1":
		return []int{1}
	case "2":
		return []int{2}
	case "1+2":
		return []int{1, 2}
	}

	return []int{1}
}
