package middleware

const (
	approved = iota
	declined
	onmoderation
)

var statusMap = map[string]int{
	"approved":      approved,
	"declined":      declined,
	"on moderation": onmoderation,
}

func StatusExists(role string) bool {
	_, exists := statusMap[role]
	return exists
}
