package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/formatters"
	test_utils "github.com/topfreegames/fluxcloud/pkg/utils/test"
)

func TestHandleV6(t *testing.T) {
	fakeExporter := &exporters.FakeExporter{}
	fakeConfig := config.NewFakeConfig()
	fakeConfig.Set("github_url", "https://github.com")

	formatter, _ := formatters.NewDefaultFormatter(fakeConfig)

	apiConfig := APIConfig{
		Server:    http.NewServeMux(),
		Exporter:  []exporters.Exporter{fakeExporter},
		Formatter: formatter,
	}

	assert.NoError(t, HandleV6(apiConfig))

	event := test_utils.NewFluxSyncEvent()
	data, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", "http://127.0.0.1:3030/v6/events", bytes.NewBuffer(data))
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	apiConfig.Server.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)

	formatted := formatter.FormatEvent(event, fakeExporter)
	assert.Equal(t, formatted.Title, fakeExporter.Sent[0].Title, formatted.Title)
	assert.Equal(t, formatted.Body, fakeExporter.Sent[0].Body, formatted.Body)
}

func TestHandleV6_multiple_exporters_fail(t *testing.T) {
	fakeConfig := config.NewFakeConfig()
	fakeConfig.Set("github_url", "https://github.com")

	formatter, _ := formatters.NewDefaultFormatter(fakeConfig)

	apiConfig := APIConfig{
		Server:    http.NewServeMux(),
		Exporter:  []exporters.Exporter{&exporters.FailingFakeExporter{}, &exporters.FakeExporter{}},
		Formatter: formatter,
	}

	assert.NoError(t, HandleV6(apiConfig))

	event := test_utils.NewFluxSyncEvent()
	data, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", "http://127.0.0.1:3030/v6/events", bytes.NewBuffer(data))
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	apiConfig.Server.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 500, resp.StatusCode)
}
