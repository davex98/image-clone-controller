package controller

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DeploymentReconciler struct {
	client.Client
	Log logr.Logger
}

func (d *DeploymentReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	dep := &appsv1.Deployment{}

	if request.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}
	err := d.Get(ctx, request.NamespacedName, dep)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

var _ reconcile.Reconciler = &DeploymentReconciler{}


func (d *DaemonsetReconciler) SetUpWithManager(mgr manager.Manager)  error {
	return controllerruntime.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).Complete(d)
}