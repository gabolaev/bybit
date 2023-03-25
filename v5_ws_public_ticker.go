package bybit

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

// SubscribeTicker :
func (s *V5WebsocketPublicService) SubscribeTicker(
	key V5WebsocketPublicTickerParamKey,
	f func(V5WebsocketPublicTickerResponse) error,
) (func() error, error) {
	if err := s.addParamTickerFunc(key, f); err != nil {
		return nil, err
	}
	param := struct {
		Op   string        `json:"op"`
		Args []interface{} `json:"args"`
	}{
		Op:   "subscribe",
		Args: []interface{}{key.Topic()},
	}
	buf, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	if err := s.writeMessage(websocket.TextMessage, buf); err != nil {
		return nil, err
	}
	return func() error {
		param := struct {
			Op   string        `json:"op"`
			Args []interface{} `json:"args"`
		}{
			Op:   "unsubscribe",
			Args: []interface{}{key.Topic()},
		}
		buf, err := json.Marshal(param)
		if err != nil {
			return err
		}
		if err := s.writeMessage(websocket.TextMessage, []byte(buf)); err != nil {
			return err
		}
		s.removeParamTickerFunc(key)
		return nil
	}, nil
}

// V5WebsocketPublicTickerParamKey :
type V5WebsocketPublicTickerParamKey struct {
	Symbol SymbolV5
}

// Topic :
func (k *V5WebsocketPublicTickerParamKey) Topic() string {
	return fmt.Sprintf("%s.%s", V5WebsocketPublicTopicTicker, k.Symbol)
}

// V5WebsocketPublicTickerResponse :
type V5WebsocketPublicTickerResponse struct {
	Topic     string                      `json:"topic"`
	Type      string                      `json:"type"`
	TimeStamp int64                       `json:"ts"`
	Data      V5WebsocketPublicTickerData `json:"data"`
}

// V5WebsocketPublicTickerData :
type V5WebsocketPublicTickerData struct {
	Symbol            string `json:"symbol"`
	TickDirection     string `json:"tickDirection"`
	Price24hPcnt      string `json:"price24hPcnt"`
	LastPrice         string `json:"lastPrice"`
	PrevPrice24h      string `json:"prevPrice24h"`
	HighPrice24h      string `json:"highPrice24h"`
	LowPrice24h       string `json:"lowPrice24h"`
	PrevPrice1h       string `json:"prevPrice1h"`
	MarkPrice         string `json:"markPrice"`
	IndexPrice        string `json:"indexPrice"`
	OpenInterest      string `json:"openInterest"`
	OpenInterestValue string `json:"openInterestValue"`
	Turnover24h       string `json:"turnover24h"`
	Volume24h         string `json:"volume24h"`
	NextFundingTime   string `json:"nextFundingTime"`
	FundingRate       string `json:"fundingRate"`
	Bid1Price         string `json:"bid1Price"`
	Bid1Size          string `json:"bid1Size"`
	Ask1Price         string `json:"ask1Price"`
	Ask1Size          string `json:"ask1Size"`
}

// Key :
func (r *V5WebsocketPublicTickerResponse) Key() V5WebsocketPublicTickerParamKey {
	topic := r.Topic
	arr := strings.Split(topic, ".")
	if arr[0] != V5WebsocketPublicTopicTicker || len(arr) != 2 {
		return V5WebsocketPublicTickerParamKey{}
	}

	return V5WebsocketPublicTickerParamKey{
		Symbol: SymbolV5(arr[1]),
	}
}

// addParamTickerFunc :
func (s *V5WebsocketPublicService) addParamTickerFunc(key V5WebsocketPublicTickerParamKey, f func(V5WebsocketPublicTickerResponse) error) error {
	if _, exist := s.paramTickerMap[key]; exist {
		return errors.New("already registered for this key")
	}
	s.paramTickerMap[key] = f
	return nil
}

// removeParamTradeFunc :
func (s *V5WebsocketPublicService) removeParamTickerFunc(key V5WebsocketPublicTickerParamKey) {
	delete(s.paramTickerMap, key)
}

// retrievePositionFunc :
func (s *V5WebsocketPublicService) retrieveTickerFunc(key V5WebsocketPublicTickerParamKey) (func(V5WebsocketPublicTickerResponse) error, error) {
	f, exist := s.paramTickerMap[key]
	if !exist {
		return nil, errors.New("func not found")
	}
	return f, nil
}
