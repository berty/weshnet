package weshnet

import (
	"encoding/hex"

	"berty.tech/weshnet/internal/queue"
	"github.com/prometheus/client_golang/prometheus"
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
	reg       prometheus.Registerer
	namespace string
}

func NewMessageMetricsTracer(reg prometheus.Registerer) *messageMetricsTracer {
	reg.MustRegister(collectorsMessageStore...)
	return &messageMetricsTracer{
		reg: reg,
	}
}

func (s *messageMetricsTracer) ItemQueued(name string, m *messageItem) {
	collectorMessageStoreQueueLength.WithLabelValues(
		name, string(hex.EncodeToString(m.headers.DevicePK)),
	).Inc()
}

func (s *messageMetricsTracer) ItemPop(name string, m *messageItem) {
	collectorMessageStoreQueueLength.WithLabelValues(
		name, string(hex.EncodeToString(m.headers.DevicePK)),
	).Dec()
}
