package user_info

import (
	"fmt"
)

type UserInfo struct {
	FullName string
	// fullName FullName
	Email  string
	Expiry string
}

func (info UserInfo) String() string {
	return fmt.Sprintf("Name:   %s\nEmail:  %s\nExpiry: %s", info.FullName, info.Email, info.Expiry)
}
