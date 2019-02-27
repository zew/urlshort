// Package storages allows multiple implementation on how to store short URLs.
package storages

// IStorage is the interface for saving and retrieving URLs
type IStorage interface {
	Save(string) (string, error)
	Load(string) (string, error)
}
