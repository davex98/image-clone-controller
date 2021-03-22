# image-clone-controller

# Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [How to](#how-to)
    * [Run the controller](#run-the-controller)
    * [Run the unit tests](#run-the-unit-tests)


# Overview


It's an implementation of a controller which watches the applications and “caches” the images by re-uploading to our
own registry repository and reconfiguring the applications to use these copies.
# Project Structure

- [Dockerfile](Dockerfile) - Dockerfile, used for production to run within the cluster (running the application as nonroot) 
- [controller](controller/) - controller code, it is written using [kubebuilder](https://book.kubebuilder.io/)
- [repository](repository/) - code that enables connecting to the repository on which the images will be re-uploaded
- [config](config/) - yml files that enable running the controller within the cluster

# Prerequisites
Make sure you have installed all of the following prerequisites on your development machine:
* Git - [Download & Install Git](https://git-scm.com/downloads). OSX and Linux machines typically have this already installed.
* Minikube - [Download & Install Minikube](https://minikube.sigs.k8s.io/docs/start/)
* Docker - [Download & Install Docker](https://docs.docker.com/engine/install/ubuntu/). Docker is used for building images and running the end-to-end tests.
* Golang - [Download & Install Golang](https://golang.org/doc/install).

# How to

## Run the controller
In order to run, you have to provide the dockersecret.yml file located in ~/.docker/config.json
Then to run the controller you can simply use:

```bash
minikube start
make deploy_app
```
It will automatically deploy the controller on local cluster eg. **Minikube**.

## Run the unit tests
In order to run the unit tests, run:

```bash
minikube start
make test
```

[Asciinema proof of the controller working!](https://asciinema.org/a/JbCECII0DU1iX3b94PTMFo3wI)

