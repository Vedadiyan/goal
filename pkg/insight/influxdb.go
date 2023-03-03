package insight

import (
	"io"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	http2 "github.com/influxdata/influxdb-client-go/v2/api/http"
)

type InfluxDbLogger struct {
	client   influxdb.Client
	writeAPI api.WriteAPI
}

func NewInfluxDbLogger(dsn string, authToken string, org string, bucket string, failover io.Writer) *InfluxDbLogger {
	client := influxdb.NewClientWithOptions(dsn, authToken, influxdb.DefaultOptions().SetBatchSize(100).SetFlushInterval(5000))
	writeAPI := client.WriteAPI(org, bucket)
	writeAPI.SetWriteFailedCallback(func(batch string, error http2.Error, retryAttempts uint) bool {
		if failover != nil {
			failover.Write([]byte(batch))
			failover.Write([]byte(error.Error()))
			failover.Write([]byte("\r\n"))
		}
		return false
	})
	return &InfluxDbLogger{
		client:   client,
		writeAPI: writeAPI,
	}
}

func (idl InfluxDbLogger) Log(id string, origin string, logLevel string, fields map[string]any) {
	tags := map[string]string{
		"id":     id,
		"origin": origin,
	}
	idl.writeAPI.WritePoint(influxdb.NewPoint(logLevel, tags, fields, time.Now()))
}

func (idl InfluxDbLogger) Close() {
	idl.client.Close()
}
