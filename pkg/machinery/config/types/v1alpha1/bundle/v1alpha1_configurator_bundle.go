// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bundle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/configpatcher"
	"github.com/talos-systems/talos/pkg/machinery/config/encoder"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
)

// ConfigBundle defines the group of v1alpha1 config files.
// docgen: nodoc
// +k8s:deepcopy-gen=false
type ConfigBundle struct {
	InitCfg         *v1alpha1.Config
	ControlPlaneCfg *v1alpha1.Config
	WorkerCfg       *v1alpha1.Config
	TalosCfg        *clientconfig.Config
}

// Init implements the ProviderBundle interface.
func (c *ConfigBundle) Init() config.Provider {
	return c.InitCfg
}

// ControlPlane implements the ProviderBundle interface.
func (c *ConfigBundle) ControlPlane() config.Provider {
	return c.ControlPlaneCfg
}

// Worker implements the ProviderBundle interface.
func (c *ConfigBundle) Worker() config.Provider {
	return c.WorkerCfg
}

// TalosConfig implements the ProviderBundle interface.
func (c *ConfigBundle) TalosConfig() *clientconfig.Config {
	return c.TalosCfg
}

// Write config files to output directory.
func (c *ConfigBundle) Write(outputDir string, commentsFlags encoder.CommentsFlags, types ...machine.Type) error {
	for _, t := range types {
		name := strings.ToLower(t.String()) + ".yaml"
		fullFilePath := filepath.Join(outputDir, name)

		var (
			configString string
			err          error
		)

		switch t {
		case machine.TypeInit:
			configString, err = c.Init().EncodeString(encoder.WithComments(commentsFlags))
			if err != nil {
				return err
			}
		case machine.TypeControlPlane:
			configString, err = c.ControlPlane().EncodeString(encoder.WithComments(commentsFlags))
			if err != nil {
				return err
			}
		case machine.TypeWorker:
			configString, err = c.Worker().EncodeString(encoder.WithComments(commentsFlags))
			if err != nil {
				return err
			}
		case machine.TypeUnknown:
			fallthrough
		default:
			return fmt.Errorf("unexpected machine type %v", t)
		}

		if err = os.WriteFile(fullFilePath, []byte(configString), 0o644); err != nil {
			return err
		}

		fmt.Printf("created %s\n", fullFilePath)
	}

	return nil
}

// ApplyPatches patches every config type with a patch.
func (c *ConfigBundle) ApplyPatches(patches []configpatcher.Patch, patchControlPlane, patchWorker bool) error {
	if len(patches) == 0 {
		return nil
	}

	apply := func(in *v1alpha1.Config) (*v1alpha1.Config, error) {
		patched, err := configpatcher.Apply(configpatcher.WithConfig(in), patches)
		if err != nil {
			return nil, err
		}

		cfg, err := patched.Config()
		if err != nil {
			return nil, err
		}

		out, ok := cfg.Raw().(*v1alpha1.Config)
		if !ok {
			return nil, fmt.Errorf("unexpected config type %T", cfg.Raw())
		}

		return out, nil
	}

	var err error

	if patchControlPlane {
		c.InitCfg, err = apply(c.InitCfg)
		if err != nil {
			return err
		}

		c.ControlPlaneCfg, err = apply(c.ControlPlaneCfg)
		if err != nil {
			return err
		}
	}

	if patchWorker {
		c.WorkerCfg, err = apply(c.WorkerCfg)
		if err != nil {
			return err
		}
	}

	return nil
}
