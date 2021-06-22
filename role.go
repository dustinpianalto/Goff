package goff

type Role struct {
	ID             string
	IsModerator    bool
	IsAdmin        bool
	SelfAssignable bool
	AutoRole       bool
	Guild          string
}

type RoleService interface {
	Role(id string) (*Role, error)
	AddRole(r *Role) (*Role, error)
	DeleteRole(r *Role) error
	MakeSelfAssignable(r *Role) error
	RemoveSelfAssignable(r *Role) error
	MakeAutoRole(r *Role) error
	RemoveAutoRole(r *Role) error
	GetAutoRoles(gID string) ([]*Role, error)
}
