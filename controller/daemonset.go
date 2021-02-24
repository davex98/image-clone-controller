package controller

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DaemonsetReconciler struct {
	client.Client
	Log logr.Logger
}

func (d *DaemonsetReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	daemon := &appsv1.DaemonSet{}

	if request.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}
	err := d.Get(ctx, request.NamespacedName, daemon)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

var _ reconcile.Reconciler = &DaemonsetReconciler{}
