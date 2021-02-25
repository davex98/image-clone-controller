package repository

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type image struct {
	v1.Image
	name string
}

func (i image) GetName() string {
	return i.name
}