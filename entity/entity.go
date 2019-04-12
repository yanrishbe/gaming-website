// Package entities collects data structures used both by a server and a database
package entity

// User struct represents a struct necessary for storing and changing user's data
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}
