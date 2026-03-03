package command

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewRecipe,
	NewUser,
)
