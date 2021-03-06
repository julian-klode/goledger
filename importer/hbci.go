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
	"strconv"
	"time"

	"github.com/julian-klode/goledger"
)

type hbciTransaction struct {
	localAccountNumber  string
	remoteAccountNumber string
	valueValue          goledger.Decimal
	valueCurrency       string
	remoteName          []string
	purposes            []string
	date                time.Time
	valutaDate          time.Time
}

func (t hbciTransaction) ID() string {
	return hashTransaction(t)
}

func (t hbciTransaction) Category() Category {
	return CategoryMisc
}

// LocalAccount returns an ID of the local account.
func (t hbciTransaction) LocalAccount() string {
	return t.localAccountNumber
}

// RemoteAccount returns an ID of the remote account (IBAN).
func (t hbciTransaction) RemoteAccount() string {
	return t.remoteAccountNumber
}

// RemoteName returns a name of the other account.
func (t hbciTransaction) RemoteName() string {
	result := ""
	for _, s := range t.remoteName {
		result += s
	}
	return result
}

// ReferenceText returns a description of the transaction.
func (t hbciTransaction) ReferenceText() string {
	result := ""
	for _, s := range t.purposes {
		result += s
	}
	return result
}

// Amount returns the amount of the transaction.
func (t hbciTransaction) Amount() goledger.Decimal {
	return t.valueValue
}

// Date returns the date of the transaction.
func (t hbciTransaction) Date() time.Time {
	return t.date
}

// ValutaDate returns the date of the transaction.
func (t hbciTransaction) ValutaDate() time.Time {
	return t.valutaDate
}

// Currency returns a currency code for the account.
func (t hbciTransaction) Currency() string {
	return t.valueCurrency
}

// RemoteNames is like RemoteName() but exposes the slice
func (t hbciTransaction) RemoteNames() []string {
	return t.remoteName
}

// Purposes is like Purpose() but exposes the slice
func (t hbciTransaction) Purposes() []string {
	return t.purposes
}

// HBCIParseFile parses a CSV file generated by acqbanking-cli listtrans.
func HBCIParseFile(path string) ([]Transaction, error) {
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

	record, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	columns := make(map[string]int)
	for i, s := range record {
		columns[s] = i
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recordStr := ""

		for _, c := range record {
			recordStr += ";" + c
		}
		var t hbciTransaction

		t.localAccountNumber = record[columns["localAccountNumber"]]
		t.remoteAccountNumber = record[columns["remoteAccountNumber"]]
		(&t.valueValue).UnmarshalJSON([]byte(record[columns["value_value"]]))
		t.valueCurrency = record[columns["value_currency"]]
		t.remoteName = append(t.remoteName, record[columns["remoteName"]])
		for i := 1; columns["remoteName"+strconv.Itoa(i)] != 0; i++ {
			if record[columns["remoteName"+strconv.Itoa(i)]] != "" {
				t.remoteName = append(t.remoteName, record[columns["remoteName"+strconv.Itoa(i)]])
			}
		}
		t.purposes = append(t.purposes, record[columns["purpose"]])
		for i := 1; columns["purpose"+strconv.Itoa(i)] != 0; i++ {
			if record[columns["purpose"+strconv.Itoa(i)]] != "" {
				t.purposes = append(t.purposes, record[columns["purpose"+strconv.Itoa(i)]])
			}
		}

		t.date, err = time.Parse("2006/01/02", record[columns["date"]])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %s: %s\n", record[columns["date"]], err)
		}
		t.valutaDate, err = time.Parse("2006/01/02", record[columns["valutadate"]])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %s: %s\n", record[columns["date"]], err)
		}
		transactions = append(transactions, &t)

	}
	return transactions, nil
}
