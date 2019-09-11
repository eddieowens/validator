package validator

import "strings"

type ValidationErrors []error

func (v ValidationErrors) Error() string {
	errMsgs := make([]string, len(v))
	for i, e := range v {
		errMsgs[i] = e.Error()
	}
	return strings.Join(errMsgs, "\n")
}
