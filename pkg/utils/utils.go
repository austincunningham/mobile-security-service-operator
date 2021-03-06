package utils

import (
	"fmt"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"os"
	"strings"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// APP_NAMESPACE_ENV_VAR is the constant for env variable APP_NAMESPACES
// which is the namespace where the APP CR can applied.
// The namespaces should be informed split by ";".
const APP_NAMESPACE_ENV_VAR = "APP_NAMESPACES"
const SERVICE_INSTANCE_NAME  = "mobile-security-service"
const SERVICE_INSTANCE_NAMESPACE  = "mobile-security-service-operator"

var log = logf.Log.WithName("mobile-security-service-operator.utils")

//GetRouteName returns an string name with the name of the router
func GetRouteName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string{
	if len(m.Spec.RouteName) > 0 {
		return m.Spec.RouteName
	}
	return m.Name
}

//GetConfigMapName returns an string name with the name of the configMap
func GetConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	if len(m.Spec.ConfigMapName) > 0 {
		return m.Spec.ConfigMapName
	}
	return m.Name
}


// GetCustomWatchNamespaces returns the namespace the operator should be watching for changes
func GetAppNamespaces() (string, error) {
	ns, found := os.LookupEnv(APP_NAMESPACE_ENV_VAR)
	if !found {
		return "", fmt.Errorf("%s must be set", APP_NAMESPACE_ENV_VAR)
	}
	return ns, nil
}


//IsValidAppNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
func IsValidAppNamespace(namespace string) (bool, error) {
	appNamespacesEnvVar,err := GetAppNamespaces()
	if err != nil {
		log.Error(err, "Unable to check if is app namespace %s is valid", namespace)
		return false, err
	}
	for _, ns := range strings.Split(appNamespacesEnvVar, ";") {
		if ns == namespace {
			return true, nil
		}
	}
	err = fmt.Errorf( "Invalid Namespace")
	return false, err
}


//IsValidAppNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
func IsValidOperatorNamespace(namespace string) (bool, error) {
	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return false, err
	}
	if ns == namespace {
		return true, nil
	}
	err = fmt.Errorf( "Invalid Namespace")
	return false, err
}