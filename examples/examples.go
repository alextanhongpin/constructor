package examples

import (
	"time"

	"github.com/alextanhongpin/constructor/examples/foo"
	"github.com/google/uuid"
)

//go:generate constructor -type User -exclude=CreatedAt,Foo
type User struct {
	ID        uuid.UUID
	Name      string
	Age       int
	SocialID  string
	CreatedAt time.Time
	Foo       foo.Foo
	Hobbies   []string
	Languages map[string]bool
	Bar       *Bar
	MaritalStatus
	Permission *Permission
	Haha       map[*int]Bar
}

//go:generate go run github.com/alextanhongpin/constructor/cmd/constructor -type Bar
type Bar struct {
	Baz string
}

type MaritalStatus struct {
	Married bool
}

type Permission struct {
	Name string
}
