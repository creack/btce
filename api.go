package btce

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	key    = "ZEX6EERL-I1WGX3CX-PNK535N7-Y8KE1LL6-RI6V9EX9"
	secret = "9fd746b3f64b21e6bb269d727769889d33f5c35aac58a6085b67da9593b426fc"
)

var (
	ErrUnkownMethod = errors.New("Unkown method")
)

func sign(params string) string {
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write([]byte(params))
	return string(hex.EncodeToString(hash.Sum(nil)))
}

// {"success":1,"return":{"funds":{"usd":0,"btc":1.310908,"ltc":10.99895802,"nmc":0,"rur":0,"eur":0,"nvc":0,"trc":0,"ppc":0,"ftc":0,"xpm":0},"rights":{"info":1,"trade":0,"withdraw":0},"transaction_count":15,"open_orders":0,"server_time":1386571226}}OK
/*
{
   "success" : 1,
   "return" : {
      "rights" : {
         "info" : 1,
         "withdraw" : 0,
         "trade" : 0
      },
      "funds" : {
         "nvc" : 0,
         "nmc" : 0,
         "btc" : 1.310908,
         "xpm" : 0,
         "usd" : 0,
         "ftc" : 0,
         "ltc" : 10.99895802,
         "trc" : 0,
         "rur" : 0,
         "ppc" : 0,
         "eur" : 0
      },
      "server_time" : 1386571226,
      "open_orders" : 0,
      "transaction_count" : 15
   }
*/

var Funds = []*Fund{
	{
		Name: "Bitcoin",
		Code: "btc",
		id:   2,
	},
	{
		Name: "Litecoin",
		Code: "ltc",
		id:   3,
	},
}

type Api struct {
	Url string
}

func (a *Api) do(method string, params map[string]string) (io.ReadCloser, error) {
	var (
		now = int(time.Now().Unix())
		v   = url.Values{}
	)

	v.Add("nonce", strconv.Itoa(now))
	v.Add("method", method)

	for key, value := range params {
		v.Add(key, value)
	}

	req, err := http.NewRequest("POST", a.Url, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Sign", sign(v.Encode()))
	req.Header.Set("Key", key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// FIXME: Use url.Value instead of a map
func (a *Api) doDecode(method string, data interface{}, params map[string]string) error {
	resp, err := a.do(method, params)
	if err != nil {
		return err
	}
	defer resp.Close()

	if err := json.NewDecoder(resp).Decode(data); err != nil {
		return err
	}

	var errStr string

	switch t := data.(type) {
	case *transHistoryResponse:
		errStr = t.Error
	case *infoResponse:
		errStr = t.Error
	default:
		return ErrUnkownMethod
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

func (a *Api) GetInfo() (*Info, error) {
	var jsonResponse infoResponse

	if err := a.doDecode("getInfo", &jsonResponse, nil); err != nil {
		return nil, err
	}
	// Once we retrieve the raw map[Currency] = amount, we populate
	// Funds objects
	for code, value := range jsonResponse.Return.FundsJ {
		for _, f := range Funds {
			if code == f.Code {
				jsonResponse.Return.Funds[f] = value
			}
		}

	}
	// Convert the timestamp in time.Time
	jsonResponse.Return.ServerTime = time.Unix(jsonResponse.Return.ServerTimestamp, 0)

	return jsonResponse.Return, nil
}

// FIXME: Use reflect with json like tags
func (a *Api) encodeOptions(options *Options) map[string]string {
	params := make(map[string]string)

	if options.Count != 0 {
		params["count"] = strconv.Itoa(options.Count)
	}
	if options.Order {
		params["order"] = "ASC"
	}
	if options.FromId != 0 {
		params["from_id"] = strconv.Itoa(options.FromId)
	}
	if options.EndId != 0 {
		params["end_id"] = strconv.Itoa(options.EndId)
	}
	if options.Since != nil {
		params["since"] = strconv.Itoa(int(options.Since.Unix()))
	}
	if options.End != nil {
		params["end"] = strconv.Itoa(int(options.End.Unix()))
	}
	return params
}

func (a *Api) TransHistory(options *Options) (map[string]*TransHistory, error) {
	var jsonResponse transHistoryResponse

	params := a.encodeOptions(options)
	if err := a.doDecode("TransHistory", &jsonResponse, params); err != nil {
		return nil, err
	}
	// For each return row, Populate TransacHistory objects
	// including the transaction id
	for k, v := range jsonResponse.Return {
		id, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		v.Id = id
		// Also convert timestamp in time.Time
		v.Time = time.Unix(v.Timestamp, 0)
	}
	return jsonResponse.Return, nil
}

func (a *Api) TradeHistory(options *Options) (map[string]*TradeHistory, error) {
	var jsonResponse transHistoryResponse

	params := a.encodeOptions(options)
	if err := a.doDecode("TradeHistory", &jsonResponse, params); err != nil {
		return nil, err
	}
	// For each return row, Populate TransacHistory objects
	// including the transaction id
	for k, v := range jsonResponse.Return {
		id, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		v.Id = id
		// Also convert timestamp in time.Time
		v.Time = time.Unix(v.Timestamp, 0)
	}
	return jsonResponse.Return, nil
}

/*
	$req['method'] = $method;
        $mt = explode(' ', microtime());
        $req['nonce'] = $mt[1];

        // generate the POST data string
        $post_data = http_build_query($req, '', '&');

        $sign = hash_hmac('sha512', $post_data, $secret);
*/
