package btce

import (
	"time"
)

type TransHistory struct {
	Id          int       `json:-`
	Type        int       `json:"type"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"desc"`
	Status      int       `json:"status"`
	Timestamp   int64     `json:"timestamp"`
	Time        time.Time `json:-`
}

type TradeHistory struct {
	Id        int       `json:-`
	Pair      string    `json:"pair"`
	Type      string    `json:"type"`
	Amount    int       `json:"amount"`
	Rate      float64   `json:"rate"`
	Timestamp int64     `json:"timestamp"`
	Time      time.Time `json:-`
}

type Fund struct {
	Name string
	Code string
	qty  float64
	id   int
}

/*
from NoThe ID of the transaction to start displaying withnumerical0
count NoThe number of transactions for displayingnumerical1,000
from_id NoThe ID of the transaction to start displaying withnumerical0
end_id NoThe ID of the transaction to finish displaying withnumerical∞
order NosortingASC or DESCDESC
since NoWhen to start displaying?UNIX time0
end NoWhen to finish displaying?UNIX time∞
*/

type Options struct {
	Since  *time.Time
	End    *time.Time
	Count  int
	FromId int
	EndId  int
	Order  bool // false: Descending, true: Ascending
}

// GetInfo types
type (
	InfoRights struct {
		Info     int `json:"info"`
		Withdraw int `json:"withdraw"`
		Trade    int `json:"trade"`
	}

	Info struct {
		Rights           InfoRights         `json:"rights"`
		FundsJ           map[string]float64 `json:"funds"`
		Funds            map[*Fund]float64  `json:-`
		ServerTime       time.Time          `json:-`
		ServerTimestamp  int64              `json:"server_time"`
		OpenOrders       int                `json:"open_orders"`
		TransactionCount int                `json:"transaction_count"`
	}
)

// Internal response types
type (
	response struct {
		Success int    `json:"success"`
		Error   string `json:"error"`
	}

	transHistoryResponse struct {
		response
		Return map[string]*TransHistory `json:"return"`
	}

	infoResponse struct {
		response
		Return *Info `json: return`
	}
)
