package key

import "time"

type Expirer interface {
	Expired() bool
}

type Key struct {
	name       string
	expiration time.Time
}

func (k *Key) Expired() bool {
	return k.expiration.Before(time.Now())
}

func New(name string) *Key {
	return &Key{
		name:       name,
		expiration: time.Now().Add(10 * time.Minute),
	}
}
