package mobilesecurityservice

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	CONFIGMAP     = "ConfigMap"
	DEEPLOYMENT   = "Deployment"
	SERVICE       = "Service"
	ROUTE         = "Route"
)

var log = logf.Log.WithName("controller_mobilesecurityservice")

// Add creates a new MobileSecurityService Controller and adds it to the Manager.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// Returns the a new Reconciler for this operator and controller
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// Add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityService
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and create the owner MobileSecurityService **/

	//ConfigMap
	if err := watchConfigMap(c); err != nil {
		return err
	}

	//Deployment
	if err := watchDeployment(c); err != nil {
		return err
	}

	//Service
	if err := watchService(c); err != nil {
		return err
	}

	//Route
	if err:= watchRoute(c); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityService{}

//ReconcileMobileSecurityService reconciles a MobileSecurityService object
type ReconcileMobileSecurityService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Update the factory object and requeue
func (r *ReconcileMobileSecurityService) update(obj runtime.Object, reqLogger logr.Logger) (reconcile.Result, error) {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Spec")
		return reconcile.Result{}, err
	}
	reqLogger.Info("Spec updated - return and create")
	return reconcile.Result{Requeue: true}, nil
}

//Create the factory object and requeue
func (r *ReconcileMobileSecurityService) create( instance *mobilesecurityservicev1alpha1.MobileSecurityService, reqLogger logr.Logger, kind string, err error) (reconcile.Result, error) {
	obj, errBuildObject := r.buildFactory(reqLogger, instance, kind)
	if errBuildObject != nil {
		return reconcile.Result{}, errBuildObject
	}
	if errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ", "kind", kind, "Namespace", instance.Namespace)
		err = r.client.Create(context.TODO(), obj)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", instance.Namespace)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Created successfully - return and create", "kind", kind, "Namespace", instance.Namespace)
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Error(err, "Failed to get", "kind", kind, "Namespace", instance.Namespace)
	return reconcile.Result{}, err

}

//buildFactory will return the resource according to the kind defined
func (r *ReconcileMobileSecurityService) buildFactory(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService, kind string) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", instance.Namespace)
	switch kind {
	case CONFIGMAP:
		return r.buildConfigMap(instance), nil
	case DEEPLOYMENT:
		return r.buildDeployment(instance), nil
	case SERVICE:
		return r.buildService(instance), nil
	case ROUTE:
		return r.buildRoute(instance), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + instance.Namespace
		panic(msg)
	}
}



// Reconcile reads that state of the cluster for a MobileSecurityService object and makes changes based on the state read
// and what is in the MobileSecurityService.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mobile Security Service ...")

	// FIXME: Check if is a valid namespace
	// We should not checked if the namespace is valid or not. It is an workaround since currently is not possible watch/cache a List of Namespaces
	// The impl to allow do it is done and merged in the master branch of the lib but not released in an stable version. It should be removed when this feature be impl.
	// See the PR which we are working on to update the deps and have this feature: https://github.com/operator-framework/operator-sdk/pull/1388
	if isValidNamespace, err:= utils.IsValidOperatorNamespace(request.Namespace); err != nil || isValidNamespace == false {
		// Stop reconcile
		operatorNamespace, _ := k8sutil.GetOperatorNamespace();
		reqLogger.Error(err, "Unable to reconcile Mobile Security Service", "Request.Namespace", request.Namespace, "isValidNamespace", isValidNamespace, "Operator.Namespace", operatorNamespace)
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Valid namespace for Mobile Security Service", "Namespace", request.Namespace)
	reqLogger.Info("Start Reconciling Mobile Security Service ...")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityService{}

	//Fetch the MobileSecurityService instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	//Check specs
	if !hasMandatorySpecs(instance, reqLogger) {
		return reconcile.Result{Requeue: true}, nil
	}

	//Check if ConfigMap for the app exist, if not create one.
	if _, err := r.fetchConfigMap(reqLogger, instance); err != nil {
		return r.create(instance, reqLogger, CONFIGMAP, err)
	}

	//Check if Deployment for the app exist, if not create one
	deployment, err := r.fetchDeployment(reqLogger, instance)
	if err != nil {
		return r.create(instance, reqLogger, DEEPLOYMENT, err)
	}

	reqLogger.Info("Ensuring the Mobile Security Service deployment size is the same as the spec")
	size := instance.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		return r.update(deployment, reqLogger)
	}

	//Check if Service for the app exist, if not create one
	if _, err := r.fetchService(reqLogger, instance); err != nil {
		return r.create(instance, reqLogger, SERVICE, err)
	}

	//Check if Route for the Service exist, if not create one
	if _, err := r.fetchRoute(reqLogger, instance); err != nil {
		return r.create(instance, reqLogger, ROUTE, err)
	}

	//Update status for ConfigMap
	configMapStatus, err := r.updateConfigMapStatus(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for deployment
	deploymentStatus, err := r.updateDeploymentStatus(reqLogger,instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for Service
	serviceStatus, err := r.updateServiceStatus(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for Route
	routeStatus, err := r.updateRouteStatus(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for App
	if err:= r.updateStatus(reqLogger, configMapStatus, deploymentStatus, serviceStatus, routeStatus, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
