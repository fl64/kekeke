/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	//"github.com/aquasecurity/trivy-operator/pkg/operator/predicate"
	coordination "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	_ "k8s.io/apimachinery/pkg/fields"
	_ "k8s.io/apimachinery/pkg/labels"

	"k8s.io/apimachinery/pkg/runtime"

	//"k8s.io/kubernetes/pkg/apis/coordination"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// LeaseReconciler reconciles a Lease object
type LeaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Lease object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *LeaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	//log.Info("Lease name", "lease", req.NamespacedName)

	lease := &coordination.Lease{}

	err := r.Get(ctx, req.NamespacedName, lease)
	if err != nil {
		return ctrl.Result{}, err
	}

	nodeName := lease.Spec.HolderIdentity
	renewTime := lease.Spec.RenewTime

	fmt.Println(*nodeName, renewTime.String())

	node := &corev1.Node{}

	r.Client.Get(ctx, client.ObjectKey{Name: *nodeName}, node)

	if node.GetName() == "virtlab-pt-0" {
		log.Info("Node name", "node", node.GetName())
		podList := &corev1.PodList{}

		listOpts := &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(fields.Set{
				"spec.nodeName": node.GetName(),
			}),
		}
		//listOpts = &client.ListOptions{LabelSelector: labels.SelectorFromSet(labels.Set{"app": "monitoring-ping"})}
		//fmt.Printf("%+v", listOpts)

		//fmt.Printf("%+v\n", listOpts)

		err = r.Client.List(ctx, podList, listOpts)

		fmt.Println(err)

		fmt.Println(len(podList.Items))

		// for _, pod := range podList.Items {
		// 	fmt.Println(pod.GetName(), pod.Spec.NodeName)

		// }

	}

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LeaseReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// add spec nodeName to indexer
	if err := mgr.GetFieldIndexer().IndexField(
		context.TODO(),
		&corev1.Pod{},
		"spec.nodeName",
		func(rawObj client.Object) []string {
			pod := rawObj.(*corev1.Pod)
			return []string{pod.Spec.NodeName}
		}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&coordination.Lease{}, builder.WithPredicates(
			predicate.NewPredicateFuncs(func(obj client.Object) bool {
				return obj.GetDeletionTimestamp() == nil
			}),
			predicate.NewPredicateFuncs(func(obj client.Object) bool {
				return obj.GetNamespace() == "kube-node-lease"
			}),
		)).
		Complete(r)
}
