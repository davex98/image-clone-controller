package controller

import k8serrors "k8s.io/apimachinery/pkg/api/errors"

func doesNotExistError(err error) bool {
	if statusError, ok := err.(*k8serrors.StatusError); ok {
		if statusError.ErrStatus.Code == 404 {
			return true
		}
	}
	return false
}

func hasBeenModifiedError(err error) bool {
	if statusError, ok := err.(*k8serrors.StatusError); ok {
		if statusError.ErrStatus.Code == 409 {
			return true
		}
	}
	return false
}