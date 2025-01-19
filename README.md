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

- Create crediential for key and password
```sh
telepath crediential add -T KEY -K "/path-to/id_rsa" cred1

telepath crediential add -T PASS -P "my-password" cred2
```

- Create Host with crediential
```sh
telepath host add -H <host-ip> -P 22 -U <username> -C cred1 mylab1
telepath host add -H <host-ip> -P 22 -U <username> -C cred2 mylab2
```

- Create tunnel for mysql connection
```sh
telepath tunnel add -L 3306 -H localhost -R 3306 -C mylab1,mylab2 tunnel1
```

- Start and stop tunnel
```sh
telepath tunnel start tunnel1
telepath tunnel stop tunnel1
```
