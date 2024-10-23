package models

import (
	"fmt"
	"io"
	"strconv"
)

type UInt32 uint32

func (u *UInt32) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case int:
		if v < 0 || v > 4294967295 {
			return fmt.Errorf("value must be a valid UInt32")
		}
		*u = UInt32(v)
		return nil
	case float64:
		if v < 0 || v > 4294967295 {
			return fmt.Errorf("value must be a valid UInt32")
		}
		*u = UInt32(v)
		return nil
	case string:
		parsedValue, err := strconv.ParseUint(v, 10, 32)
		if err != nil || parsedValue > 4294967295 {
			return fmt.Errorf("value must be a valid UInt32")
		}
		*u = UInt32(parsedValue)
		return nil
	default:
		return fmt.Errorf("value must be a valid UInt32")
	}
}

func (u UInt32) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, uint32(u))
}

func (UInt32) Type() string {
	return "UInt32"
}
