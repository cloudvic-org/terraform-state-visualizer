package main

import (
	"fmt"
	"strings"
)

// StateData represents the parsed Terraform state data
type StateData struct {
	FormatVersion    string         `json:"format_version"`
	TerraformVersion string         `json:"terraform_version"`
	Values           StateValues    `json:"values"`
	Resources        []Resource     `json:"-"`
	Outputs          []Output       `json:"-"`
	ResourceCounts   map[string]int `json:"-"`
	RootModule       RootModule     `json:"-"`
}

// StateValues represents the values section of the state
type StateValues struct {
	Outputs    map[string]OutputValue `json:"outputs"`
	RootModule RootModule             `json:"root_module"`
}

// OutputValue represents a single output value
type OutputValue struct {
	Sensitive bool        `json:"sensitive"`
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
}

// RootModule represents the root module in the state
type RootModule struct {
	Resources    []Resource             `json:"resources"`
	Outputs      map[string]OutputValue `json:"outputs"`
	ChildModules []Module               `json:"child_modules"`
}

// Module represents a Terraform module in the state
type Module struct {
	Address      string                 `json:"address"`
	Resources    []Resource             `json:"resources"`
	Outputs      map[string]OutputValue `json:"outputs"`
	ChildModules []Module               `json:"child_modules"`
}

// Resource represents a Terraform resource in the state
type Resource struct {
	Address         string                 `json:"address"`
	Mode            string                 `json:"mode"`
	Type            string                 `json:"type"`
	Name            string                 `json:"name"`
	ProviderName    string                 `json:"provider_name"`
	SchemaVersion   int                    `json:"schema_version"`
	Values          map[string]interface{} `json:"values"`
	SensitiveValues map[string]interface{} `json:"sensitive_values"`
	DependsOn       []string               `json:"depends_on"`
}

// Output represents a parsed output
type Output struct {
	Name      string      `json:"name"`
	Sensitive bool        `json:"sensitive"`
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
}

// parseStateData parses the raw state data into our structured format
func parseStateData(rawData interface{}) (*StateData, error) {
	stateMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid state data format")
	}

	state := &StateData{
		ResourceCounts: make(map[string]int),
	}

	// Parse format version
	if formatVersion, ok := stateMap["format_version"].(string); ok {
		state.FormatVersion = formatVersion
	}

	// Parse terraform version
	if terraformVersion, ok := stateMap["terraform_version"].(string); ok {
		state.TerraformVersion = terraformVersion
	}

	// Parse values section
	if valuesData, ok := stateMap["values"].(map[string]interface{}); ok {
		if err := parseValues(valuesData, state); err != nil {
			return nil, fmt.Errorf("parsing values: %v", err)
		}
	}

	return state, nil
}

// parseValues parses the values section of the state
func parseValues(valuesData map[string]interface{}, state *StateData) error {
	// Parse outputs
	if outputsData, ok := valuesData["outputs"].(map[string]interface{}); ok {
		for name, outputData := range outputsData {
			if outputMap, ok := outputData.(map[string]interface{}); ok {
				output := Output{
					Name: name,
				}

				if sensitive, ok := outputMap["sensitive"].(bool); ok {
					output.Sensitive = sensitive
				}

				if outputType, ok := outputMap["type"]; ok {
					output.Type = outputType
				}

				if value, ok := outputMap["value"]; ok {
					output.Value = value
				}

				state.Outputs = append(state.Outputs, output)
			}
		}
	}

	// Parse root module
	if rootModuleData, ok := valuesData["root_module"].(map[string]interface{}); ok {
		if err := parseRootModule(rootModuleData, state); err != nil {
			return fmt.Errorf("parsing root module: %v", err)
		}
	}

	return nil
}

// parseRootModule parses the root module section
func parseRootModule(rootModuleData map[string]interface{}, state *StateData) error {
	// Parse resources
	if resourcesData, ok := rootModuleData["resources"].([]interface{}); ok {
		for _, resourceData := range resourcesData {
			if resourceMap, ok := resourceData.(map[string]interface{}); ok {
				resource := Resource{}

				if address, ok := resourceMap["address"].(string); ok {
					resource.Address = address
				}

				if mode, ok := resourceMap["mode"].(string); ok {
					resource.Mode = mode
				}

				if resourceType, ok := resourceMap["type"].(string); ok {
					resource.Type = resourceType
				}

				if name, ok := resourceMap["name"].(string); ok {
					resource.Name = name
				}

				if providerName, ok := resourceMap["provider_name"].(string); ok {
					resource.ProviderName = providerName
				}

				if schemaVersion, ok := resourceMap["schema_version"].(float64); ok {
					resource.SchemaVersion = int(schemaVersion)
				}

				if values, ok := resourceMap["values"].(map[string]interface{}); ok {
					resource.Values = values
				}

				if sensitiveValues, ok := resourceMap["sensitive_values"].(map[string]interface{}); ok {
					resource.SensitiveValues = sensitiveValues
				}

				if dependsOn, ok := resourceMap["depends_on"].([]interface{}); ok {
					for _, dep := range dependsOn {
						if depStr, ok := dep.(string); ok {
							resource.DependsOn = append(resource.DependsOn, depStr)
						}
					}
				}

				state.Resources = append(state.Resources, resource)

				// Count resources by type
				resourceTypeKey := resource.Type
				if resource.Mode == "data" {
					resourceTypeKey = "data." + resource.Type
				}
				state.ResourceCounts[resourceTypeKey]++
			}
		}
	}

	// Parse child modules
	if childModulesData, ok := rootModuleData["child_modules"].([]interface{}); ok {
		modules, err := parseModules(childModulesData, "")
		if err != nil {
			return fmt.Errorf("parsing child modules: %v", err)
		}
		state.RootModule.ChildModules = modules

		// Add all module resources to the main resources list
		for _, module := range modules {
			addModuleResourcesToState(module, state)
		}
	}

	return nil
}

