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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bindingoperatorscoreoscomv1alpha1 "github.com/filariow/sbo-1225/api/v1alpha1"
)

var ErrServiceResourceMapNotFound = fmt.Errorf("ServiceResourceMap not found")

// ServiceProxyReconciler reconciles a ServiceProxy object
type ServiceProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
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
	_ = log.FromContext(ctx)

	// Get Service Proxy
	var sp bindingoperatorscoreoscomv1alpha1.ServiceProxy
	if err := r.Client.Get(ctx, req.NamespacedName, &sp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get related Service Resource Map
	sm, err := r.getServiceResourceMap(ctx, &sp)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Generate Service Endpoint Definition
	sed, err := r.bakeSED(ctx, &sp, sm)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create SED
	if err := r.createOrUpdateSED(ctx, sed); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ServiceProxyReconciler) getServiceResourceMap(ctx context.Context, sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy) (*bindingoperatorscoreoscomv1alpha1.ServiceResourceMap, error) {
	ag := sp.Spec.ServiceReference.ApiGroup
	k := sp.Spec.ServiceReference.Kind

	opts := []client.ListOption{
		client.InNamespace(""),
	}
	ll := bindingoperatorscoreoscomv1alpha1.ServiceResourceMapList{}
	if err := r.List(ctx, &ll, opts...); err != nil {
		return nil, err
	}

	for _, l := range ll.Items {
		if l.Spec.ServiceKindReference.ApiGroup == ag &&
			l.Spec.ServiceKindReference.Kind == k {
			return &l, nil
		}
	}

	return nil, fmt.Errorf("%w: Api Group '%s' and Kind '%s'", ErrServiceResourceMapNotFound, ag, k)
}

func (r *ServiceProxyReconciler) bakeSED(ctx context.Context, sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy, sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap) (*corev1.Secret, error) {
	// TODO

	sec := corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: sp.Namespace,
			Name:      sp.Status.Binding.Name,
		},
		StringData: map[string]string{
			"sec": "all",
		},
	}
	return &sec, nil
}

func (r *ServiceProxyReconciler) createOrUpdateSED(ctx context.Context, sed *corev1.Secret) error {
	// TODO
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

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bindingoperatorscoreoscomv1alpha1.ServiceProxy{}).
		Complete(r)
}
