package repository

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"os"
	"strings"
)

type Docker interface {
	PullImage(imageName string) (image, error)
	PushImageToPrivateRepo(imageName image)	(image, error)
	IsImageValid(imageName string) bool
}

type repository struct {
	repoName string
}

func (r repository) IsImageValid(imageName string) bool {
	if !strings.HasPrefix(imageName, r.repoName) {
		return false
	}
	return true
}

func (r repository) PullImage(imageName string) (image, error) {
	reference, err := name.ParseReference(imageName)
	if err != nil {
		return image{}, err
	}
	img, err := remote.Get(reference, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return image{}, err
	}
	i, err := img.Image()
	if err != nil {
		return image{}, err
	}
	return image{
		Image: i,
		name:  reference.String(),
	}, nil
}

func (r repository) PushImageToPrivateRepo(newImage image) (image, error) {
	path := fmt.Sprintf("%s/%s", r.repoName, newImage.name)
	ref, err := name.ParseReference(path)
	if err != nil {
		return image{}, err
	}
	err = remote.Write(ref, newImage, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return image{}, err
	}
	return image{
		Image: nil,
		name:  ref.String(),
	}, nil
}

func NewRepository() Docker {
	getenv := os.Getenv("DOCKER_REPO")
	if getenv == "" {
		panic("pass DOCKER_REPO variable")
	}

	return repository{repoName: getenv}
}
