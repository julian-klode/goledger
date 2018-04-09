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

// Package importer maintains the transaction importer.
package importer

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/julian-klode/goledger"
)

// Category represents categories of a transaction. Only supported by the
// N26Account.
type Category int

// Possible categories
const (
	CategoryMisc Category = iota
	CategoryATM
	CategoryBusiness
	CategoryFoodGroceries
	CategoryIncome
	CategoryLeisureEntertainment
	CategorySavingsInvestments
	CategoryShopping
	CategoryTransportCar
	CategoryBarsRestaurants
)

// Transaction describes a generic incoming transaction.
type Transaction interface {
	// An identifier describing the description, to filter out duplicates.
	//
	// If the bank does not provide identifiers, use hashTransaction() when
	// implementing a new transaction parser.
	ID() string
	// Category of the transaction
	Category() Category
	// Time when the transaction occured. For the difference between date
	// and valuta date search the internet, I can't explain it.
	Date() time.Time
	ValutaDate() time.Time

	// IBAN or other ID of the local account. For credit card transactions, the
	// local account describes some sort of card id.
	//
	// For LBB, this is the master card id or a related number indicating
	// incoming transactions.
	//
	// For N26, this is a UUID or something like that.
	LocalAccount() string

	// Information about the remote account (name and IBAN), and the content
	// of the transaction. For a SEPA transaction, all of these will likely
	// be set, for a credit card transaction, RemoteName() might be the only one.
	//
	// When matching keywords in transactions for categorizing, both remote
	// name and the reference text should be checked; in a canonical case,
	// and with and without spaces removed, as sometimes formatting fails in
	// a transaction and spaces go missing.
	RemoteName() string
	RemoteAccount() string
	ReferenceText() string

	// The amount and the currency of the transaction. The amount has a two
	// digit precision, and the currency is expected to be a multi-letter
	// currency code, like EUR.
	//
	// The LBB provider understands a code 'A' which means Amazon credits.
	Amount() goledger.Decimal
	Currency() string
}

// MultilineTransaction is a transaction where names and purposes can have
// multiple lines.
type MultilineTransaction interface {
	Transaction

	RemoteNames() []string
	Purposes() []string
}

// hashTransaction is a base implementation for Transaction.ID().
// It just hashes all values using SHA1.
func hashTransaction(t Transaction) string {
	hash := sha1.New()
	_, err := hash.Write([]byte(t.Date().String()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(t.ValutaDate().String()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(t.LocalAccount()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(t.RemoteName()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(t.RemoteAccount()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(t.ReferenceText()))
	if err != nil {
		panic(err)
	}
	_, err = hash.Write([]byte(string(t.Amount())))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
