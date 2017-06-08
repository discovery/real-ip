package real_ip

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

var remote_addr, _ = NewRemoteAddr([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
})

func TestUtils_get_remote_addr_net_empty_request(t *testing.T) {
	req := &http.Request{}
	_, err := remote_addr.GetIP(req)
	assert.NotNil(t, err)
}

func TestUtils_get_remote_addr_net_invalid_remote_addr_bad_ip_address(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.41234:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{},
		},
	}
	_, err := remote_addr.GetIP(req)
	assert.NotNil(t, err)
}

func TestUtils_get_remote_addr_net_invalid_remote_addr_missing_port(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4",
		Header: http.Header{
			"X-Forwarded-For": []string{},
		},
	}
	_, err := remote_addr.GetIP(req)
	assert.NotNil(t, err)
}

func TestUtils_get_remote_addr_net_invalid_remote_addr_missing_ip_address(t *testing.T) {
	req := &http.Request{
		RemoteAddr: ":1234",
		Header: http.Header{
			"X-Forwarded-For": []string{},
		},
	}
	_, err := remote_addr.GetIP(req)
	assert.NotNil(t, err)
}

func TestUtils_GetClientIP_empty_x_forwarded_for(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{""},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("1.2.3.4"), res)
}

func TestUtils_GetClientIP_invalid1_x_forwarded_for(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"1.2"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("1.2.3.4"), res)
}

func TestUtils_GetClientIP_invalid2_x_forwarded_for(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"asdfad, 4.4.4.4"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("4.4.4.4"), res)
}

func TestUtils_GetClientIP_invalid3_x_forwarded_for(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"asdfad, 4.4.4.4|"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("1.2.3.4"), res)
}

func TestUtils_GetClientIP_x_forwarded_for_test1(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.0.1"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("1.2.3.4"), res)
}

func TestUtils_GetClientIP_x_forwarded_for_test2(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"4.4.4.4"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("4.4.4.4"), res)
}

func TestUtils_GetClientIP_x_forwarded_for_test3(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.0.1, 4.4.4.4"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("4.4.4.4"), res)
}

func TestUtils_GetClientIP_x_forwarded_for_test4(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.0.1, 4.4.4.4,10.0.0.1"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("4.4.4.4"), res)
}

func TestUtils_GetClientIP_x_forwarded_for_test5(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-For": []string{"8.8.4.4,127.0.0.1, 4.4.4.4,10.0.0.1"},
		},
	}
	res, _ := remote_addr.GetIP(req)
	assert.Equal(t, net.ParseIP("4.4.4.4"), res)
}
