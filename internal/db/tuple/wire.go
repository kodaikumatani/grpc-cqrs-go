package tuple

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewStore,
)
