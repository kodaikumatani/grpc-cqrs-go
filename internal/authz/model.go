package authz

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Tuple struct {
	ID         uuid.UUID
	ObjectType ObjectType
	ObjectID   string
	Relation   Relation
	UserID     ulid.ULID
}

type ObjectType string

const (
	ObjectRecipe ObjectType = "recipe"
)

type Relation string

const (
	RelOwner  Relation = "owner"
	RelViewer Relation = "viewer"
	RelEditor Relation = "editor"
)

type Permission []Relation

var (
	PermViewRecipe   Permission = []Relation{RelViewer, RelEditor, RelOwner}
	PermEditRecipe   Permission = []Relation{RelEditor, RelOwner}
	PermDeleteRecipe Permission = []Relation{RelOwner}
	PermShareRecipe  Permission = []Relation{RelOwner}
)

func (t ObjectType) String() string {
	return string(t)
}

func NewObjectType(object string) (ObjectType, error) {
	switch object {
	case ObjectRecipe.String():
		return ObjectRecipe, nil
	default:
		return "", ErrUnknownObjectType
	}
}

func (r Relation) String() string {
	return string(r)
}

func NewRelation(relation string) (Relation, error) {
	switch relation {
	case RelOwner.String():
		return RelOwner, nil
	case RelViewer.String():
		return RelViewer, nil
	case RelEditor.String():
		return RelEditor, nil
	default:
		return "", ErrUnknownRelation
	}
}
