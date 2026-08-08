package share

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewHandler,
	NewCommand,
)
