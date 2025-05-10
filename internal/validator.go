package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Validate(val any) error {
	v := reflect.ValueOf(val)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).String()
		fieldName := v.Type().Field(i).Name
		tags := v.Type().Field(i).Tag.Get("validate")
		rules := strings.Split(tags, ",")

		if tags == "" {
			continue
		}

		for _, rule := range rules {
			if err := applyRule(rule, field, fieldName); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyRule(rule, field, fieldName string) error {
	switch {
	case strings.HasPrefix(rule, "max="):
		max, err := strconv.Atoi(strings.TrimPrefix(rule, "max="))
		if err != nil {
			return err
		}
		if len(field) > max {
			return fmt.Errorf("%s should be at most %d characters", fieldName, max)
		}

	case strings.HasPrefix(rule, "min="):
		min, err := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		if err != nil {
			return err
		}
		if len(field) < min {
			return fmt.Errorf("%s should be at least %d characters", fieldName, min)
		}

	case rule == "required":
		if len(field) == 0 {
			return fmt.Errorf("%s should not be empty", fieldName)
		}

	default:
		return fmt.Errorf("invalid validation rule: %s", rule)
	}

	return nil
}
