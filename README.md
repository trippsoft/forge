# Forge - Configuration Management Tool

Forge is a configuration management tool designed for managing and automating tasks across local and remote systems. It uses HCL (HashiCorp Configuration Language) for defining inventories and workflows.

## Features

- **Performance Focused**: Designed for speed by limiting network I/O, parallelizing tasks, and running compiled code instead of interpreted scripts
- **Security Without Sacrificing Usability**: Emphasizes secure defaults while maintaining ease of use
- **Agent-less**: No need to install agents on target hosts
- **HCL-based Configuration**: Define your infrastructure and automation workflows using HCL for readability
- **Flexible Inventory Management**: Organize hosts into groups with hierarchical support and variable inheritance
- **Workflow Automation**: Execute multi-step processes across multiple targets with conditional logic
- **Plugin Architecture**: Extensible module system for custom functionality
- **Multi-transport Support**: Connect to hosts via SSH or local execution
- **Privilege Escalation**: Built-in support for privilege escalation and user impersonation (currently Linux/Unix only)
- **Discovery Plugins**: Automatically discover host information (OS, package manager, services, etc.)
- **Cross-Platform Support**: Run on Linux, macOS, or Windows and manage systems of a variety of OS and architecture

## Getting Started

### Installing Forge

Pre-built binaries will be made available when the project reaches a more stable state. For now, you can build Forge from source.

### Building Forge

#### Prerequisites
- Go (see go.mod for minimum version)
- Git

#### Linux / macOS

Create the directory structure.

```bash
sudo mkdir -p /usr/share/forge/plugins/forge/core
sudo mkdir -p /usr/share/forge/plugins/forge/discover
chmod -R 755 /usr/share/forge
```

Clone the repository and build the CLI.

```bash
git clone https://github.com/trippsoft/forge.git
cd forge
go build -o forge ./cmd/forge
```

Install the binary.

```bash
sudo mv forge /usr/share/forge/
sudo chmod 755 /usr/share/forge/forge
sudo ln -s /usr/share/forge/forge /usr/local/bin/forge
```

Build the core plugin server(s).  Skip the platforms you don't intend to manage or run the tool from.

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o forge-core_darwin_amd64 ./cmd/forge-core
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o forge-core_darwin_arm64 ./cmd/forge-core
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forge-core_linux_amd64 ./cmd/forge-core
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o forge-core_linux_arm64 ./cmd/forge-core
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o forge-core_windows_amd64.exe ./cmd/forge-core
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o forge-core_windows_arm64.exe ./cmd/forge-core
```

Install the core plugin server(s).

```bash
sudo mv forge-core_darwin_amd64 /usr/share/forge/plugins/forge/core/
sudo mv forge-core_darwin_arm64 /usr/share/forge/plugins/forge/core/
sudo mv forge-core_linux_amd64 /usr/share/forge/plugins/forge/core/
sudo mv forge-core_linux_arm64 /usr/share/forge/plugins/forge/core/
sudo mv forge-core_windows_amd64.exe /usr/share/forge/plugins/forge/core/
sudo mv forge-core_windows_arm64.exe /usr/share/forge/plugins/forge/core/
sudo chmod 755 /usr/share/forge/plugins/forge/core/forge-core_*
```

Build the discovery plugin server(s).  Skip the platforms you don't intend to manage or run the tool from.

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o forge-discover_darwin_amd64 ./cmd/forge-discover
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o forge-discover_darwin_arm64 ./cmd/forge-discover
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forge-discover_linux_amd64 ./cmd/forge-discover
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o forge-discover_linux_arm64 ./cmd/forge-discover
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o forge-discover_windows_amd64.exe ./cmd/forge-discover
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o forge-discover_windows_arm64.exe ./cmd/forge-discover
```

Install the discovery plugin server(s).

```bash
sudo mv forge-discover_darwin_amd64 /usr/share/forge/plugins/forge/discover/
sudo mv forge-discover_darwin_arm64 /usr/share/forge/plugins/forge/discover/
sudo mv forge-discover_linux_amd64 /usr/share/forge/plugins/forge/discover/
sudo mv forge-discover_linux_arm64 /usr/share/forge/plugins/forge/discover/
sudo mv forge-discover_windows_amd64.exe /usr/share/forge/plugins/forge/discover/
sudo mv forge-discover_windows_arm64.exe /usr/share/forge/plugins/forge/discover/
sudo chmod 755 /usr/share/forge/plugins/forge/discover/forge-discover_*
```

### Basic Usage

#### Define an Inventory

Create an `inventory.hcl` file:

```hcl
vars {
    environment = "production"
    domain = "example.com"
}

transport "ssh" {
    user = "admin"
    port = 22
}

group "webservers" {
    vars {
        role = "web"
    }
}

host "web1" {
    groups = ["webservers"]
    
    vars {
        ip = "10.0.1.10"
    }
    
    transport "ssh" {
        host = "${var.ip}"
    }
}
```

#### Define a Workflow

Create a `workflow.hcl` file:

```hcl
process {
    name = "Restart myapp services"
    targets = "webservers"
    
    step "restart_service" {
        name = "Restart Service"
        module = "command"
        
        input {
            name = "systemctl"
            args = ["restart", "myapp"]
        }
    }
}
```

#### Execute the Workflow

```bash
forge run -i inventory.hcl -w workflow.hcl
```

## Development

### Project Structure

```
cmd/
  forge/              # CLI application
  forge-core/         # Core plugin server
  forge-discover/     # Discovery plugin server
pkg/
  inventory/          # Inventory types and parsing
  workflow/           # Workflow types and execution
  module/             # Module interfaces and registry
  info/               # Host information discovery
  transport/          # Transport implementations
  plugin/             # Plugin management
  result/             # Result types and handling
  hclfunction/        # HCL function implementations
  hclspec/            # HCL type specifications
  ui/                 # User interface
  util/               # Utilities
internal/
  module/             # Core plugin modules
test/
  workflow/           # Workflow parsing and resolution tests
  inventory/          # Inventory parsing and resolution tests
  integration/        # Per-platform Integration tests
```

## License

This project is licensed under the Mozilla Public License 2.0 - see the LICENSE file for details.
