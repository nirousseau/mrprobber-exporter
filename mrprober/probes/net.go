package probes

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"
)

const NETWORK_MTU = 1500

type Net struct {
	Name    string
	url     *url.URL
	timeout time.Duration
}

func NewNet(name string, args []string) (*Net, error) {

	const minArgs = 1
	var timeout, _ = time.ParseDuration("100ms")

	if len(args) < minArgs {
		return nil, fmt.Errorf(`Unsufficient arguments.

See:
* https://www.rfc-editor.org/rfc/rfc3986.html#section-1.1.2
* https://pkg.go.dev/time#ParseDuration

Usage:
net <raw_url> [timeout]
`)
	}

	// Required
	rawUrl := args[0]
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url arg %q for net probe %q", rawUrl, name)
	}

	// Optionals
	if len(args) >= 2 {
		timeout, err = time.ParseDuration(args[1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse timeout arg %q for net probe %q", timeout, name)
		}
	}

	return &Net{
		Name:    name,
		url:     u,
		timeout: timeout,
	}, nil
}

func (p Net) Exec() Result {

	switch s := p.url.Scheme; s {
	case "tcp":
		return p.vz()
	case "udp":
		return p.vzu()
	default:
		return Result{
			ReturnCode: Failure,
			Msg:        fmt.Sprintf("Scheme `%s` is not implemented by this probe", s),
		}
	}
}

// vz is a function similar to netcat -vz command
// Example :
// nc -vz 127.0.0.53 53
// Connection to 127.0.0.53 53 port [tcp/domain] succeeded!
func (p Net) vz() Result {

	src := os.Getenv("KUBERNETES_NODE_NAME")
	if len(src) == 0 {
		src = "localhost"
	}
	pr := Result{
		ProbeID:    p.Name,
		MetricName: fmt.Sprintf("net_%s", "tcp"),
		MetricLabels: map[string]string{
			"timeout": p.timeout.String(),
			"src":     src,
			"dest":    p.url.Host,
		},
		ReturnCode: Success,
		Msg:        fmt.Sprintf("Connection to %s (%s) succeeded!", p.url.Host, "tcp"),
	}

	conn, _ := net.DialTimeout("tcp", p.url.Host, p.timeout)
	if conn == nil {
		pr.ReturnCode = Failure
		pr.Msg = fmt.Sprintf("Connection to %s (%s) failed: Connection refused", p.url.Host, "tcp")
		return pr
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	return pr
}

// vzu is a function similar to netcat -vzu command
// Because UDP does not reply to connection requests, a lack of response may indicate that the
// port is open, or that the packet got dropped. If the timeout is not reached, it means that the packet has been
// explicitly refused by a device so that we can detect that the probe is ko.
// Example :
// nc -vzu 127.0.0.53 5353
// Connection to 127.0.0.53 5353 port [udp/mdns] succeeded!
func (p Net) vzu() Result {

	src := os.Getenv("KUBERNETES_NODE_NAME")
	if len(src) == 0 {
		src = "localhost"
	}
	pr := Result{
		ProbeID:    p.Name,
		MetricName: fmt.Sprintf("net_%s", "udp"),
		MetricLabels: map[string]string{
			"timeout": p.timeout.String(),
			"src":     src,
			"dest":    p.url.Host,
		},
		ReturnCode: Success,
		Msg:        fmt.Sprintf("Connection to %s (%s) succeeded!", p.url.Host, "udp"),
	}

	udpAddr, _ := net.ResolveUDPAddr("udp", p.url.Host)
	conn, _ := net.DialTimeout("udp", udpAddr.String(), p.timeout)
	if conn == nil {
		pr.ReturnCode = Failure
		pr.Msg = fmt.Sprintf("Timeout not reached %s (%s): Connection refused", p.url.Host, "udp")
		return pr
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// https://github.com/mozilla/mig/blob/master/modules/ping/ping.go#L280
	_, _ = conn.Write([]byte("Ping!Ping!Ping!"))
	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	rb := make([]byte, NETWORK_MTU)
	if _, err := conn.Read(rb); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			// We chose to be optimistic and treat lack of response (connection timeout) as an open port.
			return pr
		}

		pr.ReturnCode = Failure
		pr.Msg = fmt.Sprintf("Failed to write %s (%s): Connection refused", p.url.Host, "udp")
		return pr
	}

	return pr
}
