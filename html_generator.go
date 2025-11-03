package main

import (
	"fmt"
	"strings"
)

// generateHtml creates the complete HTML visualization for Terraform state
func generateHtml(stateData *StateData) string {
	// Generate HTML content
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Terraform State</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .source-link {
            font-size: 14px;
            font-weight: normal;
            color: #3498db;
            text-decoration: none;
        }
        .source-link:hover {
            text-decoration: underline;
        }
        .promo-message {
            text-align: center;
            margin: 20px 0;
            font-style: italic;
        }
        .promo-link {
            color: #3498db;
            text-decoration: none;
            font-weight: bold;
            font-style: normal;
        }
        .promo-link:hover {
            text-decoration: underline;
        }
        .section-header-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            width: 100%;
        }
        .section-description {
            font-size: 14px;
            font-style: italic;
            color: #6c757d;
            margin-bottom: 15px;
        }
        .section {
            margin: 20px 0;
            padding: 15px;
            background-color: #ecf0f1;
            border-radius: 5px;
        }
        .resource-item {
            margin: 10px 0;
            padding: 10px;
            background-color: white;
            border-radius: 3px;
            border-left: 4px solid #3498db;
        }
        .managed { border-left-color: #27ae60; }
        .data { border-left-color: #f39c12; }
        .resource-address {
            font-family: monospace;
            font-weight: bold;
            color: #2c3e50;
        }
        .resource-type {
            color: #7f8c8d;
            font-size: 14px;
        }
        .resource-attributes {
            margin-top: 10px;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 3px;
            font-family: monospace;
            font-size: 12px;
        }
        .collapsible {
            cursor: pointer;
            user-select: none;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .collapsible:hover {
            background-color: #f0f0f0;
        }
        .collapsible::before {
            content: "▼";
            font-size: 12px;
            transition: transform 0.2s;
            flex-shrink: 0;
        }
        .collapsible.collapsed::before {
            content: "▶";
        }
        .collapsible-content {
            overflow: hidden;
            transition: opacity 0.3s ease-out, max-height 0.3s ease-out;
        }
        .collapsible-content.collapsed {
            max-height: 0;
            opacity: 0;
        }
        .collapsible-content:not(.collapsed) {
            max-height: none;
            opacity: 1;
        }
        .attribute-item {
            margin: 5px 0;
            padding: 3px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .attribute-key {
            font-weight: bold;
            color: #495057;
        }
        .attribute-value {
            color: #6c757d;
            margin-left: 10px;
        }
        .attribute-sensitive {
            background-color: #fff3cd;
            border-left: 3px solid #ffc107;
            padding-left: 8px;
        }
        .summary {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
        }
        .summary-item {
            flex: 1;
            text-align: center;
            padding: 15px;
            background-color: white;
            border-radius: 5px;
        }
        .summary-number {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
        }
        .summary-label {
            color: #7f8c8d;
            font-size: 14px;
        }
        .module-item {
            margin: 10px 0;
            padding: 15px;
            background-color: white;
            border-radius: 5px;
            border-left: 4px solid #9b59b6;
        }
        .module-address {
            font-family: monospace;
            font-weight: bold;
            color: #8e44ad;
            font-size: 16px;
        }
        .module-resource-count {
            color: #7f8c8d;
            font-size: 14px;
            margin-top: 5px;
        }
    </style>
    <script>
        function toggleCollapsible(element) {
            const content = element.nextElementSibling;
            element.classList.toggle('collapsed');
            content.classList.toggle('collapsed');
        }
        
        // Make individual resource items collapsed by default, but keep main sections open
        document.addEventListener('DOMContentLoaded', function() {
            const collapsibles = document.querySelectorAll('.collapsible');
            collapsibles.forEach(function(element) {
                // Check if this is a main section (State Overview, Resources, etc.)
                const isMainSection = element.querySelector('h2') !== null;
                
                if (!isMainSection) {
                    // Only collapse individual resource items, not main sections
                    element.classList.add('collapsed');
                    const content = element.nextElementSibling;
                    if (content) {
                        content.classList.add('collapsed');
                    }
                } else {
                    // Check if main section has no items
                    const section = element.closest('.section');
                    const resourceItems = section.querySelectorAll('.resource-item, .module-item');
                    if (resourceItems.length === 0) {
                        element.classList.add('collapsed');
                        const content = element.nextElementSibling;
                        if (content) {
                            content.classList.add('collapsed');
                        }
                    }
                }
            });
        });
    </script>
</head>
<body>
    <div class="container">
        <h1>Terraform State</h1>
                
        <div class="section">
            <div class="collapsible" onclick="toggleCollapsible(this)">
                <div class="section-header-row">
                    <h2>State Overview</h2>
                    <p class="section-description">Summary of your Terraform state</p>
                </div>
            </div>
            <div class="collapsible-content">
                ` + generateStateOverviewHtml(stateData) + `
            </div>
        </div>
        
        <div class="section">
            <div class="collapsible" onclick="toggleCollapsible(this)">
            <div class="section-header-row">
                    <h2>Resources (` + fmt.Sprintf("%d", len(stateData.Resources)) + ` total)</h2>
                    <p class="section-description">All resources in your Terraform state</p>
                </div>
            </div>
            <div class="collapsible-content">
                ` + generateResourcesHtml(stateData) + `
            </div>
        </div>
        
        <div class="section">
            <div class="collapsible" onclick="toggleCollapsible(this)">
            <div class="section-header-row">
                    <h2>Outputs (` + fmt.Sprintf("%d", len(stateData.Outputs)) + ` total)</h2>
                    <p class="section-description">Output values from your Terraform state</p>
                </div>
            </div>
            <div class="collapsible-content">
                ` + generateOutputsHtml(stateData) + `
            </div>
        </div>
        
        <div class="section">
            <div class="collapsible" onclick="toggleCollapsible(this)">
            <div class="section-header-row">
                    <h2>Modules (` + fmt.Sprintf("%d", len(stateData.RootModule.ChildModules)) + ` total)</h2>
                    <p class="section-description">Module hierarchy and organization</p>
                </div>
            </div>
            <div class="collapsible-content">
                ` + generateModulesHtml(stateData) + `
            </div>
        </div>
    </div>
    <div class="promo-message">
        Want to visualize your Terraform plan and state changes over time and link them to your git history?<br>
        <a href="https://cloudvic.com" class="promo-link">Try CloudVIC</a>
    </div>
</body>
</html>`

	return html
}

// generateStateOverviewHtml creates the state overview section
func generateStateOverviewHtml(stateData *StateData) string {
	var html strings.Builder

	html.WriteString(`<div class="summary">
		<div class="summary-item">
			<div class="summary-number">` + fmt.Sprintf("%d", len(stateData.Resources)) + `</div>
			<div class="summary-label">Resources</div>
		</div>
		<div class="summary-item">
			<div class="summary-number">` + fmt.Sprintf("%d", len(stateData.Outputs)) + `</div>
			<div class="summary-label">Outputs</div>
		</div>
		<div class="summary-item">
			<div class="summary-number">` + stateData.FormatVersion + `</div>
			<div class="summary-label">Format Version</div>
		</div>
		<div class="summary-item">
			<div class="summary-number">` + stateData.TerraformVersion + `</div>
			<div class="summary-label">Terraform Version</div>
		</div>
	</div>`)

	// Add resource type breakdown
	if len(stateData.ResourceCounts) > 0 {
		html.WriteString(`<div style="margin-top: 20px;">
			<h3>Resources by Type</h3>
			<div style="margin-top: 10px;">`)

		// Sort resource types alphabetically
		var sortedTypes []string
		for resourceType := range stateData.ResourceCounts {
			sortedTypes = append(sortedTypes, resourceType)
		}

		// Simple alphabetical sort
		for i := 0; i < len(sortedTypes)-1; i++ {
			for j := i + 1; j < len(sortedTypes); j++ {
				if sortedTypes[i] > sortedTypes[j] {
					sortedTypes[i], sortedTypes[j] = sortedTypes[j], sortedTypes[i]
				}
			}
		}

		for _, resourceType := range sortedTypes {
			count := stateData.ResourceCounts[resourceType]
			html.WriteString(fmt.Sprintf(`
				<div style="padding: 5px 0; border-bottom: 1px solid #eee;">
					<span style="font-weight: bold; color: #2c3e50;">%s:</span>
					<span style="color: #3498db; margin-left: 10px;">%d</span>
				</div>`, resourceType, count))
		}

		html.WriteString(`</div></div>`)
	}

	return html.String()
}

// generateResourcesHtml creates the resources section
func generateResourcesHtml(stateData *StateData) string {
	if len(stateData.Resources) == 0 {
		return "<p>No resources found in state.</p>"
	}

	var html strings.Builder
	html.WriteString("<div>")

	for _, resource := range stateData.Resources {
		modeClass := "managed"
		if resource.Mode == "data" {
			modeClass = "data"
		}

		html.WriteString(fmt.Sprintf(`
			<div class="resource-item %s">
				<div class="collapsible" onclick="toggleCollapsible(this)">
					<div>%s</div>
					<div class="resource-address">%s</div>
				</div>
				<div class="collapsible-content">
					<div class="resource-attributes">
						%s
					</div>
				</div>
			</div>`,
			modeClass,
			formatResourceMode(resource.Mode),
			resource.Address,
			generateResourceAttributesHtml(resource)))
	}

	html.WriteString("</div>")
	return html.String()
}

// generateResourceAttributesHtml creates the attributes section for a resource
func generateResourceAttributesHtml(resource Resource) string {
	var html strings.Builder

	// Basic resource info
	html.WriteString(fmt.Sprintf(`
		<div class="attribute-item">
			<span class="attribute-key">Type:</span>
			<span class="attribute-value">%s</span>
		</div>`, resource.Type))

	html.WriteString(fmt.Sprintf(`
		<div class="attribute-item">
			<span class="attribute-key">Provider:</span>
			<span class="attribute-value">%s</span>
		</div>`, resource.ProviderName))

	html.WriteString(fmt.Sprintf(`
		<div class="attribute-item">
			<span class="attribute-key">Schema Version:</span>
			<span class="attribute-value">%d</span>
		</div>`, resource.SchemaVersion))

	// Resource values
	if len(resource.Values) > 0 {
		html.WriteString(`<div class="attribute-item">
			<span class="attribute-key">Configuration:</span>
		</div>`)

		for key, value := range resource.Values {
			cssClass := ""
			if isSensitiveValue(key, value, resource.SensitiveValues) {
				cssClass = "attribute-sensitive"
			}

			valueStr := formatValue(value)
			if isSensitiveValue(key, value, resource.SensitiveValues) {
				valueStr = maskSensitiveValue(value)
			}

			html.WriteString(fmt.Sprintf(`
				<div class="attribute-item %s">
					<span class="attribute-key">%s:</span>
					<span class="attribute-value">%s</span>
				</div>`, cssClass, key, valueStr))
		}
	}

	// Dependencies
	if len(resource.DependsOn) > 0 {
		html.WriteString(`<div class="attribute-item">
			<span class="attribute-key">Dependencies:</span>
		</div>`)

		for _, dep := range resource.DependsOn {
			html.WriteString(fmt.Sprintf(`
				<div class="attribute-item">
					<span class="attribute-value">%s</span>
				</div>`, dep))
		}
	}

	return html.String()
}

// generateOutputsHtml creates the outputs section
func generateOutputsHtml(stateData *StateData) string {
	if len(stateData.Outputs) == 0 {
		return "<p>No outputs found in state.</p>"
	}

	var html strings.Builder
	html.WriteString("<div>")

	for _, output := range stateData.Outputs {
		cssClass := ""
		if output.Sensitive {
			cssClass = "attribute-sensitive"
		}

		valueStr := formatValue(output.Value)
		if output.Sensitive {
			valueStr = maskSensitiveValue(output.Value)
		}

		html.WriteString(fmt.Sprintf(`
			<div class="resource-item %s">
				<div class="collapsible" onclick="toggleCollapsible(this)">
					<div>Output</div>
					<div class="resource-address">%s</div>
				</div>
				<div class="collapsible-content">
					<div class="resource-attributes">
						<div class="attribute-item">
							<span class="attribute-key">Type:</span>
							<span class="attribute-value">%s</span>
						</div>
						<div class="attribute-item %s">
							<span class="attribute-key">Value:</span>
							<span class="attribute-value">%s</span>
						</div>
						<div class="attribute-item">
							<span class="attribute-key">Sensitive:</span>
							<span class="attribute-value">%t</span>
						</div>
					</div>
				</div>
			</div>`,
			cssClass,
			output.Name,
			formatValue(output.Type),
			cssClass,
			valueStr,
			output.Sensitive))
	}

	html.WriteString("</div>")
	return html.String()
}

// generateModulesHtml creates the modules section
func generateModulesHtml(stateData *StateData) string {
	if len(stateData.RootModule.ChildModules) == 0 {
		return `<div style="text-align: center; padding: 40px 20px; background-color: #f8f9fa; border-radius: 8px; margin: 20px 0;">
			<h3 style="color: #2c3e50; margin-bottom: 15px;">No modules found</h3>
			<p style="color: #6c757d; margin-bottom: 20px; font-size: 16px;">
				This state file does not contain any child modules.
			</p>
		</div>`
	}

	var html strings.Builder
	html.WriteString("<div>")

	// Generate module hierarchy
	for _, module := range stateData.RootModule.ChildModules {
		html.WriteString(generateModuleHtml(module, 0))
	}

	html.WriteString("</div>")
	return html.String()
}

// generateModuleHtml creates HTML for a single module and its children
func generateModuleHtml(module Module, depth int) string {
	var html strings.Builder

	// Calculate margin based on depth
	marginLeft := depth * 20

	// Count total resources in this module and its children
	totalResources := countModuleResources(module)

	html.WriteString(fmt.Sprintf(`
		<div class="module-item" style="margin-left: %dpx;">
			<div class="collapsible" onclick="toggleCollapsible(this)">
				<div>
					<div class="module-address">%s</div>
					<div class="module-resource-count">%d resources</div>
				</div>
			</div>
			<div class="collapsible-content">
				<div class="resource-attributes">`, marginLeft, module.Address, totalResources))

	// Add resources from this module
	if len(module.Resources) > 0 {
		html.WriteString(`<div class="attribute-item">
			<span class="attribute-key">Resources:</span>
		</div>`)

		for _, resource := range module.Resources {
			modeClass := "managed"
			if resource.Mode == "data" {
				modeClass = "data"
			}

			html.WriteString(fmt.Sprintf(`
				<div class="resource-item %s" style="margin-left: 20px;">
					<div class="collapsible" onclick="toggleCollapsible(this)">
						<div>%s</div>
						<div class="resource-address">%s</div>
					</div>
					<div class="collapsible-content">
						<div class="resource-attributes">
							%s
						</div>
					</div>
				</div>`,
				modeClass,
				formatResourceMode(resource.Mode),
				resource.Address,
				generateResourceAttributesHtml(resource)))
		}
	}

	// Add outputs from this module
	if len(module.Outputs) > 0 {
		html.WriteString(`<div class="attribute-item">
			<span class="attribute-key">Outputs:</span>
		</div>`)

		for name, output := range module.Outputs {
			cssClass := ""
			if output.Sensitive {
				cssClass = "attribute-sensitive"
			}

			valueStr := formatValue(output.Value)
			if output.Sensitive {
				valueStr = maskSensitiveValue(output.Value)
			}

			html.WriteString(fmt.Sprintf(`
				<div class="attribute-item %s" style="margin-left: 20px;">
					<span class="attribute-key">%s:</span>
					<span class="attribute-value">%s</span>
				</div>`, cssClass, name, valueStr))
		}
	}

	html.WriteString(`</div>
			</div>
		</div>`)

	// Add child modules recursively
	for _, childModule := range module.ChildModules {
		html.WriteString(generateModuleHtml(childModule, depth+1))
	}

	return html.String()
}

// countModuleResources recursively counts resources in a module and its children
func countModuleResources(module Module) int {
	count := len(module.Resources)

	for _, childModule := range module.ChildModules {
		count += countModuleResources(childModule)
	}

	return count
}

// formatValue formats a value for display in HTML
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if len(v) > 100 {
			return v[:100] + "..."
		}
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}
		return fmt.Sprintf("[%d items]", len(v))
	case map[string]interface{}:
		return fmt.Sprintf("{%d fields}", len(v))
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
