package api

import "github.com/thezmc/go-sinch/pkg/sinch"

func Validate(v ...sinch.Validatable) error {
	var errs Errors
	for _, v := range v {
		if err := v.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
