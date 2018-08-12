package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

var (
	allowedTlsProtocols = map[string]uint16{
		"SSL3.0": tls.VersionSSL30,
		"TLS1.0": tls.VersionTLS10,
		"TLS1.1": tls.VersionTLS11,
		"TLS1.2": tls.VersionTLS12,
	}
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"keyspace": resourceKeyspace(),
		},
		ConfigureFunc: configureProvider,
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cassandra Username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cassandra Password",
				Sensitive:   true,
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Default:     9042,
				Description: "Cassandra CQL Port",
				ValidateFunc: func(i interface{}, s string) (ws []string, errors []error) {
					port := i.(int)

					if port <= 0 || port >= 65535 {
						errors = append(errors, fmt.Errorf("%d: invalid value - must be between 1 and 65535", port))
					}

					return
				},
			},
			"cqlVersion": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Default:     "3.4.4",
				Description: "CQL Version",
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				MinItems: 1,
				Required: true,
			},
			"connectionTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Default:     1000,
				Description: "Connection timeout in milliseconds",
			},
			"rootCA": &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Description: "Use root CA to connect to Cluster. Applies only when useSSL is enabled",
				ValidateFunc: func(i interface{}, s string) (ws []string, errors []error) {
					rootCA := i.(string)

					if rootCA == "" {
						return
					}

					caPool := x509.NewCertPool()
					ok := caPool.AppendCertsFromPEM([]byte(rootCA))

					if !ok {
						errors = append(errors, fmt.Errorf("%s: invalid PEM", rootCA))
					}

					return
				},
			},
			"useSSL": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Default:     false,
				Description: "Use SSL when connecting to cluster",
			},
			"minTLSVersion": &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Default:     "TLS1.2",
				Description: "Minimum TLS Version used to connect to the cluster - allowed values are SSL3.0, TLS1.0, TLS1.1, TLS1.2. Applies only when useSSL is enabled",
				ValidateFunc: func(i interface{}, s string) (ws []string, errors []error) {
					minTlsVersion := i.(string)

					if allowedTlsProtocols[minTlsVersion] == 0 {
						errors = append(errors, fmt.Errorf("%s: invalid value - must be one of SSL3.0, TLS1.0, TLS1.1, TLS1.2", minTlsVersion))
					}

					return
				},
			},
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	useSSL := d.Get("useSSL").(bool)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	port := d.Get("port").(int)
	connectionTimeout := d.Get("connectionTimeout").(int)
	cqlVersion := d.Get("cqlVersion").(string)
	hosts := d.Get("hosts").([]string)

	cluster := gocql.NewCluster()

	cluster.Hosts = hosts

	cluster.Port = port

	cluster.Consistency = gocql.All

	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	cluster.ConnectTimeout = time.Millisecond * time.Duration(connectionTimeout)

	cluster.CQLVersion = cqlVersion

	cluster.Keyspace = "system"

	if useSSL {

		rootCA := d.Get("rootCA").(string)
		minTLSVersion := d.Get("minTLSVersion").(string)

		tlsConfig := &tls.Config{
			MinVersion: allowedTlsProtocols[minTLSVersion],
		}

		if rootCA != "" {
			caPool := x509.NewCertPool()
			ok := caPool.AppendCertsFromPEM([]byte(rootCA))

			if !ok {
				return nil, errors.New("Unable to load rootCA")
			}

			tlsConfig.RootCAs = caPool
		}

		cluster.SslOpts = &gocql.SslOptions{
			Config: tlsConfig,
		}
	}

	return cluster.CreateSession()

}
