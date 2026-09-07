package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Registry collects operational metrics for Prometheus scraping.
type Registry struct {
	startTime              time.Time
	mu                     sync.RWMutex
	studiesReceived        int64
	studiesDelivered       int64
	studiesFailed          int64
	dicomCStoreReceived    int64
	networkRetries         int64
	purgedDerivativesCount int64
	storageBytesUsed       int64
	storageDiskFreeBytes   int64
	queuePendingGauge      int64
}

var DefaultRegistry = NewRegistry()

// NewRegistry initializes an operational metrics collector.
func NewRegistry() *Registry {
	return &Registry{
		startTime:            time.Now(),
		storageDiskFreeBytes: 120 * 1024 * 1024 * 1024, // 120 GB default free
	}
}

func (r *Registry) IncStudiesReceived() {
	atomic.AddInt64(&r.studiesReceived, 1)
}

func (r *Registry) IncStudiesDelivered() {
	atomic.AddInt64(&r.studiesDelivered, 1)
}

func (r *Registry) IncStudiesFailed() {
	atomic.AddInt64(&r.studiesFailed, 1)
}

func (r *Registry) IncDicomCStoreReceived() {
	atomic.AddInt64(&r.dicomCStoreReceived, 1)
}

func (r *Registry) IncNetworkRetries() {
	atomic.AddInt64(&r.networkRetries, 1)
}

func (r *Registry) AddPurgedDerivatives(count int64) {
	atomic.AddInt64(&r.purgedDerivativesCount, count)
}

func (r *Registry) SetQueuePending(val int64) {
	atomic.StoreInt64(&r.queuePendingGauge, val)
}

func (r *Registry) SetStorageBytes(used, free int64) {
	atomic.StoreInt64(&r.storageBytesUsed, used)
	if free > 0 {
		atomic.StoreInt64(&r.storageDiskFreeBytes, free)
	}
}

// Handler returns an HTTP handler serving Prometheus plaintext format.
func (r *Registry) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uptime := time.Since(r.startTime).Seconds()

		var b strings.Builder

		b.WriteString("# HELP rx_dispatch_uptime_seconds Total runtime uptime in seconds.\n")
		b.WriteString("# TYPE rx_dispatch_uptime_seconds gauge\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_uptime_seconds %.1f\n\n", uptime))

		b.WriteString("# HELP rx_dispatch_studies_total Total radiological studies processed by status.\n")
		b.WriteString("# TYPE rx_dispatch_studies_total counter\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_studies_total{status=\"RECEIVED\"} %d\n", atomic.LoadInt64(&r.studiesReceived)))
		b.WriteString(fmt.Sprintf("rx_dispatch_studies_total{status=\"DELIVERED\"} %d\n", atomic.LoadInt64(&r.studiesDelivered)))
		b.WriteString(fmt.Sprintf("rx_dispatch_studies_total{status=\"FAILED\"} %d\n\n", atomic.LoadInt64(&r.studiesFailed)))

		b.WriteString("# HELP rx_dispatch_dicom_instances_received_total Native DICOM C-STORE instances ingested.\n")
		b.WriteString("# TYPE rx_dispatch_dicom_instances_received_total counter\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_dicom_instances_received_total %d\n\n", atomic.LoadInt64(&r.dicomCStoreReceived)))

		b.WriteString("# HELP rx_dispatch_delivery_queue_pending_items Studies waiting in store-and-forward queue.\n")
		b.WriteString("# TYPE rx_dispatch_delivery_queue_pending_items gauge\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_delivery_queue_pending_items %d\n\n", atomic.LoadInt64(&r.queuePendingGauge)))

		b.WriteString("# HELP rx_dispatch_delivery_network_retries_total Resilient network retries executed.\n")
		b.WriteString("# TYPE rx_dispatch_delivery_network_retries_total counter\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_delivery_network_retries_total %d\n\n", atomic.LoadInt64(&r.networkRetries)))

		b.WriteString("# HELP rx_dispatch_storage_bytes_total Bytes stored in artifact repository.\n")
		b.WriteString("# TYPE rx_dispatch_storage_bytes_total gauge\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_storage_bytes_total %d\n\n", atomic.LoadInt64(&r.storageBytesUsed)))

		b.WriteString("# HELP rx_dispatch_storage_disk_free_bytes Available disk bytes on local storage volume.\n")
		b.WriteString("# TYPE rx_dispatch_storage_disk_free_bytes gauge\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_storage_disk_free_bytes %d\n\n", atomic.LoadInt64(&r.storageDiskFreeBytes)))

		b.WriteString("# HELP rx_dispatch_storage_purged_derivatives_total Derivatives purged by retention policy.\n")
		b.WriteString("# TYPE rx_dispatch_storage_purged_derivatives_total counter\n")
		b.WriteString(fmt.Sprintf("rx_dispatch_storage_purged_derivatives_total %d\n", atomic.LoadInt64(&r.purgedDerivativesCount)))

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(b.String()))
	}
}
