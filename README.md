# 🌐 Atlas

Atlas is a versatile, high-performance networking tool designed for creating secure communication channels across various environments. Operating in both server and client modes, it provides seamless TLS encryption through either Let's Encrypt certificates or self-signed options. With its zero-configuration approach and single-binary deployment, Atlas excels in scenarios ranging from corporate security compliance and microservices communication to remote work infrastructure and multi-cloud connectivity. Its lightweight footprint and intuitive URL-style command syntax make it the ideal solution for developers, system administrators, and security professionals seeking robust, flexible networking capabilities without complexity.

## 📋 Table of Contents

- [Features](#-features)
- [Requirements](#-requirements)
- [Installation](#-installation)
  - [Option 1: Pre-built Binaries](#-option-1-pre-built-binaries)
  - [Option 2: Using Go Install](#-option-2-using-go-install)
  - [Option 3: Building from Source](#️-option-3-building-from-source)
  - [Option 4: Using Container Image](#-option-4-using-container-image)
- [Usage](#-usage)
  - [Server Mode](#️-server-mode)
  - [Client Mode](#-client-mode)
- [Configuration](#️-configuration)
  - [Log Levels](#-log-levels)
- [Examples](#-examples)
  - [Running as a secure atlas server with a domain](#-running-as-a-secure-atlas-server-with-a-domain)
  - [Setting up a local atlas server with self-signed certificate](#-setting-up-a-local-atlas-server-with-self-signed-certificate)
  - [Connecting to an atlas server](#-connecting-to-an-atlas-server)
- [How It Works](#-how-it-works)
- [Common Use Cases](#-common-use-cases)
- [Troubleshooting](#-troubleshooting)
  - [Certificate Issues](#-certificate-issues)
  - [Connection Problems](#-connection-problems)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **🔄 Dual Operating Modes**: Run as a server to accept connections or as a client to initiate them
- **🔒 Automatic TLS Certificate Management**: Seamless integration with Let's Encrypt for valid certificates
- **🛡️ Self-Signed Certificate Generation**: Perfect for development, testing, and internal deployments
- **📊 Flexible Logging System**: Configurable verbosity with five distinct logging levels
- **🌐 HTTP/HTTPS Protocol Support**: Full support for both protocols with transparent handling
- **📦 Single-Binary Deployment**: Simple to distribute and install with no dependencies
- **⚙️ Zero Configuration Files**: Everything is specified via command-line arguments
- **🚀 Low Resource Footprint**: Minimal CPU and memory usage even under heavy load

## 📋 Requirements

- Go 1.24 or higher (for building from source)
- Network access for Let's Encrypt certificate issuance (server mode with domain)
- Admin privileges may be required for binding to ports below 1024

## 📥 Installation

### 💾 Option 1: Pre-built Binaries

Download the latest release for your platform from our [releases page](https://github.com/yosebyte/atlas/releases).

### 🔧 Option 2: Using Go Install

```bash
go install github.com/yosebyte/atlas/cmd/atlas@latest
```

### 🛠️ Option 3: Building from Source

```bash
# Clone the repository
git clone https://github.com/yosebyte/atlas.git

# Build the binary
cd atlas
go build -o atlas ./cmd/atlas

# Optional: Install to your GOPATH/bin
go install ./cmd/atlas
```

### 🐳 Option 4: Using Container Image

Atlas is available as a container image on GitHub Container Registry:

```bash
# Pull the container image
docker pull ghcr.io/yosebyte/atlas:latest

# Run in server mode
docker run -d --rm -p 443:443 ghcr.io/yosebyte/atlas server://example.org/0.0.0.0:443

# Run in client mode
docker run -d --rm -p 8080:8080 ghcr.io/yosebyte/atlas client://example.org:443/127.0.0.1:8080
```

## 🚀 Usage

Atlas can be run in either server mode or client mode with a single, intuitive URL-style command:

### 🖥️ Server Mode

```bash
atlas server://<server_addr>/<access_addr>?log=<level>
```

- `server_addr`: Domain name for TLS certificate (leave blank for self-signed certificates)
- `access_addr`: Listen address to bind to (default: random port on 127.0.0.1)
- `log`: Log level (debug, info, warn, error, fatal)

Example:
```bash
# Using Let's Encrypt certificates (requires valid domain)
atlas server://example.org/127.0.0.1:8080

# Using self-signed certificates (for development)
atlas server:///127.0.0.1:8080?log=debug
```

### 📱 Client Mode

```bash
atlas client://<server_addr>/<access_addr>?log=<level>
```

- `server_addr`: Atlas server address (hostname:port)
- `access_addr`: Local address to bind to (default: random port on 127.0.0.1)
- `log`: Log level (debug, info, warn, error, fatal)

Example:
```bash
atlas client://example.org:443/127.0.0.1:8080
```

## ⚙️ Configuration

Atlas uses a minimalist approach with command-line parameters rather than configuration files:

### 📝 Log Levels

- `debug`: Verbose debugging information - shows all operations and connections
- `info`: General operational information (default) - shows startup, shutdown, and key events
- `warn`: Warning conditions - only shows potential issues that don't affect core functionality
- `error`: Error conditions - shows only problems that affect functionality
- `fatal`: Critical conditions - shows only severe errors that cause termination

## 📚 Examples

### 🔐 Running as a secure atlas server with a domain

```bash
atlas server://example.org/0.0.0.0:443
```

This will:
1. Obtain a valid TLS certificate from Let's Encrypt for example.org
2. Start the server, binding to 0.0.0.0:443
3. Log at the INFO level

### 🧪 Setting up a local atlas server with self-signed certificate

```bash
atlas server:///0.0.0.0:8443?log=debug
```

This will:
1. Generate a self-signed certificate
2. Start the server, binding to all interfaces on port 8443
3. Log at the DEBUG level

### 🔌 Connecting to an atlas server

```bash
atlas client://example.org:443/127.0.0.1:8080
```

This will:
1. Connect to the atlas server at example.org:443
2. Bind to 127.0.0.1:8080 locally
3. Log at the INFO level

## 🔍 How It Works

Atlas creates a secure tunnel using TLS:

1. **Server Mode**: Listens for incoming connections and:
   - If using a domain name, attempts to obtain a valid certificate via Let's Encrypt
   - If using an IP address, generates a self-signed certificate
   - For CONNECT requests, establishes a tunnel to the requested destination
   - For HTTP requests, acts as a reverse gateway

2. **Client Mode**: 
   - Connects to an Atlas server using TLS
   - Exposes a local endpoint for applications to connect to
   - Forwards all traffic through the encrypted tunnel

## 💡 Common Use Cases

- **🔒 Corporate Security Compliance**: Create secure access points to internal resources with proper encryption and audit trails
- **🌍 Geographic Access**: Access region-restricted services by routing through Atlas servers in different locations
- **🧪 Development & Testing**: Create realistic production-like environments during development with TLS
- **🔄 Microservices Architecture**: Enable secure communication between distributed service components
- **⚡ Edge Computing**: Deploy at edge locations to reduce latency while maintaining TLS security
- **🏛️ Legacy System Integration**: Add modern security protocols to legacy applications without modifying them
- **📱 Remote Work Infrastructure**: Provide secure access to internal tools and services for remote teams
- **☁️ Multi-cloud Connectivity**: Establish secure channels between resources hosted across different cloud providers

## 🔧 Troubleshooting

### 📜 Certificate Issues
- Ensure your domain points to the server where Atlas is running
- Check that port 80 is accessible for Let's Encrypt verification
- For self-signed certificates, ensure clients are configured to accept them

### 🔌 Connection Problems
- Verify firewall settings allow traffic on the specified ports
- Check that the server address is correctly specified in client mode
- Increase log level to debug for more detailed connection information

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.
