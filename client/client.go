package client

import (
	yaml "gopkg.in/yaml.v2"

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest"
	"github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/helm"
)

const (
	tillerHost          = "tiller-deploy.kube-system.svc.cluster.local:44134"
	mysqlChartPath      = "/mysql-0.3.4.tgz"
	mariadbChartPath    = "/mariadb-2.1.3.tgz"
	postgresqlChartPath = "/postgresql-0.8.8.tgz"
)

// Install creates a new MariaDB chart release
func Install(releaseName string, serviceID string, namespace string, req map[string]interface{}) error {
	vals, err := yaml.Marshal(map[string]interface{}{
		"mariadbRootPassword": "root",
		"mariadbDatabase":     "dbname",
		"persistence": map[string]interface{}{
			"enabled": false,
		},
	})

	logrus.Info(serviceID)

	var chartPath string

	switch serviceID {
	case "3533e2f0-6001-xxxx-9d15-d7c0b90b75b5":
		chartPath = "/mysql-0.3.4.tgz"
	case "3533e2f0-6002-xxxx-9d15-d7c0b90b75b5":
		chartPath = "/mariadb-2.1.3.tgz"
	case "3533e2f0-6335-xxxx-9d15-d7c0b90b75b5":
		chartPath = "/postgresql-0.8.8.tgz"
	default:
		panic("unrecognized chart character")
	}

	logrus.Info(chartPath)

	// vals, err := yaml.Marshal(req)

	if err != nil {
		return err
	}
	helmClient := helm.NewClient(helm.Host(tillerHost))
	_, err = helmClient.InstallRelease(chartPath, namespace, helm.ReleaseName(releaseName), helm.ValueOverrides(vals))
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a MariaDB chart release
func Delete(releaseName string) error {
	helmClient := helm.NewClient(helm.Host(tillerHost))
	if _, err := helmClient.DeleteRelease(releaseName); err != nil {
		return err
	}
	return nil
}

// GetPassword returns the MariaDB password for a chart release
// func GetPassword(releaseName, namespace string) (string, error) {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		return "", err
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return "", err
// 	}
// 	secret, err := clientset.Core().Secrets(namespace).Get(releaseName + "-mariadb")
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(secret.Data["mariadb-root-password"]), nil
// }
