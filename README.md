![ATLAS](https://img.shields.io/badge/Yosebyte-ATLAS-blue)
![GitHub License](https://img.shields.io/github/license/yosebyte/atlas)
[![Go Report Card](https://goreportcard.com/badge/github.com/yosebyte/atlas)](https://goreportcard.com/report/github.com/yosebyte/atlas)

# ATLAS

Another Transport Layer Access Service from the Yosebyte Collections.

## Usage

```
atlas <core_mode>://<server_addr>/<access_addr>?<log=level>#<user_agent>

# Run as server
atlas server://example.org/127.0.0.1:128?log=debug#atlas

# Run as client
atlas client://example.org/127.0.0.1:8080?log=warn#atlas
```
