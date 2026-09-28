package validservice

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}
