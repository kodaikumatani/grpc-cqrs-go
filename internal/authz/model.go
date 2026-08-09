package authz

import (
	"github.com/google/uuid"
)

type Tuple struct {
	ID         uuid.UUID
	ObjectType ObjectType
	ObjectID   string
	Relation   Relation
	UserID     string
}

func NewTuple(
	id uuid.UUID,
	objectType ObjectType,
	objectID string,
	relation Relation,
	userID string,
) *Tuple {
	return &Tuple{
		ID:         id,
		ObjectType: objectType,
		ObjectID:   objectID,
		Relation:   relation,
		UserID:     userID,
	}
}

type ObjectType string

const (
	ObjectRecipe ObjectType = "recipe"
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

type Relation string

const (
	RelOwner  Relation = "owner"
	RelViewer Relation = "viewer"
	RelEditor Relation = "editor"
)

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

type Permission int

const (
	PermView Permission = iota
	PermEdit
	PermDelete
	PermShare
)

var permissionRelations = map[ObjectType]map[Permission][]Relation{
	ObjectRecipe: {
		PermView:   {RelViewer, RelEditor, RelOwner},
		PermEdit:   {RelEditor, RelOwner},
		PermDelete: {RelOwner},
		PermShare:  {RelOwner},
	},
}
