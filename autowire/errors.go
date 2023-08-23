package autowire

import (
	"fmt"
)

var (
	errMultiMatch = fmt.Errorf("multiple components meet filter condition and no instance was designated as primary")
)

func errComponentNotFound(typeName string) error {
	return fmt.Errorf("instance [%s] not found", typeName)
}

func errComponentDuplicate(name string) error {
	return fmt.Errorf("instance [%s] is duplicate", name)
}

func errPropertyScopeDuplicate(name string) error {
	return fmt.Errorf("")
}

func errValueNotFound(scope string, key string) error {
	return fmt.Errorf("property [%s/%s] not found", scope, key)
}

func errEnvNotFound(name string) error {
	return fmt.Errorf("env [%s] not found", name)
}
