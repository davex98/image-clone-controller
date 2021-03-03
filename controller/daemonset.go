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

	deepCopy := daemon.DeepCopy()
	for i, container := range deepCopy.Spec.Template.Spec.Containers {
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
		deepCopy.Spec.Template.Spec.Containers[i].Image = newImage.GetName()
	}

	err = d.Update(ctx, deepCopy, &client.UpdateOptions{})
	if err != nil {
		if hasBeenModifiedError(err) {
			d.Log.Info("the object has been modified; please apply your changes to the latest version and try again")
			return reconcile.Result{}, nil
		}
		d.Log.Error(err, "could not update daemonset image")
		return reconcile.Result{}, err
	}
	d.Log.Info(fmt.Sprintf("daemonset %s has valid images", deepCopy.Name))
	return reconcile.Result{}, nil

}

var _ reconcile.Reconciler = &DaemonsetReconciler{}
