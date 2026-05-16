package ip

import "errors"

// to find the size of stack
func (s *Stack) Size() int {
	return len(s.IPPool)
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
	s.IPPool = append(s.IPPool, ip)
}

// get ip from stack
func (s *Stack) Pop() (string, error) {
	len := len(s.IPPool)

	if s.IsEmpty() {
		return "", errors.New("Stack Empty")
	}
	ip := s.IPPool[len-1]
	s.IPPool = s.IPPool[:len-1]
	return ip, nil

}
