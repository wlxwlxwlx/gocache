package metrics

import "github.com/wlxwlxwlx/gocache/v2/codec"

// MetricsInterface represents the metrics interface for all available providers
type MetricsInterface interface {
	RecordFromCodec(codec codec.CodecInterface)
}
