package weshnet

import (
	"encoding/hex"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"berty.tech/weshnet/internal/queue"
)

const messageMetricNamespace = "bty_store_message"

var (
	collectorMessageStoreQueueLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: messageMetricNamespace,
			Name:      "message_queue_length",
			Help:      "message queue length",
		}, []string{"kind", "device_pk"},
	)
	collectorsMessageStore = []prometheus.Collector{
		collectorMessageStoreQueueLength,
	}
)

var _ queue.MetricsTracer[*messageItem] = (*messageMetricsTracer)(nil)

type messageMetricsTracer struct {
	reg prometheus.Registerer
}

func newMessageMetricsTracer(reg prometheus.Registerer) (mt *messageMetricsTracer) {
	mt = &messageMetricsTracer{reg: reg}
	for _, collector := range collectorsMessageStore {
		if err := reg.Register(collector); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(fmt.Errorf("message metrics errors: %w", err))
			}

			return
		}
	}
	// reg.MustRegister(collectorsMessageStore...)
	return
}

func (s *messageMetricsTracer) ItemQueued(name string, m *messageItem) {
	collectorMessageStoreQueueLength.WithLabelValues(
		name, hex.EncodeToString(m.headers.DevicePk),
	).Inc()
}

func (s *messageMetricsTracer) ItemPop(name string, m *messageItem) {
	collectorMessageStoreQueueLength.WithLabelValues(
		name, hex.EncodeToString(m.headers.DevicePk),
	).Dec()
}
