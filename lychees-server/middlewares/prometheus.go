package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"lychees-server/logs"
	"strconv"
	"time"
)

const (
	PrometheusUrl = "http://localhost:9091"
	PrometheusJob = "lychees_prometheus"

	PrometheusNamespace    = "lychees_data"
	EndpointsDataSubsystem = "endpoints"
)

var (
	pusher *push.Pusher

	endpointsLantencyMonitor = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: EndpointsDataSubsystem,
			Name:      "lantency_statistic",
			Help:      "统计耗时数据",
			Buckets:   []float64{1, 5, 10, 20, 50, 100, 500, 1000, 5000, 10000},
		}, []string{EndpointsDataSubsystem},
	)
)

func init() {
	pusher = push.New(PrometheusUrl, PrometheusJob)
	prometheus.MustRegister(
		endpointsLantencyMonitor,
	)
	pusher.Collector(endpointsLantencyMonitor)

	go func() {
		// 每15秒上报一次数据
		for range time.Tick(15 * time.Second) {
			if err := pusher.Add(); err != nil {
				logs.Logger.Errorf("push fail %s", err)
				continue
			}
			logs.Logger.Info("push data to prometheus")
		}
	}()
}
func HandleEndpointLantency() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.Request.URL.Path
		//logs.Logger.Infof("%s被访问了", endpoint)
		start := time.Now()

		defer func(c *gin.Context) {
			lantency := time.Since(start)
			lantencyStr := fmt.Sprintf("%0.3d", lantency.Nanoseconds()/1e6) // 记录ms数据，为小数点后3位
			lantencyFloat64, err := strconv.ParseFloat(lantencyStr, 64)     //转换成float64类型
			if err != nil {
				panic(err)
			}

			logs.Logger.Infof("耗时%fms", lantencyFloat64)

			endpointsLantencyMonitor.With(prometheus.Labels{EndpointsDataSubsystem: endpoint}).Observe(lantencyFloat64)
		}(c)

		c.Next()
	}
}
