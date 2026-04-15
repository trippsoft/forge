// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/trippsoft/forge/pkg/plugin"
)

// Registry manages a collection of modules.
//
// It allows for registering new modules and looking them up by name.
type Registry struct {
	modules map[string]Module
}

// Register adds a new module to the registry.
//
// If a module with the same name already exists, it will be overwritten.
func (r *Registry) Register(module Module) error {
	if r.modules == nil {
		r.modules = make(map[string]Module)
	}

	if module.ID() == nil {
		return errors.New("module info cannot be nil")
	}

	var name string
	if module.ID().namespace == "forge" && module.ID().pluginName == "core" {
		r.modules[module.ID().moduleName] = module
		return nil
	}

	if module.ID().namespace != "" {
		name = module.ID().namespace + "/"
	}

	if module.ID().pluginName != "" {
		name += module.ID().pluginName + "/"
	}

	name += module.ID().moduleName

	r.modules[name] = module
	return nil
}

// RegisterBuiltinModules registers the built-in core modules to the registry.
func (r *Registry) RegisterBuiltinModules() error {
	var err error
	for _, module := range builtinModules {
		e := r.Register(module)
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// RegisterPluginModules registers modules provided by plugins to the registry.
func (r *Registry) RegisterPluginModules() error {
	var err error
	e := r.registerPluginModulesAtBasePath(plugin.UserPluginBasePath)
	if e != nil {
		err = errors.Join(err, e)
	}

	e = r.registerPluginModulesAtBasePath(plugin.SharedPluginBasePath)
	if e != nil {
		err = errors.Join(err, e)
	}

	return err
}

func (r *Registry) registerPluginModulesAtBasePath(basePath string) error {
	fileInfo, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read base path info for %q: %w", basePath, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("base path %q is not a directory", basePath)
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read directory entries for %q: %w", basePath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directory entries
		}

		e := r.registerPluginModulesAtNamespacePath(basePath, entry.Name())
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func (r *Registry) registerPluginModulesAtNamespacePath(basePath, namespace string) error {
	namespacePath := path.Join(basePath, namespace)
	entries, err := os.ReadDir(namespacePath)
	if err != nil {
		return fmt.Errorf("failed to read namespace directory entries for %q: %w", namespacePath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directory entries
		}

		e := r.registerPluginModulesAtPluginPath(basePath, namespace, entry.Name())
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func (r *Registry) registerPluginModulesAtPluginPath(basePath, namespace, pluginName string) error {
	if namespace == "forge" && pluginName == "discover" {
		// Skip the discover plugin.
		return nil
	}

	path, err := plugin.FindPluginPath(basePath, namespace, pluginName, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "metadata")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe for plugin at '%s': %w", path, err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe for plugin at '%s': %w", path, err)
	}
	defer stderr.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe for plugin at '%s': %w", path, err)
	}
	defer stdin.Close()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)

		for {
			n, readErr := stderr.Read(buf)
			if n > 0 {
				accumulatedStderr.Write(buf[:n])
				text := accumulatedStderr.String()

				if strings.Contains(text, plugin.PluginReadyMessage) {
					close(readyChan)
					return
				}
			}

			if readErr != nil {
				if readErr != io.EOF {
					errChan <- fmt.Errorf("error reading stderr for plugin at '%s': %w", path, readErr)
					stdin.Close()
					cmd.Process.Kill()
					cmd.Wait()
				}
				return
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start plugin at '%s': %w", path, err)
	}
	defer cmd.Wait()

	select {
	case err := <-errChan:
		cmd.Process.Kill()
		return err
	case <-readyChan:
	}

	request := &plugin.MetadataRequest{}

	err = plugin.Write(stdin, request)
	if err != nil {
		return fmt.Errorf("failed to write metadata request to plugin at '%s': %w", path, err)
	}

	response := &plugin.MetadataResponse{}

	err = plugin.Read(stdout, response)
	if err != nil {
		return fmt.Errorf("failed to read metadata response from plugin at '%s': %w", path, err)
	}

	for name, s := range response.Modules {
		id := NewModuleID(namespace, pluginName, name)
		spec, err := s.Spec.ToSpec()
		if err != nil {
			return fmt.Errorf("failed to convert module spec for module %q in plugin at '%s': %w", name, path, err)
		}

		switch s.Type {
		case plugin.ModuleType_LOCAL:
			r.Register(NewLocalPluginModule(basePath, id, spec))
		case plugin.ModuleType_REMOTE:
			r.Register(NewRemotePluginModule(basePath, id, spec))
		default:
			return fmt.Errorf("unknown module type %q for module %q in plugin at '%s'", s.Type, name, path)
		}
	}

	return nil
}

// Lookup retrieves a module by its name from the registry.
func (r *Registry) Lookup(name string) (Module, bool) {
	module, exists := r.modules[name]

	return module, exists
}

// Modules returns a copy of the registered modules in the registry.
//
// This is used mostly for testing purposes to avoid modifying the original map.
func (r *Registry) Modules() map[string]Module {
	if r.modules == nil {
		return make(map[string]Module)
	}

	return maps.Clone(r.modules)
}

// NewRegistry creates a new module registry.
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}
