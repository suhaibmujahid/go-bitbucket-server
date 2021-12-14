package bitbucket

import (
	"fmt"
)

type Group struct {
	Name string `json:"name,omitempty"`
}

type GroupPermission struct {
	Group      *Group `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`
}

func (p *GroupPermission) String() string {
	return fmt.Sprintf("<GroupPermission: %s/%s>", p.Group.Name, p.Permission)
}

type UserPermission struct {
	Permission string `json:"permission,omitempty"`
	User       *User  `json:"user,omitempty"`
}

func (p *UserPermission) String() string {
	return fmt.Sprintf("<UserPermission: %s/%s>", p.User.Name, p.Permission)
}
