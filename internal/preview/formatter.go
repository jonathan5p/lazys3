package preview

// Formatter converts a local file path into a displayable string.
type Formatter interface {
	Format(path string) (string, error)
}
