package ssi

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key string) interface{}
	// Set sets the session value associated to the given key.
	Set(key string, val interface{}) Session
	SetID(id string) Session
	// Delete removes the session value associated to the given key.
	Delete(key string) Session
	// Clear deletes all values in the session.
	Clear() Session
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string) Session
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}

	Options(Options) Session

	// Save saves all sessions used during the current request.
	Save() error
}
