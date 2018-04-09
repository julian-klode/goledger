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
	"encoding/json"
	"os"
	"time"

	"github.com/julian-klode/goledger"
)

// n26Transaction is a transaction in an N26 account.
type n26Transaction struct {
	ID                   string           `json:"id"`
	UserID               string           `json:"userId"`
	Type                 string           `json:"type"`
	Amount               goledger.Decimal `json:"amount"`
	CurrencyCode         string           `json:"currencyCode"`
	MerchantCity         string           `json:"merchantCity,omitempty"`
	VisibleTS            int64            `json:"visibleTS"`
	Mcc                  int              `json:"mcc,omitempty"`
	MccGroup             int              `json:"mccGroup,omitempty"`
	MerchantName         string           `json:"merchantName,omitempty"`
	Recurring            bool             `json:"recurring"`
	AccountID            string           `json:"accountId"`
	Category             string           `json:"category"`
	CardID               string           `json:"cardId,omitempty"`
	UserCertified        int64            `json:"userCertified,omitempty"`
	Pending              bool             `json:"pending"`
	TransactionNature    string           `json:"transactionNature"`
	CreatedTS            int64            `json:"createdTS"`
	MerchantCountry      int              `json:"merchantCountry,omitempty"`
	SmartLinkID          string           `json:"smartLinkId"`
	LinkID               string           `json:"linkId"`
	Confirmed            int64            `json:"confirmed,omitempty"`
	PartnerBic           string           `json:"partnerBic,omitempty"`
	PartnerAccountIsSepa bool             `json:"partnerAccountIsSepa,omitempty"`
	PartnerName          string           `json:"partnerName,omitempty"`
	PartnerIban          string           `json:"partnerIban,omitempty"`
	ReferenceText        string           `json:"referenceText,omitempty"`
	UserAccepted         int64            `json:"userAccepted,omitempty"`
	PartnerBcn           string           `json:"partnerBcn,omitempty"`
	PartnerAccountBan    string           `json:"partnerAccountBan,omitempty"`
	SmartContactID       string           `json:"smartContactId,omitempty"`
	OriginalAmount       goledger.Decimal `json:"originalAmount,omitempty"`
	OriginalCurrency     string           `json:"originalCurrency,omitempty"`
	ExchangeRate         float64          `json:"exchangeRate,omitempty"`
	MerchantID           string           `json:"merchantId,omitempty"`
	TransactionTerminal  string           `json:"transactionTerminal,omitempty"`
	PartnerBankName      string           `json:"partnerBankName,omitempty"`
	BankTransferTypeText string           `json:"bankTransferTypeText,omitempty"`
}

type n26Transaction2 struct{ d *n26Transaction }

func (t n26Transaction2) ID() string {
	return t.d.ID
}

var categories = map[string]Category{
	"micro-v2-atm":                   CategoryATM,
	"micro-v2-business":			  CategoryBusiness,
	"micro-v2-food-groceries":        CategoryFoodGroceries,
	"micro-v2-income":                CategoryIncome,
	"micro-v2-leisure-entertainment": CategoryLeisureEntertainment,
	"micro-v2-miscellaneous":         CategoryMisc,
	"micro-v2-savings-investments":   CategorySavingsInvestments,
	"micro-v2-shopping":              CategoryShopping,
	"micro-v2-transport-car":         CategoryTransportCar,
	"micro-v2-bars-restaurants":	  CategoryBarsRestaurants,
	"micro-v2-travel-holidays": 	  CategoryTravelHolidays,
}

func (t n26Transaction2) Category() Category {
	return categories[t.d.Category]
}

// LocalAccount returns an ID of the local account.
func (t n26Transaction2) LocalAccount() string {
	return t.d.AccountID
}

// RemoteAccount returns an ID of the remote account (IBAN).
func (t n26Transaction2) RemoteAccount() string {
	return t.d.PartnerIban
}

// RemoteName returns a name of the other account.
func (t n26Transaction2) RemoteName() string {
	switch {
	case t.d.PartnerName != "":
		return t.d.PartnerName
	case t.d.MerchantName != "":
		return t.d.MerchantName
	default:
		return ""
	}
}

// ReferenceText returns a description of the transaction.
func (t n26Transaction2) ReferenceText() string {
	return t.d.ReferenceText
}

// Amount returns the amount of the transaction.
func (t n26Transaction2) Amount() goledger.Decimal {
	return t.d.Amount
}

// Date returns the date of the transaction.
func (t n26Transaction2) Date() time.Time {
	return time.Unix(t.d.VisibleTS/1000, 0)
}

// ValutaDate returns the date of the transaction.
func (t n26Transaction2) ValutaDate() time.Time {
	return time.Unix(t.d.CreatedTS/1000, 0)
}

// Currency returns a currency code for the account.
func (t n26Transaction2) Currency() string {
	return t.d.CurrencyCode
}

// N26ParseFile parses a N26 JSON file into a slice of transactions.
func N26ParseFile(path string) ([]Transaction, error) {
	var transactions []n26Transaction
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		r.Close()
	}()
	err = json.NewDecoder(r).Decode(&transactions)
	if err != nil {
		return nil, err
	}
	var results []Transaction
	for l := len(transactions); l > 0; l-- {
		t := Transaction(n26Transaction2{&transactions[l-1]})
		// TODO: Filter out transactions that are not completed yet.
		results = append(results, t)
	}
	return results, nil
}
