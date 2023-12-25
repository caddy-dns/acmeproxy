ACMEProxy module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records via an ACMEProxy server.

## Caddy module name

```
dns.providers.acmeproxy
```

## XCaddy
```
xcaddy build --with github.com/caddy-dns/acmeproxy@v1.0.6
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuers/acme/) like so:

```json
{
    "module": "acme",
    "challenges": {
        "dns": {
            "provider": {
                "name": "acmeproxy",
                "username": "user",
                "password": "pass",
                "endpoint": "https://example.com:9090"
            }
        }
    }
}
```

or with the Caddyfile:

```
tls {
    dns acmeproxy https://example.com:9000 {
        username user
        password pass
    }
}
```


## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/acmeproxy) for important information about credentials.


## Troubleshooting

### Error: `timed out waiting for record to fully propagate`

Some environments may have trouble querying the `_acme-challenge` TXT record from dnsproviders. Verify in the providers dashboard that the temporary record is being created.

If the record does exist, your DNS resolver may be caching an earlier response before the record was valid. You can instead configure Caddy to use an alternative DNS resolver such as [Cloudflare's official `1.1.1.1`](https://www.cloudflare.com/en-gb/learning/dns/what-is-1.1.1.1/).

Add a custom `resolver` to the [`tls` directive](https://caddyserver.com/docs/caddyfile/directives/tls):

```
tls {
  dns acmeproxy https://example.com:9000 {
    username user
    password pass
  }
  resolvers 1.1.1.1
}
```

Or with Caddy JSON to the `acme` module: [`challenges.dns.provider.resolvers: ["1.1.1.1"]`](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/challenges/dns/resolvers/).

