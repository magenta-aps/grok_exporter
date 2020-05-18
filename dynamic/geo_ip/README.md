Plugin GeoIP
============
A plugin that converts domain names or IP addresses to country codes.

## Build
```bash
go build -buildmode=plugin -o ../plugin_geoip.so
```

## Example
This example utilizes the following `grok_exporter` configuration file:
```yaml
global:
  config_version: 3
input:
  type: webhook
  webhook_path: /webhook
  webhook_format: json_single
  webhook_json_selector: .message
  webhook_text_bulk_separator: ""
metrics:
- type: counter
  name: logins_by_country
  help: Total number of logins by country via GeoIP
  match: 'Login occured'
  labels:
    ip: '{{ geoip (index .extra "ip") }}'
server:
  host: "[::]"
  port: 9144
```
Where `geoip` is used to convert `extra` JSON context IPs to country codes.

Start the `grok_exporter` and simulate some log lines using:
```bash
$ curl localhost:9144/webhook --data '{"message": "Login occured", "user": "Skeen", "ip": "40.68.218.4"}'
$ curl localhost:9144/webhook --data '{"message": "Login occured", "user": "Fabian", "ip": "52.57.35.183"}'
```

Now check the metrics output at: `localhost:9144/metrics`:
```
# HELP logins_by_country Total number of logins by country via GeoIP
# TYPE logins_by_country counter
logins_by_country{ip="DE"} 1
logins_by_country{ip="NL"} 1
```
