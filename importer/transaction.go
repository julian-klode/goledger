/*
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
	"time"

	"github.com/julian-klode/goledger"
)

// Transaction describes a generic incoming transaction.
type Transaction interface {
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
