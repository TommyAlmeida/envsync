# envsync

`envsync` is a developer tool designed to synchronize environment variables across different environments and cloud providers. It provides a flexible configuration system and supports adapters for various platforms, including AWS, Google Cloud Platform (GCP) (soon), and Azure (soon), with the ability to extend to other providers.

## Features

- Synchronize environment variables based on a schema.
- Support for multiple cloud providers via adapters (e.g., AWS).
- Configurable via YAML files or environment variables.
- Validate configurations before synchronization.
- Extensible architecture for adding new adapters.

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/tommyalmeida/envsync.git
   cd envsync
   ```

2. Build the binary:

   ```bash
   make build
   ```

3. Move the binary to your PATH (optional):
   ```bash
   sudo mv bin/envsync /usr/local/bin/
   ```

## Usage

### Configuration

Create a envsync.yaml file in your project root or specify a custom config file using the -c flag. Example envsync.yaml:

```yaml
schema:
  variables:
    PORT:
      required: true
      type: number
      default: "3000"
      description: "Server port"
rules:
  require_all: false
  allow_extra: true
  ignore_patterns:
    - "^TEMP_.*"
    - "^DEBUG_.*"
adapter:
  name: "aws"
  config:
    region: "us-west-2"
```

### Running envsync

```bash
envsync
```

### WIP: Supported Adapters:

- AWS: Syncs environment variables to AWS Systems Manager Parameter Store.

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

## Contributing

We welcome contributions! To get started:

1. **Fork** the repository.
2. **Create a feature branch**:

```bash
git checkout -b feature/your-feature
```

3. **Commit your changes**:

```bash
git commit -m "Describe your changes"
```

4. **Push to your branch**:

```bash
git push origin feature/your-feature
```

5. **Open a Pull Request** on GitHub.

Please ensure your code follows the project's style and includes relevant tests and documentation.

## Acknowledgments

Inspired by the need for seamless environment management across cloud providers and my current pain in the ass.
