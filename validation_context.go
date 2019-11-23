package jsonschema

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidationContext
type ValidationContext interface {
	AddContext(instancePtr, schemaPtr string)
	FinishInstanceContext()
}

// ValidationContextRequired is used as error context when one or more required properties are missing.
type ValidationContextRequired struct {
	// Missing contains JSON Pointers to all missing properties.
	Missing []string
}

func (r *ValidationContextRequired) AddContext(instancePtr, _ string) {
	for k, p := range r.Missing {
		r.Missing[k] = joinPtr(instancePtr, p)
	}
}

func (r *ValidationContextRequired) FinishInstanceContext() {
	for k, p := range r.Missing {
		if len(p) == 0 {
			r.Missing[k] = "#"
		} else {
			r.Missing[k] = "#/" + p
		}
	}
}

func validationRequiredError(properties []string) *ValidationError {
	missing := make([]string, len(properties))

	for k := range missing {
		missing[k] = strconv.Quote(properties[k])
		properties[k] = escape(properties[k])
	}

	return &ValidationError{
		SchemaPtr: "required",
		Message:   fmt.Sprintf("missing properties: %s", strings.Join(missing, ", ")),
		Context:   &ValidationContextRequired{Missing: properties},
	}
}
