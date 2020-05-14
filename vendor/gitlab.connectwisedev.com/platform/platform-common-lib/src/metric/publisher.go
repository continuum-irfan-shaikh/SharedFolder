package metric

import (
	"runtime"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/udp"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
	proto "github.com/golang/protobuf/proto"
)

//Collector : Interface to be implemented by all Metric types
type Collector interface {
	proto.Message
	MetricType() string
}

// Publish : Publish Metric
var Publish = func(cfg *Config, collector ...Collector) error {
	m := make([]*Message_Metric, len(collector))
	for index, c := range collector {
		data, err := proto.Marshal(c)
		if err != nil {
			return err
		}
		m[index] = &Message_Metric{Type: c.MetricType(), Value: data}
	}

	message := &Message{Metric: m, Properties: map[string]string{}}
	return publish(message, cfg)
}

// PeriodicPublish : Periodically Publish  Metric
var PeriodicPublish = func(duration time.Duration, cfg *Config, callback func() []Collector, handler func(err error)) {
	for {
		err := Publish(cfg, callback()...)
		if err != nil {
			handler(err)
		}
		time.Sleep(duration)
	}
}

func publish(message *Message, cfg *Config) error {
	message.ProcessName = util.ProcessName()
	message.HostName = util.Hostname(util.ProcessName())
	message.TimestampUnix = time.Now().UTC().UnixNano()
	message.Namespace = cfg.GetNamespace()
	message.Address = util.LocalIPAddress()
	message.Properties["GoVersion"] = runtime.Version()
	message.Properties["Architecture"] = runtime.GOARCH
	message.Properties["OS"] = runtime.GOOS
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return udp.Send(cfg.Communication, data, nil)
}
