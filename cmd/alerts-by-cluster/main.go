package main

import (
	"fmt"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/config"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/loggerpkg"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/sysdighttp"
	"github.com/aaronm-sysdig/alerts-by-cluster/structs/alerts"
	"github.com/aaronm-sysdig/alerts-by-cluster/structs/metadata"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func retrieveClusters(logger *logrus.Logger, config *configuration.Config, client sysdighttp.SysdigClient) (*metadata.ResultMetadata, error) {
	var err error
	// Get list of kubernetes clusters in environment
	configClusters := sysdighttp.DefaultSysdigRequestConfig(fmt.Sprintf("%s/api/data/entity/metadata", viper.GetString("SECURE_URL")), config.SecureAPIToken)
	configClusters.Method = "POST"
	configClusters.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	configClusters.JSON = metadata.PayloadMetadata{
		Paging: metadata.PagingPayload{
			From: 0,
			To:   9999,
		},
		Metrics: []string{"kubernetes.cluster.name"},
	}

	var objMetadataResponse *http.Response
	if objMetadataResponse, err = client.SysdigRequest(logger, configClusters); err != nil {
		logger.Fatalf("Error creating sysdig request: %s", err)
	}
	defer objMetadataResponse.Body.Close()

	jsonMetadataResponse := &metadata.ResultMetadata{}
	if err = client.ResponseBodyToJson(objMetadataResponse, jsonMetadataResponse); err != nil {
		logger.Fatalf("Error unmarshalling sysdig response: %s", err)
	}
	return jsonMetadataResponse, nil
}

func getAlerts(logger *logrus.Logger, config *configuration.Config, client sysdighttp.SysdigClient) (*alerts.AlertQuery, error) {
	var err error
	// Get list of kubernetes clusters in environment
	configAlerts := sysdighttp.DefaultSysdigRequestConfig(fmt.Sprintf("%s/api/scanning/v1/alerts", config.SecureURL), config.SecureAPIToken)
	configAlerts.Method = "GET"
	var objAlertsResponse *http.Response
	if objAlertsResponse, err = client.SysdigRequest(logger, configAlerts); err != nil {
		return nil, err
	}
	defer objAlertsResponse.Body.Close()

	jsonAlerts := &alerts.AlertQuery{}
	if err = client.ResponseBodyToJson(objAlertsResponse, jsonAlerts); err != nil {
		return nil, err
	}
	return jsonAlerts, nil
}

func alertExists(alerts *alerts.AlertQuery, clusterName string) bool {
	for _, alert := range alerts.Alerts {
		if alert.Scope == fmt.Sprintf("kubernetes.cluster.name = \"%s\"", clusterName) {
			return true
		}
	}
	return false
}

func createAlertForCluster(logger *logrus.Logger, config *configuration.Config, clusterName string, client sysdighttp.SysdigClient) error {
	var err error
	configCreateAlert := sysdighttp.DefaultSysdigRequestConfig(fmt.Sprintf("%s/api/scanning/v1/alerts", config.SecureURL), config.SecureAPIToken)
	configCreateAlert.Method = "POST"
	configCreateAlert.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	configCreateAlert.JSON = alerts.PayloadAlert{
		Enabled:      true,
		Type:         "runtime",
		Name:         fmt.Sprintf("Cluster: %s", clusterName),
		Description:  "",
		Scope:        fmt.Sprintf("kubernetes.cluster.name = \"%s\"", clusterName),
		Repositories: []string{},
		Triggers: alerts.PayloadTriggers{
			Unscanned:      true,
			AnalysisUpdate: false,
			VulnUpdate:     true,
			PolicyEval:     true,
		},
		Autoscan:               false,
		OnlyPassFail:           false,
		NotificationChannelIds: []string{},
	}

	var objAlertResponse *http.Response
	if objAlertResponse, err = client.SysdigRequest(logger, configCreateAlert); err != nil {
		return err
	}
	defer objAlertResponse.Body.Close()
	return nil
}

func main() {
	var err error
	var config *configuration.Config
	var arrClusters *metadata.ResultMetadata
	var arrAlerts *alerts.AlertQuery

	logger := loggerpkg.GetLogger()
	client := sysdighttp.NewSysdigClient()

	if config, err = configuration.LoadConfig(logger); err != nil {
		logger.Fatalf("Could not load configuration. Error: '%v'", err)
	}

	if arrClusters, err = retrieveClusters(logger, config, client); err != nil {
		logger.Fatalf("Could not retrieve clusters. error: '%v'", err)
	}

	if arrAlerts, err = getAlerts(logger, config, client); err != nil {
		logger.Fatalf("Could not retrieve alerts.  error '%v'", err)
	}
	_ = arrAlerts

	for _, cluster := range arrClusters.Data {
		if alertExists(arrAlerts, cluster.KubernetesClusterName) == true {
			logger.Debugf("Cluster '%s' exists", cluster.KubernetesClusterName)
		} else {
			logger.Debugf("Cluster '%s' doesn't exist, creating", cluster.KubernetesClusterName)
			if err = createAlertForCluster(logger, config, cluster.KubernetesClusterName, client); err != nil {
				logger.Fatalf("Could not create alert for cluster '%s'. Error: '%v'", cluster.KubernetesClusterName, err)
			}
		}
	}

	logger.Infof("Finished...")
}
