package payment

import (
	"strconv"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type IMidtrans interface {
	GetPaymentURL(transaction models.Transaction, user models.User) (string, error)
}

type MidtransConfig struct {
	ServerKey string
	EnvType   midtrans.EnvironmentType
}

type Midtrans struct {
	Config *MidtransConfig
}

func NewMidtrans(config *MidtransConfig) IMidtrans {
	return &Midtrans{Config: config}
}

func (m *Midtrans) GetPaymentURL(transaction models.Transaction, user models.User) (string, error) {
	midtransClient := snap.Client{}
	midtransClient.New(m.Config.ServerKey, m.Config.EnvType)

	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapURL, err := midtransClient.CreateTransactionUrl(request)
	if err != nil {
		return "", err
	}

	return snapURL, nil
}