// formatResourceMode returns a formatted string for the resource mode
func formatResourceMode(mode string) string {
	switch mode {
	case "managed":
		return "Managed"
	case "data":
		return "Data Source"
	default:
		return strings.Title(mode)
	}
}

// isSensitiveValue checks if a value should be masked as sensitive
func isSensitiveValue(key string, value interface{}, sensitiveValues map[string]interface{}) bool {
	// Check if the key is explicitly marked as sensitive
	if sensitiveValues != nil {
		if _, exists := sensitiveValues[key]; exists {
			return true
		}
	}

	// Check for common sensitive field names
	sensitiveKeys := []string{
		"password", "secret", "key", "token", "credential",
		"private_key", "public_key", "certificate", "ca_cert",
		"access_key", "secret_key", "api_key", "auth_token",
	}

	keyLower := strings.ToLower(key)
	for _, sensitiveKey := range sensitiveKeys {
		if strings.Contains(keyLower, sensitiveKey) {
			return true
		}
	}

	return false
}

// maskSensitiveValue returns a masked version of a sensitive value
func maskSensitiveValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if len(v) > 8 {
			return v[:4] + "..." + v[len(v)-4:]
		}
		return "***"
	case []interface{}:
		return "[***]"
	case map[string]interface{}:
		return "{***}"
	default:
		return "***"
	}
}

// parseModules recursively parses child modules
func parseModules(modulesData []interface{}, parentAddress string) ([]Module, error) {
	var modules []Module

	for _, moduleData := range modulesData {
		if moduleMap, ok := moduleData.(map[string]interface{}); ok {
			module := Module{}

			// Parse module address
			if address, ok := moduleMap["address"].(string); ok {
				module.Address = address
			}

			// Parse module resources
			if resourcesData, ok := moduleMap["resources"].([]interface{}); ok {
				for _, resourceData := range resourcesData {
					if resourceMap, ok := resourceData.(map[string]interface{}); ok {
						resource := parseResource(resourceMap)
						module.Resources = append(module.Resources, resource)
					}
				}
			}

			// Parse module outputs
			if outputsData, ok := moduleMap["outputs"].(map[string]interface{}); ok {
				module.Outputs = make(map[string]OutputValue)
				for name, outputData := range outputsData {
					if outputMap, ok := outputData.(map[string]interface{}); ok {
						output := OutputValue{}

						if sensitive, ok := outputMap["sensitive"].(bool); ok {
							output.Sensitive = sensitive
						}

						if outputType, ok := outputMap["type"]; ok {
							output.Type = outputType
						}

						if value, ok := outputMap["value"]; ok {
							output.Value = value
						}

						module.Outputs[name] = output
					}
				}
			}

			// Parse nested child modules recursively
			if childModulesData, ok := moduleMap["child_modules"].([]interface{}); ok {
				childModules, err := parseModules(childModulesData, module.Address)
				if err != nil {
					return nil, fmt.Errorf("parsing child modules for %s: %v", module.Address, err)
				}
				module.ChildModules = childModules
			}

			modules = append(modules, module)
		}
	}

	return modules, nil
}

// parseResource parses a single resource from JSON data
func parseResource(resourceMap map[string]interface{}) Resource {
	resource := Resource{}

	if address, ok := resourceMap["address"].(string); ok {
		resource.Address = address
	}

	if mode, ok := resourceMap["mode"].(string); ok {
		resource.Mode = mode
	}

	if resourceType, ok := resourceMap["type"].(string); ok {
		resource.Type = resourceType
	}

	if name, ok := resourceMap["name"].(string); ok {
		resource.Name = name
	}

	if providerName, ok := resourceMap["provider_name"].(string); ok {
		resource.ProviderName = providerName
	}

	if schemaVersion, ok := resourceMap["schema_version"].(float64); ok {
		resource.SchemaVersion = int(schemaVersion)
	}

	if values, ok := resourceMap["values"].(map[string]interface{}); ok {
		resource.Values = values
	}

	if sensitiveValues, ok := resourceMap["sensitive_values"].(map[string]interface{}); ok {
		resource.SensitiveValues = sensitiveValues
	}

	if dependsOn, ok := resourceMap["depends_on"].([]interface{}); ok {
		for _, dep := range dependsOn {
			if depStr, ok := dep.(string); ok {
				resource.DependsOn = append(resource.DependsOn, depStr)
			}
		}
	}

	return resource
}

// addModuleResourcesToState recursively adds all resources from a module and its children to the state
func addModuleResourcesToState(module Module, state *StateData) {
	// Add resources from this module
	for _, resource := range module.Resources {
		state.Resources = append(state.Resources, resource)

		// Count resources by type
		resourceTypeKey := resource.Type
		if resource.Mode == "data" {
			resourceTypeKey = "data." + resource.Type
		}
		state.ResourceCounts[resourceTypeKey]++
	}

	// Recursively add resources from child modules
	for _, childModule := range module.ChildModules {
		addModuleResourcesToState(childModule, state)
	}
}
