package storageconfigloader

import (
	"encoding/json"
	"fmt"

	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(new(StorageLoader))
}

const defaultKey = "config/caddy.json"

// StorageLoader is a dynamic configuration loader that reads the configuration from a Caddy storage. If
// the storage is not configured, the default storage is used, which may be the file-system if none is configured
// If the `key` is not configured, the default key is `config/caddy.json`.
type StorageLoader struct {
	// StorageRaw is a storage module that defines how/where Caddy
	// stores assets (such as TLS certificates). The default storage
	// module is `caddy.storage.file_system` (the local file system),
	// and the default path
	// [depends on the OS and environment](/docs/conventions#data-directory).
	StorageRaw json.RawMessage `json:"storage,omitempty" caddy:"namespace=caddy.storage inline_key=module"`

	// The storage key at which the configuration is to be found
	Key string `json:"key,omitempty"`

	storage certmagic.Storage
	logger  *zap.Logger
}

// CaddyModule implements caddy.Module.
func (*StorageLoader) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.config_loaders.storage",
		New: func() caddy.Module { return new(StorageLoader) },
	}
}

// Provision implements caddy.Provisioner.
func (sl *StorageLoader) Provision(ctx caddy.Context) error {
	sl.logger = ctx.Logger()
	if sl.Key == "" {
		sl.Key = defaultKey
	}
	if sl.StorageRaw != nil {
		val, err := ctx.LoadModule(sl, "StorageRaw")
		if err != nil {
			return fmt.Errorf("loading storage module: %v", err)
		}
		cmStorage, err := val.(caddy.StorageConverter).CertMagicStorage()
		if err != nil {
			return fmt.Errorf("creating storage configuration: %v", err)
		}
		sl.storage = cmStorage
	}
	if sl.storage == nil {
		sl.storage = ctx.Storage()
	}
	return nil
}

// LoadConfig reads the configuration from the storage
func (sl *StorageLoader) LoadConfig(ctx caddy.Context) ([]byte, error) {
	sl.logger.Info("loading config from storage", zap.String("key", sl.Key))
	return sl.storage.Load(ctx, sl.Key)
}

var _ caddy.Module = (*StorageLoader)(nil)
var _ caddy.Provisioner = (*StorageLoader)(nil)
var _ caddy.ConfigLoader = (*StorageLoader)(nil)
