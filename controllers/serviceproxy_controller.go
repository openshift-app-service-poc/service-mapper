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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bindingoperatorscoreoscomv1alpha1 "github.com/filariow/sbo-1225/api/v1alpha1"
	"github.com/filariow/sbo-1225/pkg/binding"
)

var ErrServiceResourceMapNotFound = fmt.Errorf("ServiceResourceMap not found")
var ErrProxableServiceNotFound = fmt.Errorf("ProxableService not found")

// ServiceProxyReconciler reconciles a ServiceProxy object
type ServiceProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	config *rest.Config
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceresourcemaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=binding.operators.coreos.com,resources=serviceproxies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// Get Service Proxy
	l.Info("reading ServiceProxy", "resource", req.NamespacedName)
	var sp bindingoperatorscoreoscomv1alpha1.ServiceProxy
	if err := r.Client.Get(ctx, req.NamespacedName, &sp); err != nil {
		// service proxy has been cancelled, we also have to delete the SED
		if err := r.deleteSED(ctx, req.Name, req.Namespace); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get related Proxable Service
	l.Info("reading get ProxableService", "service-proxy", sp)
	ps, err := r.getProxableService(ctx, &sp)
	if err != nil {
		// proxable service has been cancelled, we also have to delete the SED
		if err := r.deleteSED(ctx, req.Name, req.Namespace); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get Service Resource Map
	l.Info("reading ServiceResourceMap", "proxable-service", ps)
	var sm bindingoperatorscoreoscomv1alpha1.ServiceResourceMap
	okey := client.ObjectKey{Namespace: ps.Spec.ServiceResourceMapRef.Namespace, Name: ps.Spec.ServiceResourceMapRef.Name}
	if err := r.Get(ctx, okey, &sm); err != nil {
		l.Error(err, "error retrieving service-map", "object-key", okey)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get Service Instance
	clusterClient, err := dynamic.NewForConfig(r.config)
	if err != nil {
		return ctrl.Result{}, err
	}

	gvr := schema.GroupVersionResource{
		Group:    strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[0],
		Version:  strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[1],
		Resource: sm.Spec.ServiceKindReference.Kind,
	}

	si, err := clusterClient.
		Resource(gvr).
		Namespace(ps.Spec.ServiceInstance.Namespace).
		Get(ctx, ps.Spec.ServiceInstance.Name, metav1.GetOptions{})
	if err != nil {
		l.Error(err, "error retrieving service instance", "object-key", okey)
		return ctrl.Result{}, err
	}

	l.Info("retrieved unstructured service instance", "unstructured", si, "unstructured content", si.UnstructuredContent())

	// Create SED
	l.Info("creating or updating Service Endpoint Definition")
	if err := r.createOrUpdateSED(ctx, &sp, &sm, si); err != nil {
		return ctrl.Result{}, err
	}

	// Update ServiceProxy
	l.Info("updating service proxy status.binding.name")
	if sp.Status.Binding.Name != sp.GetName()+"-sed" {
		sp.Status = bindingoperatorscoreoscomv1alpha1.ServiceProxyStatus{
			Binding: bindingoperatorscoreoscomv1alpha1.ServiceProxyStatusBinding{
				Name: sp.GetName() + "-sed",
			},
		}

		if err := r.Status().Update(ctx, &sp); err != nil {
			return ctrl.Result{}, err
		}
	}

	if err := controllerutil.SetControllerReference(&sp, &sm, r.Client.Scheme()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ServiceProxyReconciler) deleteSED(ctx context.Context, spName, spNamespace string) error {
	var sec corev1.Secret
	okey := client.ObjectKey{Namespace: spNamespace, Name: spName + "-sed"}
	if err := r.Get(ctx, okey, &sec); err != nil {
		return client.IgnoreNotFound(err)
	}

	if err := r.Delete(ctx, &sec); err != nil {
		return err
	}

	return nil
}

func (r *ServiceProxyReconciler) getProxableService(ctx context.Context, sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy) (*bindingoperatorscoreoscomv1alpha1.ProxableService, error) {
	var ps bindingoperatorscoreoscomv1alpha1.ProxableService
	okey := client.ObjectKey{Namespace: corev1.NamespaceAll, Name: sp.Spec.ProxableService}
	if err := r.Get(ctx, okey, &ps); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: ProxableService '%s'", ErrProxableServiceNotFound, sp.Spec.ProxableService)
		}
		return nil, err
	}

	return &ps, nil
}

func (r *ServiceProxyReconciler) createOrUpdateSED(
	ctx context.Context,
	sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	obj *unstructured.Unstructured) error {
	// Generate Service Endpoint Definition
	sed, err := r.bakeSED(ctx, sp, sm, obj)
	if err != nil {
		return err
	}

	okey := client.ObjectKey{Namespace: sed.ObjectMeta.Namespace, Name: sed.ObjectMeta.Name}
	var s corev1.Secret

	if err := r.Get(ctx, okey, &s); err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, sed)
		}
		return err
	}

	return r.Update(ctx, sed)
}

