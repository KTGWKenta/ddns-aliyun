package config

import "time"

domains: {
	"kenta.im": {
		provider: "aliyun"
		authArgs: {
			region:       "cn-shanghai"
			accessKey:    ""
			accessSecret: ""
		}
		prefixes: [
			{
				ipType:     "v4"
				record:     "*.s86510hwg"
				recordType: "A"
			},
			{
				ipType:     "v6"
				record:     "*.s86510hwg"
				recordType: "AAAA"
			},
			{
				ipType:     "v4"
				record:     "s86510hwg"
				recordType: "A"
			},
			{
				ipType:     "v6"
				record:     "s86510hwg"
				recordType: "AAAA"
			},
			{
				ipType:     "v4"
				record:     "v4.s86510hwg"
				recordType: "A"
			},
			{
				ipType:     "v6"
				record:     "v6.s86510hwg"
				recordType: "AAAA"
			},
		]
	}
}

lookups: {
	timeout:  time.ParseDuration("50s")
	interval: time.ParseDuration("1m")
	v4Enable: true
	v4Addr:   "http://ipv4.lookup.test-ipv6.com/ip/"
	v4Path:   "ip"
	v6Enable: false
	v6Addr:   "http://ipv6.lookup.test-ipv6.com/ip/"
	v4Path:   "ip"
}
