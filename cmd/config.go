package main

import (
	"flag"
	"fmt"
	"log"
)

// SolrConfig wraps up the config for solr acess
type SolrConfig struct {
	URL  string
	Core string
}

func (sc *SolrConfig) selectURL() string {
	return fmt.Sprintf("%s/%s/select", sc.URL, sc.Core)
}

// ServiceConfig defines all of the JRML pool configuration parameters
type ServiceConfig struct {
	Port int
	Solr SolrConfig
}

// LoadConfiguration will load the service configuration from the commandline
// and return a pointer to it. Any failures are fatal.
func LoadConfiguration() *ServiceConfig {
	log.Printf("INFO: loading configuration...")
	var cfg ServiceConfig
	flag.IntVar(&cfg.Port, "port", 8080, "API service port (default 8080)")
	flag.StringVar(&cfg.Solr.URL, "solr", "", "Solr URL")
	flag.StringVar(&cfg.Solr.Core, "core", "test_core", "Solr core")

	flag.Parse()

	if cfg.Solr.URL == "" {
		log.Fatal("Parameter solr is required")
	}
	if cfg.Solr.Core == "" {
		log.Fatal("Parameter core is required")
	}

	log.Printf("[CONFIG] port          = [%d]", cfg.Port)
	log.Printf("[CONFIG] solr          = [%s]", cfg.Solr.URL)
	log.Printf("[CONFIG] core          = [%s]", cfg.Solr.Core)

	return &cfg
}
