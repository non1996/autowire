package model

// Identifier 类型名，包含包路径和类名
type Identifier struct {
	Package string
	Name    string
}

type Condition struct {
	Key   string
	Value string
}
