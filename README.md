# mini-httpd

`mini-httpd` is a super-simple, minimal http server used to serve static files
from a directory.

---

**Caution: `mini-httpd` is a development tool and not designed for production systems.**

---

# Install

## Build yourself

`mini-httpd` is written using [Go](https://golang.org/). You need a working Go
installation. To install the latest version, run

```
$ go install github.com/halimath/mini-httpd
```

# Usage

To simply serve the content of the current directory, run

```
$ mini-httpd
```

This will start the server listening on `:8080` for HTTP connections.

The following command line switches can be used to customize the behaviour:


Switch | Default Value | Description
-- | -- | --
`doc-root` | `.` | Document root
`http-address` | `:8080` | Network address to bind to listen for incoming requests
`no-log` | _no argument_ | Disable logging

# Author

`mini-httpd` is written by [Alexander Metzner](mailto:alexander.metzner@gmail.com)

# License

> 
> Copyright (c) 2021 Alexander Metzner.
> 
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
> 
>     http://www.apache.org/licenses/LICENSE-2.0
> 
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.
>
