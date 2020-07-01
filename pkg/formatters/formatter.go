package formatters

import (
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/msg"
)

// Formats a flux event for an exporter
type Formatter interface {
	FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message
}
