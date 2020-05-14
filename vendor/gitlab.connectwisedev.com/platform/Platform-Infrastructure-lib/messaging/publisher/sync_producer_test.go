package publisher

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func Test_syncProducer_Publish(t *testing.T) {
	t.Run("No Messages", func(t *testing.T) {
		s := &syncProducer{cfg: NewConfig()}
		if err := s.Publish(context.Background(), "transaction"); err == nil {
			t.Errorf("syncProducer.Publish() want error got nil")
		}
	})

	t.Run("One Message", func(t *testing.T) {
		s := &syncProducer{cfg: NewConfig()}
		if err := s.Publish(context.Background(), "transaction", &Message{}); err == nil {
			t.Errorf("syncProducer.Publish() want error got nil")
		}
	})

	t.Run("Blank Valid Errors", func(t *testing.T) {
		cfg := NewConfig()
		cfg.CircuitProneErrors = []string{}
		cfg.Address = []string{"localhost"}
		s := &syncProducer{cfg: cfg}
		if err := s.Publish(context.Background(), "transaction", &Message{}); err == nil {
			t.Errorf("syncProducer.Publish() want error got nil")
		}
	})

	t.Run("Invalid Address Error", func(t *testing.T) {
		cfg := NewConfig()
		cfg.Address = []string{"localhost"}
		s := &syncProducer{cfg: cfg}
		if err := s.Publish(context.Background(), "transaction", &Message{Topic: "test"}); err == nil {
			t.Errorf("syncProducer.Publish() want error got nil")
		}
	})
}

func Test_syncProducer_publish(t *testing.T) {
	cntx, cancle := context.WithTimeout(context.Background(), time.Nanosecond)
	cancle()
	type fields struct {
		cfg          *Config
		producerType ProducerType
		producer     sarama.SyncProducer
		existing     sarama.SyncProducer
	}
	type args struct {
		ctx         context.Context
		transaction string
		messages    []*Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background()}, wantErr: true},
		{
			name: "2", args: args{ctx: context.Background()}, wantErr: true,
			fields: fields{cfg: &Config{}, producerType: BigKafkaProducer},
		},
		{
			name: "3", args: args{ctx: context.Background()}, wantErr: false,
			fields: fields{cfg: &Config{Address: []string{"localhost"}}, producerType: BigKafkaProducer},
		},
		{
			name: "4", args: args{ctx: context.Background(), messages: []*Message{&Message{}}}, wantErr: true,
			fields: fields{cfg: &Config{Address: []string{"localhost"}}, producerType: BigKafkaProducer},
		},
		{
			name: "5", args: args{ctx: cntx, messages: []*Message{&Message{}}}, wantErr: true,
			fields: fields{cfg: &Config{Address: []string{"localhost"}}, producerType: BigKafkaProducer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &syncProducer{
				cfg:          tt.fields.cfg,
				producerType: tt.fields.producerType,
				producer:     tt.fields.producer,
				existing:     tt.fields.existing,
			}
			if err := s.publish(tt.args.ctx, tt.args.transaction, tt.args.messages...); (err != nil) != tt.wantErr {
				t.Errorf("syncProducer.publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
