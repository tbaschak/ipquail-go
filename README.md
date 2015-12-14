# ipquail-go 1.0 #
Here is a go version of ipquail.com. It uses Go-lang and traffic on the backend.

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
			proxy_pass		http://127.0.0.1:3000/ip;
			proxy_set_header	X-Forwarded-For $remote_addr;
			include			proxy_params;
		}
		location /ptr {
			proxy_pass		http://127.0.0.1:3000/ptr;
			proxy_set_header	X-Forwarded-For $remote_addr;
			include			proxy_params;
		}
		location /api/ {
			proxy_pass		http://127.0.0.1:3000/api/;
			proxy_set_header	X-Forwarded-For $remote_addr;
			include			proxy_params;
		}
	}

## dns zone config ##

The front-end web application uses specific IPv4 and IPv6 only hostnames. The base hostname should be both IPv4 and IPv6 capable.

	$ORIGIN ipquail.com.
		A		<ip4>
		AAAA	<ip6>
	4	A		<ip4>
	6	AAAA	<ip6>

### traffic ###

https://github.com/pilu/traffic

(short version `go get github.com/pilu/traffic`)


