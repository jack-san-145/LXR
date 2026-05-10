package ip

import "errors"

func (s Stack) Size() int {
	return len(s.Stack)
}

func (s Stack) IsEmpty() bool {
	if s.Size() != 0 {
		return false
	}
	return true
}

func (s Stack) Push(ip string) {
	s.Stack = append(s.Stack, ip)
}

func (s Stack) Pop() (string, error) {
	len := len(s.Stack)

	if s.IsEmpty() {
		ip := s.Stack[len-1]
		s.Stack = s.Stack[:len-1]
		return ip, nil
	}
	return "", errors.New("Stack Empty")
}
