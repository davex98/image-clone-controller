# image-clone-controller

### Guten tag!

This repo is used as an example how to implement an image-clone-controller!
First you need to have a cluster running locally, eg. minikube, this controller is running in-cluster as a simple deployment.


In order to run, you have to provide the the docker secret file located in ~/.docker/config.json, and then run : **make deploy**.
The controller automatically picks up the deployments and daemonsets, checks whether the images exist in **burghardtkubermatic** registry on dockerhub and if not, automatically makes copy of them, pushes them to that regitsry and replaces the image in particular deployment or daemonset.

[Asciinema proof of the controller working!](https://asciinema.org/a/JbCECII0DU1iX3b94PTMFo3wI)


If you have any questions or suggestions, please let me know :D
