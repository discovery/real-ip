package real_ip

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	. "testing"
)

var remote_addr, _ = NewRemoteAddr([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
})

func TestUtils_get_remote_addr_net_empty_request(t *T) {
	req := &http.Request{}
	_, err := remote_addr.GetIP(req)
	assert.NotNil(t, err)
}

func TestStuff(t *T) {
	DoTest(t, "bad_ip_address",            "", 	  "1.2.3.41234:1234", []string{}, nil)
	DoTest(t, "missing_port",              "", 	  "1.2.3.4",          []string{}, nil)
	DoTest(t, "missing_ip_address",        "", 	  ":1234",            []string{}, nil)
	DoTest(t, "empty_x_forwarded_for",     "1.2.3.4", "1.2.3.4:1234",     []string{""}, nil)
	DoTest(t, "invalid_x_forwarded_for 1", "1.2.3.4", "1.2.3.4:1234",     []string{"1.2"}, nil)
	DoTest(t, "invalid_x_forwarded_for 2", "4.4.4.4", "1.2.3.4:1234",     []string{"asdfad, 4.4.4.4"}, nil)
	DoTest(t, "invalid_x_forwarded_for 3", "1.2.3.4", "1.2.3.4:1234",     []string{"asdfad, 4.4.4.4|"}, nil)
	DoTest(t, "x_forwarded_for_test 1",    "1.2.3.4", "1.2.3.4:1234",     []string{"127.0.0.1"}, nil)
	DoTest(t, "x_forwarded_for_test 2",    "4.4.4.4", "1.2.3.4:1234",     []string{"4.4.4.4"}, nil)
	DoTest(t, "x_forwarded_for_test 3",    "4.4.4.4", "1.2.3.4:1234",     []string{"127.0.0.1, 4.4.4.4"}, nil)
	DoTest(t, "x_forwarded_for_test 4",    "4.4.4.4", "1.2.3.4:1234",     []string{"127.0.0.1, 4.4.4.4,10.0.0.1"}, nil)
	DoTest(t, "x_forwarded_for_test 5",    "4.4.4.4", "1.2.3.4:1234",     []string{"8.8.4.4,127.0.0.1, 4.4.4.4,10.0.0.1"}, nil)
}

func DoTest(t *T, test_name string, expected string, remoteAddr string, x_forwarded_for []string, extra_headers *http.Header) {
	t.Run(test_name, func(t *T) {
		if extra_headers == nil {
			extra_headers = &http.Header{}
		}
		(*extra_headers)["X-Forwarded-For"] = x_forwarded_for
		req := &http.Request{
			RemoteAddr: remoteAddr,
			Header: *extra_headers,
		}
		res, err := remote_addr.GetIP(req)
		if expected == "" {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, net.ParseIP(expected), res)
		}
	})
}
