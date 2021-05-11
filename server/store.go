// Utils for interacting with KVStore
package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

// KVEnsure makes sure the initial value for a key is set to the value provided, if it does not already exists
// Returns the value set for the key in kv-store and error
func (p *Plugin) KVEnsure(key string, newValue []byte) ([]byte, error) {
	value, err := p.KVLoad(key)
	switch err {
	case nil:
		return value, nil
	case ErrNotFound:
		break
	default:
		return nil, err
	}

	err = p.KVStore(key, newValue)
	if err != nil {
		return nil, err
	}

	// Load again in case we lost the race to another server
	value, err = p.KVLoad(key)
	if err != nil {
		return newValue, nil
	}
	return value, nil
}

func (p *Plugin) KVLoad(key string) ([]byte, error) {
	data, appErr := p.API.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}
	if data == nil {
		return nil, ErrNotFound
	}
	return data, nil
}

func (p *Plugin) KVStore(key string, data []byte) error {
	var appErr *model.AppError
	if appErr = p.API.KVSet(key, data); appErr != nil {
		return errors.WithMessagef(appErr, "failed plugin KVSet %q", key)
	}
	return nil
}
