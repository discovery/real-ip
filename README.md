real-ip
=======

Real IP is a minimal Go library to extract the correct IP from an
`X-Forwarded-For` http header similar to Nginx's [realip
module][real_ip]. It is based on a heavily modified version of the
[xff][xff] code.

[real_ip]: https://nginx.org/en/docs/http/ngx_http_realip_module.html
[xff]: https://github.com/sebest/xff

Example
-------

```go
package main

import (
    "net/http"
    "github.com/discovery/real-ip"
)
func main() {
    remote_addr, _ := real_ip.NewRemoteAddr([]string{"10.0.0.0/8", "192.168.0.0/16", "1.2.3.4/32"})
    // (Insert error checking here)

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _ := remote_addr.GetIP(r)
        // (Insert error checking here)
        w.Write([]byte("hello from " + ip.String() + "\n"))
    })

    http.ListenAndServe(":8080", http.Handler(handler))
}
```

Copyright & License
===================

Copyright (c) 2018 Discovery Communications

Licensed under the MIT License. See LICENSE file for details.
