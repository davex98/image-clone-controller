package main

import (
	ctl "github.com/davex98/image-clone-controller/controller"
	"github.com/davex98/image-clone-controller/repository"
	appsv1 "k8s.io/api/apps/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	log := ctrl.Log.WithName("image-clone-controller")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Namespace: ""})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	docker := repository.NewRepository()

	_, err = builder.
		ControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Build(&ctl.DeploymentReconciler{Client: mgr.GetClient(), Log: logf.Log.WithName("deployment-controller"), Repository: docker})

	if err != nil {
		log.Error(err, "could not create deployment controller")
		os.Exit(1)
	}

	_, err = builder.ControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Build(&ctl.DaemonsetReconciler{Client: mgr.GetClient(), Log: logf.Log.WithName("daemonset-controller"), Repository: docker})

	if err != nil {
		log.Error(err, "could not create daemonset controller")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}