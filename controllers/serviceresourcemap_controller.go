/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	bindingoperatorscoreoscomv1alpha1 "github.com/openshift-app-service-poc/service-mapper/api/v1alpha1"
	"github.com/openshift-app-service-poc/service-mapper/pkg/binding"
)

// ServiceResourceMapReconciler reconciles a ServiceResourceMap object
type ServiceResourceMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	config    *rest.Config
	informers map[string]informer
}

type informer struct {
	informer   cache.SharedIndexInformer
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (i *informer) run() {
	i.informer.Run(i.ctx.Done())
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies/finalizers,verbs=update
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceResourceMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// Get ServiceResourceMap
	l.Info("get ServiceResourceMap", "srm name", req.Name)
	var sm bindingoperatorscoreoscomv1alpha1.ServiceResourceMap
	if err := r.Get(ctx, req.NamespacedName, &sm); err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		// delete serviceproxies and Service Endpoint definitions
		l.Info("ServiceResourceMap deleted, deleting also ServiceProxy", "srm name", req.Name)
		return ctrl.Result{}, r.deleteLinkedResources(ctx, req.Name)
	}

	// reconciling resources
	return ctrl.Result{}, r.reconcileLinkedResources(ctx, &sm)
}

func (r *ServiceResourceMapReconciler) reconcileLinkedResources(
	ctx context.Context,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap) error {
	l := log.FromContext(ctx)

	clusterClient, err := dynamic.NewForConfig(r.config)
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[0],
		Version:  strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[1],
		Resource: sm.Spec.ServiceKindReference.Kind,
	}

	crds, err := clusterClient.
		Resource(gvr).
		Namespace(corev1.NamespaceAll).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		l.Error(err, "error listing resource", "GroupVersionResource", gvr)
		return err
	}

	for _, c := range crds.Items {
		if err := r.createOrUpdateServiceProxyAndSED(ctx, sm, &c); err != nil {
			return err
		}
	}

	// running informer for monitored resources if not running
	if err := r.runInformer(ctx, clusterClient, gvr, sm); err != nil {
		return err
	}

	return nil
}

func (r *ServiceResourceMapReconciler) runInformer(
	ctx context.Context,
	clusterClient dynamic.Interface,
	gvr schema.GroupVersionResource,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap) error {
	l, _ := logr.FromContext(ctx)
	if _, ok := r.informers[sm.Name]; ok {
		// informer already running for this GVR
		l.Info("informer yet running", "GroupVersionResource", gvr)
		return nil
	}

	// run dynamic informer
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, time.Minute, sm.Namespace, nil)
	i := factory.ForResource(gvr).Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			r.createOrUpdateServiceProxyAndSED(ctx, sm, obj)
		},
		UpdateFunc: func(past, future interface{}) {
			r.createOrUpdateServiceProxyAndSED(ctx, sm, future)
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			r.deleteLinkedResources(ctx, u.GetName())
		},
	})

	l.Info("run informer", "GroupVersionResource", gvr)
	c, fc := context.WithCancel(ctx)

	li := informer{informer: i, ctx: c, cancelFunc: fc}
	r.informers[sm.GetName()] = li
	go li.run()

	return nil
}

func (r *ServiceResourceMapReconciler) createOrUpdateServiceProxyAndSED(
	ctx context.Context,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	obj interface{}) error {
	sp, err := r.createOrUpdateServiceProxy(ctx, sm, obj)
	if err != nil {
		return err
	}

	sec, err := r.createOrUpdateSED(ctx, sp, sm, obj)
	if err != nil {
		return err
	}

	sp.Status.Binding.Name = sec.Name
	if err := r.Status().Update(ctx, sp); err != nil {
		return fmt.Errorf("error updating serviceproxy.status.binding.name to '%s': %w", sec.Name, err)
	}

	return nil
}

func (r *ServiceResourceMapReconciler) createOrUpdateServiceProxy(
	ctx context.Context,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	obj interface{}) (*bindingoperatorscoreoscomv1alpha1.ServiceProxy, error) {
	l, _ := logr.FromContext(ctx)
	u := obj.(*unstructured.Unstructured)

	var sp bindingoperatorscoreoscomv1alpha1.ServiceProxy
	spSpec := bindingoperatorscoreoscomv1alpha1.ServiceProxySpec{
		ServiceResourceMapRef: sm.GetName(),
		ServiceInstance: bindingoperatorscoreoscomv1alpha1.NamespacedName{
			Name:      u.GetName(),
			Namespace: u.GetNamespace(),
		},
	}

	// check if ServiceProxy already exists
	spkey := client.ObjectKey{Namespace: u.GetNamespace(), Name: u.GetName()}
	if err := r.Get(ctx, spkey, &sp); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting ServiceProxy %s: %w", sp.Namespace+"/"+sp.Name, err)
		}

		// create ServiceProxy
		sp = bindingoperatorscoreoscomv1alpha1.ServiceProxy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      u.GetName(),
				Namespace: u.GetNamespace(),
			},
			Spec: spSpec,
		}
		if err := r.Create(ctx, &sp); err != nil {
			return nil, fmt.Errorf("error creating ServiceProxy %s: %w", sp.Namespace+"/"+sp.Name, err)
		}

		return &sp, nil
	}

	// update ServiceProxy
	sp.Spec = spSpec
	if err := r.Update(ctx, &sp); err != nil {
		l.Error(err, "error updating ServiceProxy")
	}
	return &sp, nil
}

