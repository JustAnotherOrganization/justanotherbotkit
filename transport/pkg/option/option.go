package option

type (
	Text struct {
		Value  string
		Escape bool // slack support, not implemented in current discord package.
	}
)
