/*
MIT License

Copyright (c) 2023 Kalle R. Aagaard

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package acmeproxy

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/acmeproxy"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *acmeproxy.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.acmeproxy",
		New: func() caddy.Module { return &Provider{&acmeproxy.Provider{}} },
	}
}

// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	ctx.Logger(Provider{}).Warn("Provisioning acmeproxy with endpoint " + p.Provider.Username)
	p.Provider.Username = repl.ReplaceAll(p.Provider.Username, "")
	ctx.Logger(Provider{}).Warn("Provisioning acmeproxy with endpoint " + p.Provider.Username)
	p.Provider.Password = repl.ReplaceAll(p.Provider.Password, "")
	p.Provider.Endpoint = repl.ReplaceAll(p.Provider.Endpoint, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	acmeproxy [<endpoint>] {
//		endpoint [<endpoint>]
//	    username <username>
//	    password <password>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.Endpoint = d.Val()
			if d.NextArg() {
				return d.ArgErr()
			}
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "endpoint":
				if d.NextArg() {
					if p.Provider.Endpoint != "" {
						return d.Err("API token already set")
					}
					p.Provider.Endpoint = d.Val()
					if d.NextArg() {
						return d.ArgErr()
					}
				} else {
					return d.ArgErr()
				}
			case "username":
				if d.NextArg() {
					p.Provider.Username = d.Val()
					if d.NextArg() {
						return d.ArgErr()
					}
				} else {
					return d.ArgErr()
				}

			case "password":
				if d.NextArg() {
					p.Provider.Password = d.Val()
					if d.NextArg() {

						return d.ArgErr()
					}
				} else {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.Endpoint == "" {
		return d.Err("endpoint must be specified")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
