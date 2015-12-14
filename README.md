# ipquail-go 1.0 #
Here is a go version of ipquail.com. It uses Go-lang and traffic on the backend.

    cd ipquail
    go build
    ./ipquail

## nginx config ##

The backend application is setup to use the `x-forwarded-for` header, the following nginx reverse proxy config sets up this header properly.

	server {
		listen 80 default_server;
		listen [::]:80 default_server ipv6only=on;

		root /var/www/html;

		server_name _;

		location / {
			proxy_pass	http://127.0.0.1:3000/;
			proxy_set_header X-Forwarded-For $remote_addr;
			include		proxy_params;
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


