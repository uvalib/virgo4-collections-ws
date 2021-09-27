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

// DBConfig wraps up all of the DB configuration
type DBConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

func (sc *SolrConfig) selectURL() string {
	return fmt.Sprintf("%s/%s/select", sc.URL, sc.Core)
}

// ServiceConfig defines all of the JRML pool configuration parameters
type ServiceConfig struct {
	Port int
	Solr SolrConfig
	DB   DBConfig
}

// LoadConfiguration will load the service configuration from the commandline
// and return a pointer to it. Any failures are fatal.
func LoadConfiguration() *ServiceConfig {
	log.Printf("INFO: loading configuration...")
	var cfg ServiceConfig
	flag.IntVar(&cfg.Port, "port", 8080, "API service port (default 8080)")
	flag.StringVar(&cfg.Solr.URL, "solr", "", "Solr URL")
	flag.StringVar(&cfg.Solr.Core, "core", "test_core", "Solr core")

	// DB connection params
	flag.StringVar(&cfg.DB.Host, "dbhost", "localhost", "Database host")
	flag.IntVar(&cfg.DB.Port, "dbport", 5432, "Database port")
	flag.StringVar(&cfg.DB.Name, "dbname", "virgo4", "Database name")
	flag.StringVar(&cfg.DB.User, "dbuser", "v4user", "Database user")
	flag.StringVar(&cfg.DB.Pass, "dbpass", "pass", "Database password")

	flag.Parse()

	if cfg.Solr.URL == "" {
		log.Fatal("Parameter solr is required")
	}
	if cfg.Solr.Core == "" {
		log.Fatal("Parameter core is required")
	}

	if cfg.DB.Host == "" {
		log.Fatal("Parameter dbhost is required")
	}
	if cfg.DB.Name == "" {
		log.Fatal("Parameter dbname is required")
	}
	if cfg.DB.User == "" {
		log.Fatal("Parameter dbuser is required")
	}
	if cfg.DB.Pass == "" {
		log.Fatal("Parameter dbpass is required")
	}

	log.Printf("[CONFIG] port          = [%d]", cfg.Port)
	log.Printf("[CONFIG] solr          = [%s]", cfg.Solr.URL)
	log.Printf("[CONFIG] core          = [%s]", cfg.Solr.Core)
	log.Printf("[CONFIG] dbhost        = [%s]", cfg.DB.Host)
	log.Printf("[CONFIG] dbport        = [%d]", cfg.DB.Port)
	log.Printf("[CONFIG] dbname        = [%s]", cfg.DB.Name)
	log.Printf("[CONFIG] dbuser        = [%s]", cfg.DB.User)

	return &cfg
}
