package tests

import (
	"bytes"
	"homework-3/internal/pkg/kafka"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteToKafka(t *testing.T) {
	kafka.InitKafka()
	t.Run("smoke test", func(t *testing.T) {
		event := kafka.Event{
			Timestamp: time.Now(),
			Method:    "test",
			Request:   "test",
		}
		err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic)

		require.NoError(t, err)
	})
}

func TestReadFromKafka(t *testing.T) {
	kafka.InitKafka()
	event := kafka.Event{
		Timestamp: time.Date(2024, time.April, 2, 16, 54, 10, 0, time.UTC),
		Method:    "test",
		Request:   "test",
	}
	t.Run("smoke test", func(t *testing.T) {
		go kafka.ReadFromKafka(kafka.KafPrCo.Consumer, kafka.Topic)
		time.Sleep(2 * time.Second)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "Received message: 2024-04-02T16:54:10Z test: test\n"
		assert.Equal(t, expected, buf.String())
	})
}
