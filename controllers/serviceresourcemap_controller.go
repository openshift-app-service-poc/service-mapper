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

	bindingoperatorscoreoscomv1alpha1 "github.com/filariow/sbo-1225/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ServiceResourceMapReconciler reconciles a ServiceResourceMap object
type ServiceResourceMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	config    *rest.Config
	informers map[schema.GroupVersionResource]cache.SharedIndexInformer
}

//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=proxableservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceResourceMap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceResourceMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// Get ServiceResourceMap
	l.Info("get service map")
	var sm bindingoperatorscoreoscomv1alpha1.ServiceResourceMap
	if err := r.Get(ctx, req.NamespacedName, &sm); err != nil {
		// TODO: delete proxable service
		l.Info("service map deleted, deleting also ProxableService", "service map namespace", req.Namespace, "service map name", req.Name)
		var pss bindingoperatorscoreoscomv1alpha1.ProxableServiceList
		opts := &client.MatchingFields{".spec.service_resource_map": req.Namespace + "," + req.Name}
		if err := r.List(ctx, &pss, opts); err != nil {
			return ctrl.Result{}, err
		}

		l.Info("retrieved linked ProxableServices")

		for _, ps := range pss.Items {
			l.Info("processing linked ProxableService", "proxableservice name", ps.Name)
			if err := r.Delete(ctx, &ps); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// reconciling service map
	clusterClient, err := dynamic.NewForConfig(r.config)
	if err != nil {
		return ctrl.Result{}, err
	}

	gvr := schema.GroupVersionResource{
		Group:    strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[0],
		Version:  strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[1],
		Resource: sm.Spec.ServiceKindReference.Kind,
	}

	crds, err := clusterClient.Resource(gvr).Namespace(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		l.Error(err, "listing resource")
		return ctrl.Result{}, err
	}

	for _, c := range crds.Items {
		r.createUpdateProxableService(ctx, &sm, &c)
	}

	// running informer for monitored resources if not running
	if err := r.runInformer(ctx, clusterClient, gvr, &sm); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ServiceResourceMapReconciler) runInformer(ctx context.Context, clusterClient dynamic.Interface, gvr schema.GroupVersionResource, sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap) error {
	l, _ := logr.FromContext(ctx)
	if _, ok := r.informers[gvr]; ok {
		// informer already running for this GVR
		l.Info("informer yet running", "GroupVersionResource", gvr)
		return nil
	}

	// run dynamic informer
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, time.Minute, sm.Namespace, nil)
	informer := factory.ForResource(gvr).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			r.createUpdateProxableService(ctx, sm, obj)
		},
		UpdateFunc: func(past, future interface{}) {
			r.createUpdateProxableService(ctx, sm, future)
		},
		DeleteFunc: func(obj interface{}) {
			r.deleteProxableService(ctx, obj)
		},
	})

	l.Info("run informer", "GroupVersionResource", gvr)
	r.informers[gvr] = informer
	go informer.Run(ctx.Done())

	return nil
}

func (r *ServiceResourceMapReconciler) deleteProxableService(ctx context.Context, obj interface{}) {
	l, _ := logr.FromContext(ctx)
	u := obj.(*unstructured.Unstructured)

	l.Info("deleting proxable service")

	okey := client.ObjectKey{Namespace: corev1.NamespaceAll, Name: u.GetName() + "-" + u.GetNamespace()}
	var ps bindingoperatorscoreoscomv1alpha1.ProxableService
	if err := r.Get(ctx, okey, &ps); err != nil {
		l.Error(err, "error getting proxable service for target service", "target service key", okey)
		return
	}

	if err := r.Delete(ctx, &ps); err != nil {
		l.Error(err, "error deleting proxable service", "target service key", okey)
		return
	}
}

func (r *ServiceResourceMapReconciler) createUpdateProxableService(ctx context.Context, sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap, obj interface{}) {
	l, _ := logr.FromContext(ctx)
	u := obj.(*unstructured.Unstructured)

	var ps bindingoperatorscoreoscomv1alpha1.ProxableService
	psSpec := bindingoperatorscoreoscomv1alpha1.ProxableServiceSpec{
		ServiceResourceMapRef: bindingoperatorscoreoscomv1alpha1.NamespacedName{
			Name:      sm.GetName(),
			Namespace: sm.GetNamespace(),
		},
		ServiceInstance: bindingoperatorscoreoscomv1alpha1.NamespacedName{
			Name:      u.GetName(),
			Namespace: u.GetNamespace(),
		},
	}

	psname := u.GetName() + "-" + u.GetNamespace()
	if err := r.Get(ctx, client.ObjectKey{Namespace: corev1.NamespaceAll, Name: psname}, &ps); err != nil {
		if !apierrors.IsNotFound(err) {
			l.Error(err, "error getting ProxableService", "resource", u, "proxableservice", ps)
			return
		}

		ps = bindingoperatorscoreoscomv1alpha1.ProxableService{
			ObjectMeta: metav1.ObjectMeta{
				Name: psname,
			},
			Spec: psSpec,
		}
		if err := r.Create(ctx, &ps); err != nil {
			l.Error(err, "error creating ProxableService", "resource", u, "proxableservice", ps)
		}

		return
	}

	ps.Spec = psSpec
	if err := r.Update(ctx, &ps); err != nil {
		l.Error(err, "error updating ProxableService", "resource", u, "proxableservice", ps)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceResourceMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.config = mgr.GetConfig()
	r.informers = make(map[schema.GroupVersionResource]cache.SharedIndexInformer)

	mgr.
		GetFieldIndexer().
		IndexField(context.Background(),
			&bindingoperatorscoreoscomv1alpha1.ProxableService{},
			".spec.service_resource_map",
			func(o client.Object) []string {
				sm := o.(*bindingoperatorscoreoscomv1alpha1.ProxableService)
				return []string{sm.Spec.ServiceResourceMapRef.Name + "," + sm.Spec.ServiceResourceMapRef.Namespace}
			})

	return ctrl.NewControllerManagedBy(mgr).
		For(&bindingoperatorscoreoscomv1alpha1.ServiceResourceMap{}).
		Complete(r)
}
