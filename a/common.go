package a

type PrimaryAnno interface {
	isPrimary()
}

type SecondaryAnno interface {
	isSecondary()
}

type Annotations []any
