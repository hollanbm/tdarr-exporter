package collector

import (
	"encoding/json"

	"github.com/homeylab/tdarr-exporter/internal/client"
	"github.com/homeylab/tdarr-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var METRIC_PREFIX = "tdarr"

type TdarrCollector struct {
	config           config.Config
	payload          TdarrMetricRequest
	totalFilesMetric *prometheus.Desc
}

func NewTdarrCollector(runConfig config.Config) *TdarrCollector {
	return &TdarrCollector{
		config:  runConfig,
		payload: getRequestPayload(),
		totalFilesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(METRIC_PREFIX, "", "total_files"),
			"Tdarr totalFileCount",
			nil,
			prometheus.Labels{"instance": runConfig.Url},
		),
	}
}

func (collector *TdarrCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalFilesMetric
}

func (collector *TdarrCollector) Collect(ch chan<- prometheus.Metric) {
	httpClient, err := client.NewClient(collector.config.Url, collector.config.VerifySsl, collector.config.HttpTimeoutSeconds)
	if err != nil {
		log.Fatal().
			Err(err)
	}
	log.Debug().Interface("payload", collector.payload).Msg("Requesting statistics data from Tdarr")
	// Marshal it into JSON prior to requesting
	payload, err := json.Marshal(collector.payload)
	if err != nil {
		log.Fatal().
			Err(err)
	}
	responseData := &TdarrDataResponse{}
	httpErr := httpClient.DoPostRequest(collector.config.TdarrMetricsPath, responseData, payload)
	if httpErr != nil {
		log.Error().Err(httpErr).Msg("Failed to get data for Tdarr exporter")
		return
	}
	log.Info().Interface("response", responseData).Msg("Output")
	ch <- prometheus.MustNewConstMetric(collector.totalFilesMetric, prometheus.GaugeValue, 12.221)
	// time.Sleep(50 * time.Second)
}

func getRequestPayload() TdarrMetricRequest {
	return TdarrMetricRequest{
		Data: TdarrDataRequest{
			Collection: "StatisticsJSONDB",
			Mode:       "getById",
			DocId:      "statistics",
		},
	}
}

// func getRootPieMetrics() map{
// 	return {

// 	}
// }

// func getPieParseMap() map[int]func() {
// 	return {
// 		0: getTranscodeMetrics(),
// 		1: getHealthCheckMetrics(),
// 		2: getVideoCodesMetrics(),
// 	}
// }
