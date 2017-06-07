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
	return fmt.Sprintf("%v", t.remoteName)
}

// ReferenceText returns a description of the transaction.
func (t hbciTransaction) ReferenceText() string {
	return fmt.Sprintf("%v", t.purposes)
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
