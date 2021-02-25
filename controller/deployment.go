package controller

import (
	"context"
	"fmt"
	"github.com/davex98/image-clone-controller/repository"
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
	Repository repository.Docker
}

func (d *DeploymentReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	dep := &appsv1.Deployment{}

	if request.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}
	err := d.Get(ctx, request.NamespacedName, dep)
	if err != nil {
		if doesNotExistError(err) {
			d.Log.Info("the object does not exist anymore")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	depSpecCopy := dep.Spec.Template.Spec.DeepCopy()
	for i, container := range dep.Spec.Template.Spec.Containers {
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
		depSpecCopy.Containers[i].Image = newImage.GetName()
	}
	depCopy := dep.DeepCopy()
	depCopy.Spec.Template.Spec = *depSpecCopy
	err = d.Update(ctx, depCopy, &client.UpdateOptions{})
	if err != nil {
		if hasBeenModifiedError(err) {
			d.Log.Info("the object has been modified; please apply your changes to the latest version and try again")
			return reconcile.Result{}, nil
		}
		d.Log.Error(err, "could not update deployment image")
		return reconcile.Result{}, err
	}
	d.Log.Info(fmt.Sprintf("deployment %s has valid images", depCopy.Name))
	return reconcile.Result{}, nil
}


var _ reconcile.Reconciler = &DeploymentReconciler{}


func (d *DaemonsetReconciler) SetUpWithManager(mgr manager.Manager)  error {
	return controllerruntime.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).Complete(d)
}