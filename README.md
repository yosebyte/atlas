![ATLAS](https://img.shields.io/badge/Yosebyte-ATLAS-blue)
![GitHub License](https://img.shields.io/github/license/yosebyte/atlas)
[![Go Report Card](https://goreportcard.com/badge/github.com/yosebyte/atlas)](https://goreportcard.com/report/github.com/yosebyte/atlas)

# ATLAS

Another Transport Layer Access Service from the Yosebyte Collections.

## Usage

```
atlas <core_mode>://<server_addr>#<access_addr>?<log=level>

# Run as server
atlas server://10.1.0.1:10101?log=debug

# Run as client
atlas client://10.1.0.1:10101#127.0.0.1:8080
```
