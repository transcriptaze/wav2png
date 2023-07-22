package audio

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
	ss := strings.ToUpper(s)
	match := regexp.MustCompile(`^(L|R|L\+R)$`).FindStringSubmatch(ss)

	if len(match) > 1 {
		*m = Mix(match[1])
		return nil
	}

	return nil
}

func (m Mix) Channels() []int {
	switch m {
	case "L":
		return []int{1}
	case "R":
		return []int{2}
	case "L+R":
		return []int{1, 2}
	}

	return []int{1, 2}
}
