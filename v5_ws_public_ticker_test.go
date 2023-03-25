package bybit

import (
	"encoding/json"
	"testing"

	"github.com/hirokisan/bybit/v2/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketV5Public_Ticker(t *testing.T) {
	respBody := map[string]interface{}{
		"topic": "ticker.BTCUSDT",
		"type":  "snapshot",
		"ts":    1672324988882,
		"data": []map[string]interface{}{
			{
				"symbol":            "BTCUSDT",
				"tickDirection":     "PlusTick",
				"price24hPcnt":      "0.017103",
				"lastPrice":         "17216.00",
				"prevPrice24h":      "16926.50",
				"highPrice24h":      "17281.50",
				"lowPrice24h":       "16915.00",
				"prevPrice1h":       "17238.00",
				"markPrice":         "17217.33",
				"indexPrice":        "17227.36",
				"openInterest":      "68744.761",
				"openInterestValue": "1183601235.91",
				"turnover24h":       "1570383121.943499",
				"volume24h":         "91705.276",
				"nextFundingTime":   "1673280000000",
				"fundingRate":       "-0.000212",
				"bid1Price":         "17215.50",
				"bid1Size":          "84.489",
				"ask1Price":         "17216.00",
				"ask1Size":          "83.020",
			},
		},
	}
	bytesBody, err := json.Marshal(respBody)
	require.NoError(t, err)

	category := CategoryV5Linear

	server, teardown := testhelper.NewWebsocketServer(
		testhelper.WithWebsocketHandlerOption(V5WebsocketPublicPathFor(category), bytesBody),
	)
	defer teardown()

	wsClient := NewTestWebsocketClient().
		WithBaseURL(server.URL)

	svc, err := wsClient.V5().Public(category)
	require.NoError(t, err)

	{
		_, err := svc.SubscribeTicker(
			V5WebsocketPublicTickerParamKey{
				Symbol: SymbolV5BTCUSDT,
			},
			func(response V5WebsocketPublicTickerResponse) error {
				assert.Equal(t, respBody["topic"], response.Topic)
				return nil
			},
		)
		require.NoError(t, err)
	}

	assert.NoError(t, svc.Run())
	assert.NoError(t, svc.Ping())
	assert.NoError(t, svc.Close())
}
