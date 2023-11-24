package models

type KeyAlreadyExists struct{}
func (k *KeyAlreadyExists) Error() string {
	return "error: key already exists"
}