// Package importer maintains the transaction importer.
package importer

import (
	"time"

	"github.com/julian-klode/goledger"
)

// Transaction describes a generic incoming transaction.
type Transaction interface {
	// Time when the transaction occured.
	Date() time.Time
	ValutaDate() time.Time

	// IBAN or other ID of the local account.
	LocalAccount() string

	// Information about the remote account (name and IBAN).
	RemoteName() string
	RemoteAccount() string

	// The reference text describing the transaction. May be optional.
	ReferenceText() string

	// The amount and the currency of the transaction.
	Amount() goledger.Decimal
	Currency() string
}
