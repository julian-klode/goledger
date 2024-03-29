/*
Copyright (C) 2017 Julian Andres Klode <jak@jak-linux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package goledger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Posting describes a part of a ledger transaction
type Posting struct {
	Account    string
	Value      decimal.Decimal
	Currency   string
	AtValue    decimal.Decimal
	AtCurrency string
}

// Transaction represents a transaction in a ledger file
type Transaction struct {
	Date        time.Time
	ValutaDate  time.Time
	Description string
	Postings    []Posting
}

func renderDecimal(d decimal.Decimal) string {
	value := d.String()
	if !strings.Contains(value, ".") {
		value += ".00"
	} else if len(value) >= 2 && value[len(value)-2] == '.' {
		value += "0"
	}
	return value
}

// Print prints the ledger transaction to the writer
func (l *Transaction) Print(w io.Writer) {
	switch {
	case l.ValutaDate.Year() > 1000 && l.Date.Year() > 1000 && l.ValutaDate != l.Date:
		fmt.Printf("%d/%02d/%02d=%d/%02d/%02d", l.Date.Year(), l.Date.Month(), l.Date.Day(), l.ValutaDate.Year(), l.ValutaDate.Month(), l.ValutaDate.Day())
	case l.Date.Year() > 1000:
		fmt.Printf("%d/%02d/%02d", l.Date.Year(), l.Date.Month(), l.Date.Day())
	case l.ValutaDate.Year() > 1000:
		fmt.Printf("%d/%02d/%02d", l.ValutaDate.Year(), l.ValutaDate.Month(), l.ValutaDate.Day())
	default:
		fmt.Printf("1970/01/01")
	}
	fmt.Printf(" %s\n", l.Description)
	for _, p := range l.Postings {
		if !p.AtValue.IsZero() {
			fmt.Printf("    %s  %v %s @ %v %s\n", p.Account, renderDecimal(p.Value), p.Currency, renderDecimal(p.AtValue), p.AtCurrency)
		} else {
			fmt.Printf("    %s  %v %s\n", p.Account, renderDecimal(p.Value), p.Currency)
		}
	}
	fmt.Println()
}
