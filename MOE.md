kubebuilder init --domain fencing-controller --repo fencing-controller
kubebuilder create api --group coordination.k8s.io --version v1 --kind Lease --resource=false --controller=true

https://enix.io/en/blog/kubebuilder/

https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/

https://github.com/kubernetes-sigs/controller-runtime/issues/2378

lease controller

```go
// SetupWithManager sets up the controller with the Manager.
func (r *LeaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
       return ctrl.NewControllerManagedBy(mgr).
               // Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
               For(&coordination.Lease{}).
               Complete(r)
}
