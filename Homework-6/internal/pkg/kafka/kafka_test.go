package kafka

import (
	"homework-3/internal/pkg/kafka/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWriteToKafka(t *testing.T) {
	event := Event{
		Timestamp: time.Now(),
		Method:    "GET",
		Request:   "/api/v1/resource",
	}
	topic := "test-topic"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProducer := mocks.NewMockSyncProducer(ctrl)

	mockProducer.EXPECT().
		SendMessage(gomock.Any()).
		Return(int32(0), int64(0), nil).
		Times(1)

	err := WriteToKafka(event, mockProducer, topic)
	assert.NoError(t, err, "Expected no error when writing to Kafka")
}
