# mr-prober

> Mr Prober runs system probes and exposes the results on a Prometheus endpoint.
> This program is meant to be run as a DaemonSet on all nodes in `Kubernetes/Openshift` clusters.

[[_TOC_]]

## Docker

## Kubernetes

```bash
make install
#make uninstall
```

## Development

### Configuration

`Rules` must have a unique name. This is mandatory.

### net

```
rules:
  - name: "Network connectivity test"
    probe: "net"
    args:
      - "<tcp|udp>://<host:port>"
      - "[timeout]"
```

#### About UDP

> Because UDP does not reply to connection requests, a lack of response may indicate that the port is open, or that the packet got dropped.
> We chose to be optimistic and treat lack of response (connection timeout) as an open port.

* [mozilla/mig](https://github.com/mozilla/mig/blob/master/modules/ping/ping.go#L280)

## LICENSE

> MIT License

## Methodology
* [USEmethod/use-linux.html](https://www.brendangregg.com/USEmethod/use-linux.html)

## References
* [prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter/tree/master)
* [M/MONIT](https://mmonit.com/)
* [vfedoroff/go-netcat](https://github.com/vfedoroff/go-netcat/blob/master/main.go)
* [Go by Example: WaitGroups](https://gobyexample.com/waitgroups)
* [metrics - lightweight package for exporting metrics in Prometheus format](https://github.com/VictoriaMetrics/metrics#faq)
* [Effortless Hot Reloading in Golang](https://medium.com/@adamszpilewicz/effortless-hot-reloading-in-golang-harnessing-the-power-of-viper-4b54703f7424)
* [Prometheus/Getting Started](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md)
* [go-runtime-mixin](https://github.com/grafana/jsonnet-libs/blob/master/go-runtime-mixin/dashboards/go-runtime.json)
* [*http.Server in Go 1.8 supports graceful shutdown](https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2)
