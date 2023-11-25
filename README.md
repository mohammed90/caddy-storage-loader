Caddy Config Loader from Storage
===

Caddy supports [dynamic config loading](https://caddyserver.com/docs/json/admin/config/load/) as part of an experimental feature that was [introduced in v2.4.0](https://github.com/caddyserver/caddy/releases/tag/v2.4.0). This module tells Caddy to load its configuration from any Caddy storage module. This means it's possible to store Caddy configuration in any caddy-compatible storage, e.g. [database](https://caddyserver.com/docs/modules/caddy.storage.postgres), [s3](https://caddyserver.com/docs/modules/caddy.storage.s3), file_system, memory, [encrypted](https://caddyserver.com/docs/modules/caddy.storage.encrypted), or [any module that's part of the `caddy.storage.` namespace](https://caddyserver.com/docs/modules/), and then load it into Caddy dynamically. This is useful for use-cases where you want to store your Caddy configuration in a [database](https://caddyserver.com/docs/modules/caddy.storage.postgres), or any form of shared location with multiple Caddy instances.

## Example

<details>
<summary>This configuration file combination eventually configures Caddy to respond with `OK`. Store this configuration file in caddy storage under the key `config/caddy.json`</summary>

```json
{
	"admin": {
		"listen": "localhost:2999"
	},
	"apps": {
		"http": {
			"http_port": 9080,
			"https_port": 9443,
			"servers": {
				"srv0": {
					"listen": [
						":8443"
					],
					"routes": [
						{
							"match": [
								{
									"host": [
										"localhost"
									]
								}
							],
							"handle": [
								{
									"handler": "subroute",
									"routes": [
										{
											"handle": [
												{
													"body": "OK!",
													"handler": "static_response"
												}
											]
										}
									]
								}
							],
							"terminal": true
						}
					]
				}
			}
		},
		"pki": {
			"certificate_authorities": {
				"local": {
					"install_trust": false
				}
			}
		}
	}
}
```

</details>

Run Caddy with the following config:
```json
{
	"admin": {
		"listen":"localhost:2019",
		"config": {
			"load": {
				"module": "storage"
			}
		}
	}
}
```
