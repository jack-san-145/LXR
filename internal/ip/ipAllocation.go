package ip

import (
	"fmt"
	// "strconv"
	// "strings"
)

// to create the ip with all 4 octets
func MakeIp(first, second, third, fourth int) string {
	return fmt.Sprintf("%v.%v.%v.%v", first, second, third, fourth)
}
