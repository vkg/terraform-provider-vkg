package event

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func in(vals []string) schema.SchemaValidateFunc {
	return validation.StringInSlice(vals, false)
}

func datetime() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if _, err := time.Parse("2006-01-02 15:04:05", v); err != nil {
			es = append(es, err)
			return
		}

		return
	}
}
