package service

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}
