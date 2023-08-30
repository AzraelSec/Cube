package evaluator

import "github.com/AzraelSec/cube/pkg/object"

func isTruthy(obj object.Object) bool {
	if obj == NULL || obj == FALSE {
		return false
	}
	return true
}

func nativeBooleanMap(exp bool) *object.Boolean {
	if exp {
		return TRUE
	}
	return FALSE
}
