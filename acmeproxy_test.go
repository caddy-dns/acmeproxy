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
	"reflect"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/acmeproxy"
)

func TestProvider_UnmarshalCaddyfile(t *testing.T) {

	tests := []struct {
		name    string
		cfile   string
		want    *acmeproxy.Provider
		wantErr bool
	}{
		{
			name: "Basic",
			cfile: `acmeproxy https://example.com:9000 {
				username user
				password pass
			}`,
			want: &acmeproxy.Provider{
				Endpoint: "https://example.com:9000",
				Credentials: acmeproxy.Credentials{
					Username: "user",
					Password: "pass",
				},
			},
			wantErr: false,
		},
		{
			name: "Alt basic",
			cfile: `acmeproxy {
				endpoint https://example.com:9000
				username user
				password pass
			}`,
			want: &acmeproxy.Provider{
				Endpoint: "https://example.com:9000",
				Credentials: acmeproxy.Credentials{
					Username: "user",
					Password: "pass",
				},
			},
			wantErr: false,
		},
		{
			name: "Double endpoint error",
			cfile: `acmeproxy https://example.com:9000 {
				endpoint https://example.com:9000
				username user
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Extra head arg error",
			cfile: `acmeproxy https://example.com:9000 dummy {
				username user
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Extra user arg error",
			cfile: `acmeproxy https://example.com:9000 {
				username user dummy
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Extra pass arg error",
			cfile: `acmeproxy https://example.com:9000 {
				username user 
				password pass dummy
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Extra endpoint arg error",
			cfile: `acmeproxy {
				endpoint https://example.com:9000 dummy
				username user 
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing endpoint error",
			cfile: `acmeproxy {
				username user
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing endpoint value error",
			cfile: `acmeproxy {
				endpoint
				username user
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing username value error",
			cfile: `acmeproxy {
				endpoint https://example.com:9000
				username 
				password pass
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing password value error",
			cfile: `acmeproxy {
				endpoint https://example.com:9000
				username user
				password 
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unknown value error",
			cfile: `acmeproxy {
				endpoint https://example.com:9000
				username user
				password pass
				certificates 1
			}`,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				&acmeproxy.Provider{},
			}
			err := p.UnmarshalCaddyfile(caddyfile.NewTestDispenser(tt.cfile))
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if got := p.Provider; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Provider.CaddyModule() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
