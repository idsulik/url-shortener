package alias

import "github.com/idsulik/url-shortener/internal/lib/random"

type Alias struct {
}

func (a *Alias) NewAlias(size int) string {
	return random.NewRandomString(size)
}

func NewAliasGenerator() *Alias {
	return &Alias{}
}
