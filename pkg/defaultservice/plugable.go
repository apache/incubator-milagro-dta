package defaultservice

import (
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
)

// Plugable service methods
type Plugable interface {
	// service
	Name() string
	Vendor() string

	// order
	ValidateOrderRequest(req *api.OrderRequest) error
	PrepareOrderPart1(order *documents.OrderDoc, reqExtension map[string]string) (fulfillExtension map[string]string, err error)
	PrepareOrderResponse(orderPart2 *documents.OrderDoc, reqExtension, fulfillExtension map[string]string) (commitment string, extension map[string]string, err error)
	ProduceBeneficiaryEncryptedData(blsSK []byte, order *documents.OrderDoc, req *api.OrderSecretRequest) (encrypted []byte, extension map[string]string, err error)
	ProduceFinalSecret(seed, sikeSK []byte, order, orderPart4 *documents.OrderDoc, req *api.OrderSecretRequest, fulfillSecretRespomse *api.FulfillOrderSecretResponse) (secret, commitment string, extension map[string]string, err error)
}
