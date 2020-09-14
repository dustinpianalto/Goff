package main

type User struct {
	ID       string
	Banned   bool
	Logging  bool
	SteamID  string
	IsActive bool
	IsStaff  bool
	IsAdmin  bool
	Guilds   []string
}

type UserService interface {
	User(id string) (*User, error)
	CreateUser(u *User) error
	DeleteUser(u *User) error
	MarkUserInactive(u *User) error
	AddUserToGuild(u *User, g *Guild) error
	RemoveUserFromGuild(u *User, g *Guild) error
	UpdateUser(u *User) error
}