func (r *ServiceResourceMapReconciler) createOrUpdateSED(
	ctx context.Context,
	sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	o interface{}) (*corev1.Secret, error) {

	obj := o.(*unstructured.Unstructured)

	// Generate Service Endpoint Definition
	sed := binding.NewServiceEndpointDefinition(ctx, r.Client, sm, sp, obj.UnstructuredContent())

	okey := client.ObjectKey{Namespace: sed.ObjectMeta.Namespace, Name: sed.ObjectMeta.Name}
	var s corev1.Secret

	if err := r.Get(ctx, okey, &s); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}

		if err2 := r.Create(ctx, sed); err2 != nil {
			return nil, err2
		}

		return sed, nil
	}

	if err := r.Update(ctx, sed); err != nil {
		return nil, err
	}
	return sed, nil
}

func (r *ServiceResourceMapReconciler) deleteLinkedResources(ctx context.Context, smName string) error {
	l := log.FromContext(ctx)

	// retrieve serviceproxies
	l.Info("service map deleted, deleting also ServiceProxy", "serviceresourcemap name", smName)
	var sps bindingoperatorscoreoscomv1alpha1.ServiceProxyList
	opts := &client.MatchingFields{".spec.service_resource_map": smName}
	if err := r.List(ctx, &sps, opts); err != nil {
		return err
	}

	// delete serviceproxies and seds
	for _, sp := range sps.Items {
		l.Info("processing impacted ServiceProxies", "serviceproxy namespace", sp.Namespace, "serviceproxy name", sp.Name)

		if sn := sp.Status.Binding.Name; sn != "" {
			l.Info("deleting linked Service Endpoint Definition", "sed", sn, "namespace", sp.Namespace)
			if err := r.deleteSecretIfExists(ctx, sp.Namespace, sn); err != nil {
				return err
			}
		}

		// delete ServiceProxy
		if err := r.Delete(ctx, &sp); err != nil {
			return err
		}
	}

	if i, ok := r.informers[smName]; ok {
		i.cancelFunc()
		delete(r.informers, smName)
	}
	return nil
}

func (r *ServiceResourceMapReconciler) deleteSecretIfExists(ctx context.Context, namespace, name string) error {
	var sec corev1.Secret
	skey := client.ObjectKey{Namespace: namespace, Name: name}
	if err := r.Get(ctx, skey, &sec); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := r.Delete(ctx, &sec); err != nil {
		return err
	}

	return nil
}

func (r *ServiceResourceMapReconciler) deleteServiceProxy(ctx context.Context, obj interface{}) {
	l, _ := logr.FromContext(ctx)
	u := obj.(*unstructured.Unstructured)

	l.Info("deleting service proxy")

	okey := client.ObjectKey{Namespace: u.GetNamespace(), Name: u.GetName()}
	var sp bindingoperatorscoreoscomv1alpha1.ServiceProxy
	if err := r.Get(ctx, okey, &sp); err != nil {
		l.Error(err, "error getting proxable service for target service", "target service key", okey)
		return
	}

	if err := r.Delete(ctx, &sp); err != nil {
		l.Error(err, "error deleting proxable service", "target service key", okey)
		return
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceResourceMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.config = mgr.GetConfig()
	r.informers = make(map[string]informer)

	mgr.
		GetFieldIndexer().
		IndexField(context.Background(),
			&bindingoperatorscoreoscomv1alpha1.ServiceProxy{},
			".spec.service_resource_map",
			func(o client.Object) []string {
				sm := o.(*bindingoperatorscoreoscomv1alpha1.ServiceProxy)
				return []string{sm.Spec.ServiceResourceMapRef}
			})

	return ctrl.NewControllerManagedBy(mgr).
		For(&bindingoperatorscoreoscomv1alpha1.ServiceResourceMap{}).
		Complete(r)
}
