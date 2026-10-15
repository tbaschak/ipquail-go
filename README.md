# ipquail-go #

Here is a Go version of ipquail.com. The backend uses Go's standard `net/http` package.

Go 1.26 or newer is required.

    cd ipquail
    go build
    ./ipquail

## nginx config ##

The backend application is setup to use the `x-forwarded-for` header, the following nginx reverse proxy config sets up this header properly.

	map $http_user_agent $ipquailindex {
		default /index.html;
		~curl /ip;
		~Python-urllib /ip;
	}

	server {
		listen 80 default_server;
		listen [::]:80 default_server ipv6only=on;

		root /home/website/ipquail-go/ipquail/public;

		charset utf-8;

		server_name ipquail.com 4.ipquail.com 6.ipquail.com www.ipquail.com;

		# Special endpoint depending on useragent
		location = / { rewrite ^ $ipquailindex; }

		# staticly serve files
		location / {
			try_files $uri $uri/ =404;
		}

        # revproxy endpoints
        location /ip {
                proxy_pass              http://127.0.0.1:3000/ip;
                proxy_set_header        X-Forwarded-For $remote_addr;
                include                 proxy_params;
                proxy_headers_hash_bucket_size  256;
        }
        location /ptr {
                proxy_pass              http://127.0.0.1:3000/ptr;
                proxy_set_header        X-Forwarded-For $remote_addr;
                include                 proxy_params;
                proxy_headers_hash_bucket_size  256;
        }
        location /asn {
                proxy_pass              http://127.0.0.1:3000/asn;
                proxy_set_header        X-Forwarded-For $remote_addr;
                include                 proxy_params;
                proxy_headers_hash_bucket_size  256;
        }
        location /api/ {
                proxy_pass              http://127.0.0.1:3000/api/;
                proxy_set_header        X-Forwarded-For $remote_addr;
                include                 proxy_params;
                proxy_headers_hash_bucket_size  256;
        }
	}

## dns zone config ##

The front-end web application uses specific IPv4 and IPv6 only hostnames. The base hostname should be both IPv4 and IPv6 capable.

	$ORIGIN ipquail.com.
		A		<ip4>
		AAAA	<ip6>
	4	A		<ip4>
	6	AAAA	<ip6>

## Backend

The backend listens on port `3000` through the nginx reverse proxy shown above. It has no third-party Go dependencies.

