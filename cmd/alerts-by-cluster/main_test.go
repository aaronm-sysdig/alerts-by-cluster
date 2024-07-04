package main

import (
	"bytes"
	"encoding/json"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/config"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/loggerpkg"
	"github.com/aaronm-sysdig/alerts-by-cluster/pkg/sysdighttp"
	"github.com/aaronm-sysdig/alerts-by-cluster/structs/alerts"
	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"testing"
)

func TestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Main Suite")
}

var _ = ginkgo.Describe("Main", func() {
	var (
		ctrl             *gomock.Controller
		mockSysdigClient *sysdighttp.MockSysdigClient
		logger           *logrus.Logger
		config           *configuration.Config
	)

	ginkgo.BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		mockSysdigClient = sysdighttp.NewMockSysdigClient(ctrl)
		logger = loggerpkg.GetLogger()

		// Load configuration for tests
		var err error
		config, err = configuration.LoadConfig(logger)
		if err != nil {
			logger.Warnf("Could not load configuration. Error: '%v'", err)
		}
	})

	ginkgo.AfterEach(func() {
		ctrl.Finish()
	})

	ginkgo.It("should retrieve clusters successfully", func() {
		mockResponse := `{"metrics":["kubernetes.cluster.name"],"time":{"from":1720047000000000,"to":1720068600000000,"sampling":600000000},"data":[{"kubernetes.cluster.name":"aamiles-onprem5"}],"paging":{"from":0,"to":9999,"total":1}}`
		httpResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
		}

		mockSysdigClient.EXPECT().SysdigRequest(gomock.Any(), gomock.Any()).Return(httpResponse, nil).Times(1)
		mockSysdigClient.EXPECT().ResponseBodyToJson(httpResponse, gomock.Any()).DoAndReturn(func(resp *http.Response, target interface{}) error {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return json.Unmarshal(body, target)
		}).Times(1)

		result, err := retrieveClusters(logger, config, mockSysdigClient)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(len(result.Data)).Should(gomega.Equal(1))
		gomega.Expect(result.Data[0].KubernetesClusterName).Should(gomega.Equal("aamiles-onprem5"))
	})

	ginkgo.It("should get alerts successfully", func() {
		mockResponse := `{"alerts":[{"enabled":true,"type":"runtime","name":"Cluster: aamiles-onprem5","description":"","scope":"kubernetes.cluster.name = \"aamiles-onprem5\"","repositories":[],"triggers":{"unscanned":true,"analysis_update":false,"vuln_update":true,"policy_eval":true},"autoscan":false,"onlyPassFail":false,"notificationChannelIds":[]}]}`

		httpResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
		}

		mockSysdigClient.EXPECT().SysdigRequest(gomock.Any(), gomock.Any()).Return(httpResponse, nil).Times(1)
		mockSysdigClient.EXPECT().ResponseBodyToJson(httpResponse, gomock.Any()).DoAndReturn(func(resp *http.Response, target interface{}) error {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return json.Unmarshal(body, target)
		}).Times(1)

		result, err := getAlerts(logger, config, mockSysdigClient)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(len(result.Alerts)).Should(gomega.Equal(1))
		gomega.Expect(result.Alerts[0].Name).Should(gomega.Equal("Cluster: aamiles-onprem5"))
	})

	ginkgo.It("should detect if an alert exists", func() {
		alerts := &alerts.AlertQuery{
			Alerts: []alerts.Alert{
				{Name: "Cluster: aamiles-onprem5", Scope: "kubernetes.cluster.name = \"aamiles-onprem5\""},
			},
		}
		exists := alertExists(alerts, "aamiles-onprem5")
		gomega.Expect(exists).Should(gomega.BeTrue())
	})

	ginkgo.It("should detect if an alert does not exist", func() {
		alerts := &alerts.AlertQuery{
			Alerts: []alerts.Alert{
				{Name: "Cluster: other-cluster", Scope: "kubernetes.cluster.name = \"other-cluster\""},
			},
		}
		exists := alertExists(alerts, "aamiles-onprem5")
		gomega.Expect(exists).Should(gomega.BeFalse())
	})

	ginkgo.It("should create alert for cluster", func() {
		clusterName := "aamiles-onprem5"
		httpResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		}

		mockSysdigClient.EXPECT().SysdigRequest(gomock.Any(), gomock.Any()).Return(httpResponse, nil).Times(1)
		err := createAlertForCluster(logger, config, clusterName, mockSysdigClient)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
})
