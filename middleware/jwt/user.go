package jwt

// DefaultUserKey user
const DefaultUserKey = "user"

// User is for authentication
type User struct {
	ID              string
	IsAuthenticated bool
}
