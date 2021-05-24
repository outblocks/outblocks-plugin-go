package util

import "fmt"

func ActionDesc(action, target, id string, a ...interface{}) string {
	if len(a) == 0 {
		return fmt.Sprintf("%s %s '%s'", action, target, id)
	}

	if len(a) == 1 {
		return fmt.Sprintf("%s %s '%s' (%s)", action, target, id, a[0])
	}

	return fmt.Sprintf("%s %s '%s' (%s)", action, target, id, fmt.Sprintf(a[0].(string), a[1:]...))
}

func AddDesc(target, id string, a ...interface{}) string {
	return ActionDesc("add", target, id, a...)
}

func UpdateDesc(target, id string, a ...interface{}) string {
	return ActionDesc("update", target, id, a...)
}

func DeleteDesc(target, id string, a ...interface{}) string {
	return ActionDesc("delete", target, id, a...)
}
