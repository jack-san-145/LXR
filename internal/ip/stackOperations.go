package ip

import "errors"

// to find the size of stack
func (s *Stack) Size() int {
	return len(s.Stack)
}

// to check whether the stack is empty or not
func (s *Stack) IsEmpty() bool {
	if s.Size() != 0 {
		return false
	}
	return true
}

// to insert ip to the stack
func (s *Stack) Push(ip string) {
	s.Stack = append(s.Stack, ip)
}

// get ip from stack
func (s *Stack) Pop() (string, error) {
	len := len(s.Stack)

	if s.IsEmpty() {
		ip := s.Stack[len-1]
		s.Stack = s.Stack[:len-1]
		return ip, nil
	}
	return "", errors.New("Stack Empty")
}
