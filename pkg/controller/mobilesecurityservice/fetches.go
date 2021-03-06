package mobilesecurityservice

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	routev1 "github.com/openshift/api/route/v1"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func fetch(r *ReconcileMobileSecurityService, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service App resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service App")
	return reconcile.Result{}, err
}

//fetchRoute returns the Route resource created for this instance
func (r *ReconcileMobileSecurityService) fetchRoute(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService) (*routev1.Route, error) {
	reqLogger.Info("Checking if the route already exists")
	route := &routev1.Route{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetRouteName(instance), Namespace: instance.Namespace}, route)
	return route, err
}

//fetchService returns the service resource created for this instance
func (r *ReconcileMobileSecurityService) fetchService(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService) (*corev1.Service, error) {
	reqLogger.Info("Checking if the service already exists")
	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	return service, err
}

//fetchDeployment returns the deployment resource created for this instance
func (r *ReconcileMobileSecurityService) fetchDeployment(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService) (*v1beta1.Deployment, error) {
	reqLogger.Info("Checking if the deployment already exists")
	deployment := &v1beta1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	return deployment, err
}

//fetchConfigMap returns the config map resource created for this instance
func (r *ReconcileMobileSecurityService) fetchConfigMap(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService) (*corev1.ConfigMap, error) {
	reqLogger.Info("Checking if the ConfigMap already exists")
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetConfigMapName(instance), Namespace: instance.Namespace}, configMap)
	return configMap, err
}