func (r *ServiceProxyReconciler) bakeSED(
	ctx context.Context,
	sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	obj *unstructured.Unstructured) (*corev1.Secret, error) {

	sec := binding.NewServiceEndpointDefinition(ctx, r.Client, sm, sp, obj.UnstructuredContent())
	return sec, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.config = mgr.GetConfig()

	mgr.
		GetFieldIndexer().
		IndexField(context.Background(),
			&bindingoperatorscoreoscomv1alpha1.ServiceProxy{},
			".spec.proxable_service",
			func(o client.Object) []string {
				sm := o.(*bindingoperatorscoreoscomv1alpha1.ServiceProxy)
				return []string{sm.Spec.ProxableService}
			})

	return ctrl.NewControllerManagedBy(mgr).
		For(&bindingoperatorscoreoscomv1alpha1.ServiceProxy{}).
		Watches(&source.Kind{Type: &bindingoperatorscoreoscomv1alpha1.ServiceResourceMap{}},
			handler.EnqueueRequestsFromMapFunc(func(o client.Object) []reconcile.Request {
				// fmt.Printf("change in ServiceResourceMap '%s/%s'\n", o.GetNamespace(), o.GetName())
				c := mgr.GetClient()
				clusterClient, err := dynamic.NewForConfig(mgr.GetConfig())
				if err != nil {
					return nil
				}

				sm := o.(*bindingoperatorscoreoscomv1alpha1.ServiceResourceMap)
				gvr := schema.GroupVersionResource{
					Group:    strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[0],
					Version:  strings.Split(sm.Spec.ServiceKindReference.ApiGroup, "/")[1],
					Resource: sm.Spec.ServiceKindReference.Kind,
				}

				crds, err := clusterClient.Resource(gvr).Namespace(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return nil
				}

				rr := []reconcile.Request{}
				for _, crd := range crds.Items {
					// get ServiceProxies for ProxableServices for services related to updated ServiceResourceMap
					var sps bindingoperatorscoreoscomv1alpha1.ServiceProxyList
					opts := &client.ListOptions{
						FieldSelector: fields.OneTermEqualSelector(".spec.proxable_service", crd.GetName()+"-"+crd.GetNamespace()),
						Namespace:     corev1.NamespaceAll,
					}
					if err := c.List(context.TODO(), &sps, opts); err != nil {
						// fmt.Printf("error retrieving list of service proxies: %s\n", err)
						return nil
					}

					for _, sp := range sps.Items {
						rr = append(rr, reconcile.Request{
							NamespacedName: types.NamespacedName{
								Name:      sp.GetName(),
								Namespace: sp.GetNamespace(),
							},
						})
					}
				}

				// map service proxies to reconcile requests
				// fmt.Printf("asking for reconciling service proxies: %v\n", rr)
				return rr
			}),
		).
		Complete(r)
}
