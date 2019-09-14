package defaultservice

import (
	"io"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/config"
)

// ServiceOption function to set Service properties
type ServiceOption func(s *Service) error

// Init sets-up the service options. It's called when the plugin gets registered
func (s *Service) Init(plugin Plugable, options ...ServiceOption) error {
	s.Plugin = plugin

	for _, opt := range options {
		if err := opt(s); err != nil {
			return err
		}
	}

	return nil
}

// WithLogger adds logger to the Service
func WithLogger(logger *logger.Logger) ServiceOption {
	return func(s *Service) error {
		s.Logger = logger
		return nil
	}
}

// WithRng adds rng to the Service
func WithRng(rng io.Reader) ServiceOption {
	return func(s *Service) error {
		s.Rng = rng
		return nil
	}
}

// WithStore adds store to the Service
func WithStore(store *datastore.Store) ServiceOption {
	return func(s *Service) error {
		s.Store = store
		return nil
	}
}

// WithIPFS adds ipfs connector to the Service
func WithIPFS(ipfsConnector ipfs.Connector) ServiceOption {
	return func(s *Service) error {
		s.Ipfs = ipfsConnector
		return nil
	}
}

// WithMasterFiduciary adds master fiduciary connector to the Service
func WithMasterFiduciary(masterFiduciaryServer api.ClientService) ServiceOption {
	return func(s *Service) error {
		s.MasterFiduciaryServer = masterFiduciaryServer
		return nil
	}
}

// WithConfig adds config settings to the Service
func WithConfig(cfg *config.Config) ServiceOption {
	return func(s *Service) error {
		s.SetNodeID(cfg.Node.NodeID)
		s.SetMasterFiduciaryNodeID(cfg.Node.MasterFiduciaryNodeID)
		return nil
	}
}
