package exporters

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/topfreegames/fluxcloud/pkg/msg"
)

type FailingFakeExporter struct {
	Sent []msg.Message
}

func (f *FailingFakeExporter) Send(_ context.Context, _ *http.Client, _ msg.Message) error {
	log.Print("Sending to exporter: ", f.Name())
	return errors.New(fmt.Sprintf("Could not Send, failed on purpose!"))
}

func (f *FailingFakeExporter) NewLine() string {
	return "\n"
}

func (f *FailingFakeExporter) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

func (f *FailingFakeExporter) Name() string {
	return "FailingFake"
}
