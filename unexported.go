package svgxml

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	re "regexp"
	"strconv"
	s "strings"
)

func traverseObjectById(elem any, elemType string, id any, findFirst bool) ([]any, error) {
	if elem == nil {
		return nil, nil
	}
	var (
		r      reflect.Value
		err    error
		result []any
		idStr  string
		id_re  *re.Regexp
		ok     bool
	)

	if idStr, ok = id.(string); !ok {
		if id_re, ok = id.(*re.Regexp); !ok {
			return nil, fmt.Errorf("cannot determine type of id: %#v", id)
		}
	}

	results := make([]any, 0)

	if e, ok := elem.(SVG); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.(*SVG); ok {
		r = reflect.Indirect(reflect.ValueOf(*e))
	} else if e, ok := elem.(GroupDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.(*GroupDef); ok {
		r = reflect.Indirect(reflect.ValueOf(*e))
	} else if e, ok := elem.(PathDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.(*PathDef); ok {
		r = reflect.Indirect(reflect.ValueOf(*e))
	} else if e, ok := elem.(TextDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.(*TextDef); ok {
		r = reflect.Indirect(reflect.ValueOf(*e))
	} else if e, ok := elem.([]GroupDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.([]PathDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else if e, ok := elem.([]TextDef); ok {
		r = reflect.Indirect(reflect.ValueOf(e))
	} else {
		r = reflect.ValueOf(e)
		if r == (reflect.Value{}) {
			return results, fmt.Errorf("no reflect.Value for element")
		} else {
			return results, fmt.Errorf("element is unhandled Kind '%s' and Type '%s'", r.Kind().String(), r.Type().String())
		}
	}

	t := r.Kind().String()
	switch t {
	case "struct":
		var idCompare bool
		if id_re == nil {
			idCompare = (r.FieldByName("Id").String() == idStr)
		} else {
			idCompare = id_re.MatchString(r.FieldByName("Id").String())
		}
		if idCompare && r.Type().Name() == elemType {
			results = append(results, elem)
			if findFirst {
				return results, nil
			}
		}
		for _, field := range []string{"A", "G", "Circle", "Rect", "Path", "Text", "TSpan"} {
			if r.FieldByName(field) == (reflect.Value{}) {
				// typically occurs due to accessing a nonexistent field
				continue
			}
			result, err = traverseObjectById(r.FieldByName(field).Interface(), elemType, id, findFirst)
			if err == nil && len(result) > 0 {
				results = append(results, result...)
			}
		}
	case "slice", "array":
		if r.Len() > 0 {
			for ind := 0; ind < r.Len(); ind++ {
				switch r.Type().String() {
				case "[]svgxml.GroupDef":
					result, err = traverseObjectById(
						&(elem.([]GroupDef)[ind]),
						elemType,
						id,
						findFirst,
					)
				case "[]svgxml.TextDef":
					result, err = traverseObjectById(&(elem.([]TextDef)[ind]), elemType, id, findFirst)
				case "[]svgxml.PathDef":
					result, err = traverseObjectById(&(elem.([]PathDef)[ind]), elemType, id, findFirst)
				default:
					err = fmt.Errorf("unknown %s type %s", t, r.Type().String())
				}
				if err == nil && len(result) > 0 {
					results = append(results, result...)
				}
				if len(results) > 0 && findFirst {
					break
				}
			}
		}
	default:
		err = fmt.Errorf("object must be struct or slice, got '%s'", r.Kind().String())
	}

	return results, err
}

func anyNumberToString[N string | NumberAttr](num N, floatPrecision int) string {
	floatBits := 64
	if floatPrecision < 0 {
		floatPrecision = 0
	}
	numType := reflect.TypeOf(num)
	switch numType.String() {
	case "string":
		numStr := reflect.ValueOf(num).String()
		if floatPrecision == 0 && s.Index(numStr, ".") < 0 {
			return numStr
		}
		strFloatVal, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			return strconv.FormatFloat(strFloatVal, 'f', floatPrecision, floatBits)
		}
	case "float32":
		floatBits = 32
		fallthrough
	case "float64":
		return strconv.FormatFloat(reflect.ValueOf(num).Float(), 'f', floatPrecision, floatBits)
	case "int", "int32", "int64":
		if floatPrecision == 0 {
			return strconv.FormatInt(reflect.ValueOf(num).Int(), 10)
		} else {
			return strconv.FormatFloat(float64(reflect.ValueOf(num).Int()), 'f', floatPrecision, floatBits)
		}
	}
	return ""
}

func placeBackgroundRect(svg_in *SVG) *SVG {
	if svg_in.background == (RectDef{}) {
		return svg_in
	}
	newS := *svg_in
	g := GroupDef{Rect: []RectDef{newS.background}}
	newS.G = append([]GroupDef{g}, newS.G...)
	return &newS
}

func doWriteFile(content []byte, filename string) error {
	tmpFileExt := rand.Text()
	outFile := filepath.Clean(filename)
	tmpFilename := outFile + "." + string(tmpFileExt)

	err := os.WriteFile(tmpFilename, content, 0666)
	if err != nil {
		return fmt.Errorf("write: %w (file may need to be removed)", err.(*os.PathError))
	}
	err = os.Remove(outFile)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		err = os.Link(tmpFilename, outFile)
		if err == nil {
			err = os.Remove(tmpFilename)
		}
	}
	if err != nil {
		err = fmt.Errorf("rename temp file: %w (temp file may beed to be removed)", err.(*os.PathError))
	}
	return err
}
