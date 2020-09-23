domains : {
	"kenta.im" : {
		provider : "aliyun",
		authArgs : {
			region : "cn-shanghai",
			accessKey : "LTA******************S3i",
			accessSecret : "HG4************************pr5",
		}
		prefixes : [
		    {
		        ipType : "v4"
		        record : "s86510hwg"
		        recordType : "A"
		    },
		    {
		        ipType : "v4"
		        record : "v4.s86510hwg"
		        recordType : "A"
		    },
		    {
		        ipType : "v6"
		        record : "s86510hwg"
		        recordType : "AAAA"
		    },
		    {
		        ipType : "v6"
		        record : "v6.s86510hwg"
		        recordType : "AAAA"
		    }
		]
	}
}

lookups : {
	v4Addr: "http://ipv4.lookup.test-ipv6.com/ip/"
	v4Path: "ip"
	v6Addr: "http://ipv6.lookup.test-ipv6.com/ip/"
	v4Path: "ip"
}