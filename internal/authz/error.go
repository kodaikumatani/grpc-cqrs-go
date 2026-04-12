package authz

import "errors"

var (
	ErrUnknownObjectType = errors.New("unknown object type")
	ErrUnknownRelation   = errors.New("unknown relation")
	ErrPermissionDenied  = errors.New("permission denied")
)
