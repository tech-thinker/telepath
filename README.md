# telepath
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tech-thinker/telepath)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/tech-thinker/telepath/release.yaml)
![GitHub](https://img.shields.io/github/license/tech-thinker/telepath)
![GitHub All Releases](https://img.shields.io/github/downloads/tech-thinker/telepath/total)
![GitHub last commit](https://img.shields.io/github/last-commit/tech-thinker/telepath)
![GitHub forks](https://img.shields.io/github/forks/tech-thinker/telepath)
![GitHub top language](https://img.shields.io/github/languages/top/tech-thinker/telepath)

Telepath is a modern, intelligent CLI tool for seamless and secure port forwarding.
Designed with versatility and ease of use in mind, Telepath enables developers and system administrators to create complex forwarding paths across multiple hosts effortlessly.
Whether you're working with password-based or keyfile authentication, single or multiple jump hosts, Telepath has you covered.

## Why Choose Telepath?
- Multi-Jump Host Support: Telepath allows seamless port forwarding through multiple intermediate hosts, making it ideal for accessing restricted networks.
- Secure Authentication: Supports password-based and keyfile authentication, ensuring flexibility and security.
- CLI Simplicity: Its intuitive command-line interface simplifies complex operations with straightforward commands.
- Daemon Mode: Runs as a background service for continuous operation, perfect for long-running port-forwarding tasks.
- Customizable & Open Source: It's open source and developer-friendly, so you can tweak it for your needs.

## Installation
Download and install executable binary from GitHub releases page.

### Using homebrew
```sh
brew tap tech-thinker/tap
brew install telepath
```

### Linux Installation
```sh
# Use latest tag name from release page
TAG=<tag-name>

curl -sL "https://github.com/tech-thinker/telepath/releases/download/${TAG}/telepath-linux-amd64" -o telepath
chmod +x telepath
sudo mv telepath /usr/bin
```

### MacOS Installation
```sh
# Use latest tag name from release page
TAG=<tag-name>

curl -sL "https://github.com/tech-thinker/telepath/releases/download/${TAG}/telepath-darwin-amd64" -o telepath
chmod +x telepath
sudo mv telepath /usr/bin
```

### Windows Installation
```sh
# Use latest tag name from release page
TAG=<tag-name>

curl -sL "https://github.com/tech-thinker/telepath/releases/download/${TAG}/telepath-windows-amd64.exe" -o telepath.exe
telepath.exe
```

### Verify checksum
```sh
# Use latest tag name from release page
TAG=<tag-name>

# Using sha256sum
curl -sL "https://github.com/tech-thinker/telepath/releases/download/${TAG}/checksum-sha256sum.txt" -o checksum-sha256sum.txt
sha256sum --ignore-missing --check checksum-sha256sum.txt

# Using md5sum
curl -sL "https://github.com/tech-thinker/telepath/releases/download/${TAG}/checksum-md5sum.txt" -o checksum-md5sum.txt
md5sum --ignore-missing --check checksum-md5sum.txt
```
Output:
```sh
telepath-darwin-amd64: OK
telepath-darwin-amd64.tar.gz: OK
telepath-darwin-arm64: OK
telepath-darwin-arm64.tar.gz: OK
telepath-linux-amd64: OK
telepath-linux-amd64.tar.gz: OK
telepath-linux-arm: OK
telepath-linux-arm.tar.gz: OK
telepath-linux-arm64: OK
telepath-linux-arm64.tar.gz: OK
telepath-windows-amd64.exe: OK
telepath-windows-i386.exe: OK
```

## CLI Guide
- `telepath` help
```sh
telepath -h
```

- Start tunnel using cli
```sh
telepath -f /etc/telepath/telepath.json
```

## Define Config file
Config file is a JSON file which contains list of config. Here I have attached a sample config file-

```json
[
  {
    "name": "mongodb",
    "type": "L",
    "localPort": 27017,
    "localHost": "0.0.0.0",
    "remotePort": 27017,
    "remoteHost": "0.0.0.0",
    "server": {
      "host": "final-host-ip",
      "port": 22,
      "username": "user",
      "authType": "KEY",
      "password": "",
      "key": "/etc/autossh/id_rsa",
      "passphrase": "passphrase",
      "jump": {
        "host": "jump-host-ip",
        "port": 22,
        "username": "user",
        "authType": "KEY",
        "password": "",
        "key": "/etc/autossh/id_rsa",
        "passphrase": "passphrase"
      }
    }
  },
  {
    "name": "mysql",
    "type": "R",
    "localPort": 3306,
    "localHost": "0.0.0.0",
    "remotePort": 3306,
    "remoteHost": "0.0.0.0",
    "server": {
      "host": "final-host-ip",
      "port": 22,
      "username": "user",
      "authType": "KEY",
      "password": "",
      "key": "/etc/autossh/id_rsa",
      "passphrase": "passphrase",
      "jump": {
        "host": "jump-host-ip",
        "port": 22,
        "username": "user",
        "authType": "KEY",
        "password": "",
        "key": "/etc/autossh/id_rsa",
        "passphrase": "passphrase"
      }
    }
  }
]
```
