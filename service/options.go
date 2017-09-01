package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

const (
	DefaultServer         = "localhost:2181"
	DefaultConnectTimeout = 5 * time.Second
	DefaultSessionTimeout = 1 * time.Hour
	DefaultPath           = "/xmidt"
	DefaultServiceName    = "test"
	DefaultRegistration   = "http://localhost:8080"
	DefaultVnodeCount     = 211
)

// Options represents the set of configurable attributes for service discovery and registration
type Options struct {
	// Logger is used by any component configured via this Options.  If unset, a default
	// logger is used.
	Logger log.Logger `json:"-"`

	// Connection is the comma-delimited Zookeeper connection string.  Both this and
	// Servers may be set, and they will be merged together when connecting to Zookeeper.
	Connection string `json:"connection,omitempty"`

	// Servers is the array of Zookeeper servers.  Both this and Connection may be set,
	// and they will be merged together when connecting to Zookeeper.
	Servers []string `json:"servers,omitempty"`

	// ConnectTimeout is the Zookeeper connection timeout.
	ConnectTimeout time.Duration `json:"connectTimeout"`

	// SessionTimeout is the Zookeeper session timeout.
	SessionTimeout time.Duration `json:"sessionTimeout"`

	// UpdateDelay specifies the period of time between a service discovery update and when a client
	// is notified.  Updates during the wait time simply replace the waiting set of instances.
	// There is no default for this field.  If unset, all updates are immediately processed.
	UpdateDelay time.Duration `json:"updateDelay"`

	// Path is the base path for all znodes created via this Options.
	Path string `json:"path,omitempty"`

	// ServiceName is the name of the service being registered.
	ServiceName string `json:"serviceName,omitempty"`

	// Registration is the data stored about this service, typically host:port or scheme://host:port.
	Registration string `json:"registration,omitempty"`

	// VnodeCount is used to tune the underlying consistent hash algorithm for servers.
	VnodeCount uint `json:"vnodeCount"`
}

func (o *Options) String() string {
	if o == nil {
		return "<nil>"
	}

	if data, err := json.Marshal(o); err == nil {
		return string(data)
	} else {
		return err.Error()
	}
}

func (o *Options) logger() log.Logger {
	if o != nil && o.Logger != nil {
		return o.Logger
	}

	return log.NewNopLogger()
}

func (o *Options) servers() []string {
	servers := make([]string, 0, 10)

	if o != nil {
		if len(o.Connection) > 0 {
			for _, server := range strings.Split(o.Connection, ",") {
				servers = append(servers, strings.TrimSpace(server))
			}
		}

		if len(o.Servers) > 0 {
			servers = append(servers, o.Servers...)
		}
	}

	if len(servers) == 0 {
		servers = append(servers, DefaultServer)
	}

	return servers
}

func (o *Options) connectTimeout() time.Duration {
	if o != nil && o.ConnectTimeout > 0 {
		return o.ConnectTimeout
	}

	return DefaultConnectTimeout
}

func (o *Options) sessionTimeout() time.Duration {
	if o != nil && o.SessionTimeout > 0 {
		return o.SessionTimeout
	}

	return DefaultSessionTimeout
}

func (o *Options) updateDelay() time.Duration {
	if o != nil && o.UpdateDelay > 0 {
		return o.UpdateDelay
	}

	return 0
}

func (o *Options) path() string {
	if o != nil && len(o.Path) > 0 {
		return o.Path
	}

	return DefaultPath
}

func (o *Options) serviceName() string {
	if o != nil && len(o.ServiceName) > 0 {
		return o.ServiceName
	}

	return DefaultServiceName
}

func (o *Options) registration() string {
	if o != nil && len(o.Registration) > 0 {
		return o.Registration
	}

	return DefaultRegistration
}

func (o *Options) vnodeCount() int {
	if o != nil && o.VnodeCount > 0 {
		return int(o.VnodeCount)
	}

	return DefaultVnodeCount
}
