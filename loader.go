package superlint

import "github.com/coder/flog"

// Loader describes the user-defined rule loader.
type Loader func(log *flog.Logger, set *RuleSet)
