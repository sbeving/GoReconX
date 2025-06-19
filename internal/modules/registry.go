package modules

import (
	"sync"
)

// ModuleRegistry manages all available reconnaissance modules
type ModuleRegistry struct {
	modules map[string]Module
	mutex   sync.RWMutex
}

// Global registry instance
var GlobalRegistry = NewModuleRegistry()

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
	}
}

// Register registers a module in the registry
func (r *ModuleRegistry) Register(module Module) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	info := module.GetInfo()
	r.modules[info.Name] = module
}

// Get returns a module by name
func (r *ModuleRegistry) Get(name string) (Module, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	module, exists := r.modules[name]
	return module, exists
}

// GetAll returns all registered modules
func (r *ModuleRegistry) GetAll() map[string]Module {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]Module)
	for name, module := range r.modules {
		result[name] = module
	}
	return result
}

// GetByCategory returns modules filtered by category
func (r *ModuleRegistry) GetByCategory(category string) map[string]Module {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]Module)
	for name, module := range r.modules {
		if module.GetInfo().Category == category {
			result[name] = module
		}
	}
	return result
}

// GetCategories returns all available categories
func (r *ModuleRegistry) GetCategories() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	categories := make(map[string]bool)
	for _, module := range r.modules {
		categories[module.GetInfo().Category] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// GetModuleInfo returns information about all modules
func (r *ModuleRegistry) GetModuleInfo() map[string]ModuleInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]ModuleInfo)
	for name, module := range r.modules {
		result[name] = module.GetInfo()
	}
	return result
}

// Count returns the number of registered modules
func (r *ModuleRegistry) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.modules)
}

// init registers all available modules
func init() {
	// Register all modules
	GlobalRegistry.Register(NewDomainEnumModule())
	GlobalRegistry.Register(NewPortScanModule())
	GlobalRegistry.Register(NewWebEnumModule())
	GlobalRegistry.Register(NewEmailEnumModule())
	GlobalRegistry.Register(NewNetworkReconModule())

	// Additional modules can be registered here
}
