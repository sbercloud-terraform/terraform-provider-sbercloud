package sbercloud

import (
	"fmt"
	"reflect"
	"strings"
)

func convertToStr(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func navigateValue(d interface{}, index []string, arrayIndex map[string]int) (interface{}, error) {
	for n, i := range index {
		if d == nil {
			return nil, nil
		}
		if d1, ok := d.(map[string]interface{}); ok {
			d, ok = d1[i]
			if !ok {
				msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
				return nil, fmt.Errorf("%s: '%s' may not exist", msg, i)
			}
		} else {
			msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
			return nil, fmt.Errorf("%s: Can not convert (%s) to map", msg, reflect.TypeOf(d))
		}

		if arrayIndex != nil {
			if j, ok := arrayIndex[strings.Join(index[:n+1], ".")]; ok {
				if d == nil {
					return nil, nil
				}
				if d2, ok := d.([]interface{}); ok {
					if len(d2) == 0 {
						return nil, nil
					}
					if j >= len(d2) {
						msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
						return nil, fmt.Errorf("%s: The index is out of array", msg)
					}

					d = d2[j]
				} else {
					msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
					return nil, fmt.Errorf("%s: Can not convert (%s) to array, index=%s.%v", msg, reflect.TypeOf(d), i, j)
				}
			}
		}
	}

	return d, nil
}
