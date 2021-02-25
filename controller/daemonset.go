package controller

import (
	"context"
	"fmt"
	"github.com/davex98/image-clone-controller/repository"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DaemonsetReconciler struct {
	client.Client
	Log logr.Logger
	Repository repository.Docker
}

func (d *DaemonsetReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	daemon := &appsv1.DaemonSet{}

	if request.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}
	err := d.Get(ctx, request.NamespacedName, daemon)
	if err != nil {
		if doesNotExistError(err) {
			d.Log.Info("the object does not exist anymore")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	daeomonSpecCopy := daemon.Spec.Template.Spec.DeepCopy()
	for i, container := range daemon.Spec.Template.Spec.Containers {
		currentImage := container.Image
		if valid := d.Repository.IsImageValid(currentImage); valid {
			continue
		}
		oldImage, err := d.Repository.PullImage(currentImage)
		if err != nil {
			d.Log.Error(err, "could not pull image")
			return reconcile.Result{}, err
		}

		newImage, err := d.Repository.PushImageToPrivateRepo(oldImage)
		if err != nil {
			d.Log.Error(err, "could not push image to private repo")
			return reconcile.Result{}, err
		}
		daeomonSpecCopy.Containers[i].Image = newImage.GetName()
	}
	daemonCopy := daemon.DeepCopy()
	daemonCopy.Spec.Template.Spec = *daeomonSpecCopy
	err = d.Update(ctx, daemonCopy, &client.UpdateOptions{})
	if err != nil {
		if hasBeenModifiedError(err) {
			d.Log.Info("the object has been modified; please apply your changes to the latest version and try again")
			return reconcile.Result{}, nil
		}
		d.Log.Error(err, "could not update daemonset image")
		return reconcile.Result{}, err
	}
	d.Log.Info(fmt.Sprintf("daemonset %s has valid images", daemonCopy.Name))
	return reconcile.Result{}, nil

}

var _ reconcile.Reconciler = &DaemonsetReconciler{}
