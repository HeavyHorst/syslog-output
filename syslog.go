package syslog

import (
	"encoding/json"
	"os"
	"time"

	"github.com/micro/go-platform/log"
	"github.com/papertrail/remote_syslog2/syslog"
)

type output struct {
	opts   log.OutputOptions
	logger *syslog.Logger
	err    error
}

type msg struct {
	Level   string     `json:"level"`
	Fields  log.Fields `json:"fields"`
	Message string     `json:"message"`
}

func (o *output) Send(e *log.Event) error {
	if o.err != nil {
		return o.err
	}

	dat, err := json.Marshal(msg{
		Level:   log.Levels[e.Level],
		Fields:  e.Fields,
		Message: e.Message,
	})
	if err != nil {
		return err
	}

	f, _ := syslog.Facility("user")
	packet := syslog.Packet{
		Hostname: o.logger.ClientHostname,
		Tag:      log.Levels[e.Level],
		Time:     time.Unix(0, e.Timestamp),
		Message:  string(dat),
		Facility: f,
	}

	switch e.Level {
	case log.InfoLevel:
		packet.Severity = syslog.SevInfo
	case log.DebugLevel:
		packet.Severity = syslog.SevDebug
	case log.ErrorLevel:
		packet.Severity = syslog.SevErr
	case log.WarnLevel:
		packet.Severity = syslog.SevWarning
	case log.FatalLevel:
		packet.Severity = syslog.SevCrit
	}

	o.logger.Packets <- packet

	return nil
}

func (o *output) Flush() error {
	return nil
}
func (o *output) Close() error {
	return nil
}

func (o *output) String() string {
	return "syslog"
}

func NewOutput(protocol string, opts ...log.OutputOption) log.Output {
	var options log.OutputOptions
	out := &output{
		err: nil,
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Name) == 0 {
		options.Name = "127.0.0.1:514"
	}

	hostname, err := os.Hostname()
	if err != nil {
		out.err = err
	}

	logger, err := syslog.Dial(options.Name, protocol, options.Name, nil, 5*time.Second, 5*time.Second, 0)
	if err != nil {
		out.err = err
	}

	logger.ClientHostname = hostname
	out.opts = options
	out.logger = logger

	return out
}
