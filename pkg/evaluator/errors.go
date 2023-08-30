package evaluator

import (
	"fmt"

	"github.com/AzraelSec/cube/pkg/object"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil && obj.Type() == object.ERROR_OBJ {
		return true
	}
	return false
}
