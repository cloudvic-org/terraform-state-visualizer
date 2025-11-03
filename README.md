# Terraform State Visualizer

A free, open source Terraform state visualization tool that converts Terraform state JSON files into succinct, interactive HTML pages. Perfect for sharing state information in CI/CD pipelines, pull requests, and team collaboration.

## Quick Start

### Using the Binary

1. **Download the latest release** from the [Releases page](https://github.com/cloudvic-org/terraform-state-visualizer/releases)

2. **Generate a Terraform state JSON file**:
   ```bash
   terraform show -json > state.json
   ```

3. **Generate the visualization**:
   ```bash
   ./terraform-state-visualizer -i state.json -o visualization.html
   ```

4. **Open the HTML file** in your browser to view the interactive visualization

### Using Docker

```bash
# Generate state JSON (as above)
terraform show -json > state.json

# Run with Docker
docker run --rm -v $(pwd):/workspace \
  ghcr.io/cloudvic-org/terraform-state-visualizer:latest \
  -i /workspace/state.json -o /workspace/visualization.html
```

### Using GitHub Action

```yaml
name: Terraform State Visualization
on:
  pull_request:
    paths:
      - 'terraform/**'

jobs:
  state:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: 1.6.0
      
      - name: Terraform State Pull
        run: |
          cd terraform
          terraform init
          terraform state pull > state.json
      
      - name: Generate Visualization
        uses: cloudvic-org/terraform-state-visualizer@v1
        with:
          state-file: terraform/state.json
          output-file: terraform-state-visualization.html
      
      - name: Upload Visualization
        uses: actions/upload-artifact@v4
        with:
          name: terraform-state-visualization
          path: terraform-state-visualization.html
```

## Installation

### From Source

```bash
git clone https://github.com/cloudvic-org/terraform-state-visualizer.git
cd terraform-state-visualizer
go build -o terraform-state-visualizer .
```

### Using Go Install

```bash
go install github.com/cloudvic-org/terraform-state-visualizer@latest
```

## Usage

### Command Line Options

```bash
terraform-state-visualizer [OPTIONS]

Options:
  -i, -input string        Input Terraform state JSON file (required)
  -o, -output string       Output HTML file path (default: state-visualization.html)
  --output-html-path string
                           Output HTML file path (alternative to -o)
  -h, -help               Show help information
  -v, -version            Show version information
```

### Examples

```bash
# Basic usage
terraform-state-visualizer -i state.json

# Custom output file
terraform-state-visualizer -i state.json -o my-state.html

# Using long-form flags
terraform-state-visualizer --input state.json --output-html-path visualization.html
```

## Integration Examples

### GitHub Actions

```yaml
- name: Generate State Visualization
  uses: cloudvic-org/terraform-state-visualizer@v1
  with:
    state-file: terraform/state.json
    output-file: state-visualization.html
    upload-artifact: true
```

### GitLab CI

```yaml
generate_visualization:
  image: ghcr.io/cloudvic-org/terraform-state-visualizer:latest
  script:
    - terraform-state-visualizer -i state.json -o visualization.html
  artifacts:
    paths:
      - visualization.html
```

### Jenkins

```groovy
pipeline {
  agent any
  stages {
    stage('Visualize State') {
      steps {
        sh 'terraform-state-visualizer -i state.json -o visualization.html'
        archiveArtifacts artifacts: 'visualization.html'
      }
    }
  }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [CloudVIC](https://cloudvic.com) - Advanced Terraform plan and state visualizations, with drift detection, git history integration, and more