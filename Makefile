build_image:
	 docker build . -t image-clone-controller

push_image:
ifdef version
	 docker tag image-clone-controller:latest burghardtkubermatic/image-clone-controller:$(version) &&  docker push burghardtkubermatic/image-clone-controller:$(version)
else
	@echo 'provide version, eg: make push_image version=v5'
endif


deploy_app:
	kubectl apply -f config/namespace.yml
	kubectl apply -f config/serviceaccount.yml
	kubectl apply -f config/configmap.yml
	kubectl apply -f config/clusterrole-deamonset.yml
	kubectl apply -f config/clusterrole-deployment.yml
	kubectl apply -f config/crb-daemon.yml
	kubectl apply -f config/crb-deployment.yml
	kubectl apply -f config/dockersecret.yml
	kubectl apply -f config/image-clone-deployment.yml

test:
	go test ./...
