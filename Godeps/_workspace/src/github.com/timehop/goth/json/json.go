package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Solves the precision loss problem when numbers
// are sometimes sent over the wire in scientific notation
// See: http://play.golang.org/p/VpnD88yhSo
type Number float64

func (n Number) Int64() int64 {
	return int64(n)
}

func (n Number) Int() int {
	return int(n)
}

func (n Number) Float64() float64 {
	return float64(n)
}

// String returns a decimal notation without padded zeros.
func (n Number) String() string {
	if n == 0 {
		return "0"
	}
	s := fmt.Sprintf("%f", n)
	s = strings.Trim(s, "0")
	splut := strings.Split(s, ".")
	if len(splut[1]) == 0 {
		return splut[0]
	}
	return s
}

func (n Number) MarshalJSON() ([]byte, error) {
	return []byte(n.String()), nil
}

// Required to enforce that string values are attempted to be parsed as numbers
func (n *Number) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return &json.InvalidUnmarshalError{
			Type: reflect.TypeOf(*n),
		}
	}

	var f float64
	var err error
	var errVal string
	switch data[0] {
	case '"': // Strings
		if f, err = strconv.ParseFloat(string(data[1:len(data)-1]), 64); err != nil {
			errVal = fmt.Sprintf("string %s", data)
		}
	case '[': // Arrays
		errVal = fmt.Sprintf("array %s", data)
	case 't', 'f': // Booleans
		errVal = fmt.Sprintf("boolean %s", data)
	default: // Objects and numbers
		if err = json.Unmarshal(data, &f); err != nil {
			errVal = fmt.Sprintf("object %s", data)
		}
	}
	if errVal != "" {
		return &json.UnmarshalTypeError{
			Value: errVal,
			Type:  reflect.TypeOf(*n),
		}
	}
	*n = Number(f)
	return nil
}

type Data map[string]interface{}

// Parses the value the expected path as a Number.
// If the value is a string then an attempt is made to convert
// it to a float64 before casting to Number.
func (s Data) Number(dotPath string) (Number, error) {
	value := s.Value(dotPath)
	if value == nil {
		return 0, MissingFieldError(dotPath)
	}

	switch val := value.(type) {
	case float64:
		return Number(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, NewTypeError(dotPath, "Number", value)
		}
		return Number(f), nil
	default:
		return 0, NewTypeError(dotPath, "Number", value)
	}
}

func (s Data) Float64(dotPath string) (float64, error) {
	value := s.Value(dotPath)
	if value == nil {
		return 0, MissingFieldError(dotPath)
	}

	if val, ok := value.(float64); ok {
		return val, nil
	} else {
		return 0, NewTypeError(dotPath, "float64", value)
	}
}

func (s Data) Bool(dotPath string) (bool, error) {
	value := s.Value(dotPath)
	if value == nil {
		return false, MissingFieldError(dotPath)
	}

	if val, ok := value.(bool); ok {
		return val, nil
	}

	return false, NewTypeError(dotPath, "bool", value)
}

func (s Data) Int64(dotPath string) (int64, error) {
	i, err := s.Float64(dotPath)
	return int64(i), err
}

func (s Data) Int(dotPath string) (int, error) {
	i, err := s.Float64(dotPath)
	return int(i), err
}

func (s Data) String(dotPath string) (string, error) {
	value := s.Value(dotPath)
	if value == nil {
		return "", MissingFieldError(dotPath)
	}

	if val, ok := value.(string); ok {
		return val, nil
	} else {
		return "", NewTypeError(dotPath, "string", value)
	}
}

func (s Data) Time(dotPath string) (*time.Time, error) {
	value := s.Value(dotPath)
	if value == nil {
		return nil, MissingFieldError(dotPath)
	}

	switch t := value.(type) {
	case string:
		// 2012-04-23T18:25:43.511Z
		if jsonTime, err := time.Parse("2006-01-02T15:04:05.000Z", t); err == nil {
			return &jsonTime, nil
		}

		if facebookEventTime1, err := time.Parse("2006-01-02T15:04:05-0700", t); err == nil {
			return &facebookEventTime1, nil
		}

		if facebookEventTime2, err := time.Parse("2006-01-02T15:04:05", t); err == nil {
			return &facebookEventTime2, nil
		}

		if facebookEventTime3, err := time.Parse("2006-01-02", t); err == nil {
			return &facebookEventTime3, nil
		}

		if unixTime, err := strconv.ParseInt(t, 10, 64); err == nil {
			ts := time.Unix(unixTime, 0)
			return &ts, nil
		}

		return nil, NewTypeError(dotPath, "time.Time", value)
	case float64:
		ts := time.Unix(int64(t), 0)
		if ts.IsZero() {
			return nil, NewTypeError(dotPath, "time.Time", value)
		}
		return &ts, nil
	case int64:
		ts := time.Unix(t, 0)
		if ts.IsZero() {
			return nil, NewTypeError(dotPath, "time.Time", value)
		}
		return &ts, nil
	case int:
		ts := time.Unix(int64(t), 0)
		if ts.IsZero() {
			return nil, NewTypeError(dotPath, "time.Time", value)
		}
		return &ts, nil
	default:
		return nil, NewTypeError(dotPath, "time.Time", value)
	}
}

func (s Data) Data(dotPath string) (Data, error) {
	value := s.Value(dotPath)
	if value == nil {
		return nil, MissingFieldError(dotPath)
	}

	if val, ok := value.(map[string]interface{}); ok {
		return Data(val), nil
	} else {
		return nil, NewTypeError(dotPath, "json.Data", value)
	}
}

func (s Data) Array(dotPath string) ([]interface{}, error) {
	value := s.Value(dotPath)
	if value == nil {
		return nil, MissingFieldError(dotPath)
	}

	if val, ok := value.([]interface{}); ok {
		return val, nil
	} else {
		return nil, NewTypeError(dotPath, "[]interface{}", value)
	}
}

func (s Data) Value(dotPath string) interface{} {
	keys := strings.Split(dotPath, ".")
	var current interface{}
	current = s
	for i := 0; i < len(keys); i++ {
		switch current.(type) {
		case Data:
			current = (current.(Data))[keys[i]]
		case map[string]interface{}:
			current = (current.(map[string]interface{}))[keys[i]]
		default:
			return nil
		}
	}
	return current
}

func (s Data) JSON() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
