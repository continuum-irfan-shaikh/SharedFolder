package kafka

import (
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka"
	. "github.com/onsi/gomega"
)

func init() {
	logger.Load(config.Config.Log)
}

func TestPush(t *testing.T) {
	RegisterTestingT(t)
	type args struct {
		msg     interface{}
		msgType string
	}
	tests := []struct {
		name      string
		args      args
		publisher func(transaction string, producerType publisher.ProducerType, cfg config.Configuration, messages ...*publisher.Message) error
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				msg: models.KafkaMessage{
					PartnerID: "1",
					UID:       "2",
					Entity: struct {
						Msg string
					}{
						Msg: `TEST`,
					},
				},
				msgType: MsgTypeTaskCreated,
			},
			publisher: func(transaction string, producerType publisher.ProducerType, cfg config.Configuration, messages ...*publisher.Message) error {
				return nil
			},
		},
		{
			name: "Failed to marshal",
			args: args{
				msg: func() {},
			},
			publisher: func(transaction string, producerType publisher.ProducerType, cfg config.Configuration, messages ...*publisher.Message) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "Failed to publish",
			args: args{
				msg: models.KafkaMessage{
					PartnerID: "1",
					UID:       "2",
					Entity: struct {
						Msg string
					}{
						Msg: `TEST`,
					},
				},
				msgType: MsgTypeTaskCreated,
			},
			publisher: func(transaction string, producerType publisher.ProducerType, cfg config.Configuration, messages ...*publisher.Message) error {
				return errors.New("fail")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		kafka.PostMessage = tt.publisher
		p := NewTasking(config.Config)
		err := p.Push(tt.args.msg, tt.args.msgType)
		if tt.wantErr {
			Ω(err).ShouldNot(BeNil(), tt.name)
		} else {
			Ω(err).Should(BeNil(), tt.name)
		}
	}
}
