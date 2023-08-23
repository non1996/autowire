package a

type Application[_ any] []SecondaryAnno

func (Application[_]) isPrimary() {}

type Component[_ any] []SecondaryAnno

func (c Component[_]) isPrimary() {}
