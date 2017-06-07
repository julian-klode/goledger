package goledger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Decimal is a decimal type for amounts of money (2 digits after separator).
type Decimal int

// UnmarshalJSON reads a JSON number as a decimal value.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	var value json.Number
	var sign = 1
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	sections := strings.SplitN(string(value), ".", 2)
	if sections[0][0] == '-' {
		sections[0] = sections[0][1:]
		sign = -1
	}
	main, err := strconv.ParseInt(sections[0], 10, 0)
	if err != nil {
		return err
	}

	cent, err := strconv.ParseInt(sections[1], 10, 0)
	if err != nil {
		return err
	}
	if len(sections[1]) == 1 {
		cent *= 10
	} else {
		for i := 3; i < len(sections[1]); i++ {
			cent /= 10
		}
	}
	if cent < 0 || main < 0 {
		cent *= -1
	}
	*d = Decimal(sign * int(100*main+cent))
	return nil
}

// String converts the decimal into a canonical string form.
func (d Decimal) String() string {
	sign := ""
	euros := d / 100
	cents := d - euros*100

	if d < 0 {
		sign = "-"
		euros *= -1
		cents *= -1
	}
	return fmt.Sprintf("%s%d.%02d", sign, euros, cents)
}
