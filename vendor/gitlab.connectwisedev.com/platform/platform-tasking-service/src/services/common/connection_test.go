package common

import (
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetOutboundConnectionStatus(t *testing.T) {
	got := GetOutboundConnectionStatus()

	var test = []OutboundConnectionStatus{
		{
			TimeStampUTC:   time.Now(),
			Type:           dbMessageType,
			Name:           fmt.Sprintf("%s-%s", config.Config.Version.ServiceName, dbNameSuffix),
			ConnectionType: "Cassandra",
			ConnectionURLs: strings.SplitN(config.Config.CassandraURL, ",", -1),
		},
		{
			TimeStampUTC:   time.Now(),
			Type:           dbMessageType,
			Name:           fmt.Sprintf("%s-%s", config.Config.Version.ServiceName, dbNameSuffix),
			ConnectionType: "Kafka",
			ConnectionURLs: strings.SplitN(config.Config.KafkaBrokers, ",", -1),
		},
	}

	test[0].ConnectionStatus = connMethods["Cassandra"](test[0]).ConnectionStatus
	test[1].ConnectionStatus = connMethods["Kafka"](test[1]).ConnectionStatus

	for i := range got {
		if got[i].TimeStampUTC.After(test[i].TimeStampUTC) {
			t.Fatalf("Got %v, want %v", got[i].TimeStampUTC, test[i].TimeStampUTC)
		}
		if !reflect.DeepEqual(got[i].ConnectionURLs, test[i].ConnectionURLs) {
			t.Fatalf("Got %v, want %v", got[i].ConnectionURLs, test[i].ConnectionURLs)
		}
		if got[i].Type != test[i].Type {
			t.Fatalf("Got %v, want %v", got[i].Type, test[i].Type)
		}
		if got[i].Name != test[i].Name {
			t.Fatalf("Got %v, want %v", got[i].Name, test[i].Name)
		}
		if got[i].ConnectionStatus != test[i].ConnectionStatus {
			t.Fatalf("Got %v, want %v", got[i].ConnectionStatus, test[i].ConnectionStatus)
		}
	}
}
