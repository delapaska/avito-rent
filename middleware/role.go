package middleware

import "fmt"

func RoleGuard(userType string) error {
	if userType != "moderator" && userType != "client" {
		return fmt.Errorf("invalid role: %s found, moderator or client required", userType)
	}
	return nil
}
