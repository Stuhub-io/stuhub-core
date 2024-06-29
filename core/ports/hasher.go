package ports

type Hasher interface {
	GenerateSalt() string
	Hash(val, salt string) (string, error)
	Compare(val, hash, salt string) bool
}
