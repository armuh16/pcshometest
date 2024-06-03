package enum

import (
	"fmt"

	"pcstakehometest/static"
)

type TransactionStatusType int

const (
	TransactionStatusTypePending TransactionStatusType = 1
	TransactionStatusTypeAccept  TransactionStatusType = 2
)

func (t TransactionStatusType) String() string {
	switch t {
	case TransactionStatusTypePending:
		return "Pending"
	case TransactionStatusTypeAccept:
		return "Accept"
	default:
		return "Unknown"
	}
}

func (t TransactionStatusType) IsValid() error {
	switch t {
	case TransactionStatusTypePending, TransactionStatusTypeAccept:
		return nil
	}
	return fmt.Errorf(static.DataNotFound, "Tipe Status Transaksi")
}
