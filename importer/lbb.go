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

package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/julian-klode/goledger"
)

// lbbTransaction describes transactions as contained in a CSV file exported
// by the LBB (Landesbank Berlin).
type lbbTransaction struct {
	CardNumber     string
	valutaDate     time.Time
	date           time.Time
	Merchant       string
	OriginalAmount goledger.Decimal
	ExchangeRate   float64
	amount         goledger.Decimal
	currency       string
}

func (t lbbTransaction) ID() string {
	return hashTransaction(t)
}

// LocalAccount returns an ID of the local account.
func (t lbbTransaction) LocalAccount() string {
	return t.CardNumber
}

// RemoteAccount returns an ID of the remote account (IBAN).
func (t lbbTransaction) RemoteAccount() string {
	return ""
}

// RemoteName returns a name of the other account.
func (t lbbTransaction) RemoteName() string {
	return t.Merchant
}

// ReferenceText returns a description of the transaction.
func (t lbbTransaction) ReferenceText() string {
	return ""
}

// Amount returns the amount of the transaction.
func (t lbbTransaction) Amount() goledger.Decimal {
	return t.amount
}

// Date returns the date of the transaction.
func (t lbbTransaction) Date() time.Time {
	return t.date
}

// ValutaDate returns the date of the transaction.
func (t lbbTransaction) ValutaDate() time.Time {
	return t.valutaDate
}

// Currency returns a currency code for the account.
func (t lbbTransaction) Currency() string {
	return t.currency
}

func lbbParseTransaction(record []string, t *lbbTransaction) bool {
	var err error

	t.CardNumber = record[0]
	t.valutaDate, err = time.Parse("02.01.2006", record[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s: %s\n", record[1], err)
	}
	t.date, err = time.Parse("02.01.2006", record[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s: %s\n", record[2], err)
	}
	if matched, _ := regexp.MatchString("[+-] .* AMAZON(.DE)? PUNKTE", record[3]); matched {
		var sign rune
		var value int
		// New format
		n, err := fmt.Sscanf(strings.TrimSpace(record[3]), "%c %d.0 AMAZON PUNKTE", &sign, &value)
		if n != 2 {
			// Old format
			n2, err2 := fmt.Sscanf(strings.TrimSpace(record[3]), "%c %d.0 AMAZON.DE PUNKTE", &sign, &value)
			if n2 != 2 {
				fmt.Fprintf(os.Stderr, "Error parsing '%s': %s, %s\n", record[3], err, err2)
				return false
			}
		}
		t.amount = goledger.Decimal(value * 100)
		if sign == '-' {
			t.amount = -t.amount
		}
		t.currency = "A"
		t.Merchant = "AMAZON PUNKTE"
	} else {
		t.currency = "EUR"
		t.Merchant = record[3]

		(&t.amount).UnmarshalJSON([]byte(strings.Replace(record[6], ",", ".", 1)))
	}
	return true
}

// LBBParseFile parses a CSV file generated by the Landesbank Berlin for their
// Amazon credit cards.
func LBBParseFile(path string) ([]Transaction, error) {
	fr, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		fr.Close()
	}()

	var transactions []Transaction

	r := csv.NewReader(fr)
	r.Comma = ';'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if !strings.ContainsRune("0123456789", rune(record[0][0])) {
			continue
		}

		var t lbbTransaction
		if !lbbParseTransaction(record, &t) {
			continue
		}
		transactions = append(transactions, &t)

	}
	return transactions, nil
}
