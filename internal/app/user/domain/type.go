package domain

import (
	"errors"

	"github.com/oklog/ulid/v2"
)

var ErrAlreadyExists = errors.New("user already exits")

type User struct {
	id    ulid.ULID
	name  string
	email string
}

// NewUser は業務フィールドを受け取ってユーザーを構築する。
// id は呼び出し側(app 層)が用意する。
func NewUser(id ulid.ULID, name, email string) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}

func (u *User) ID() ulid.ULID { return u.id }
func (u *User) Name() string  { return u.name }
func (u *User) Email() string { return u.email }
