package defaultservice

import (
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
)

// FulfillOrder -
func (s *Service) FulfillOrder(req *api.FulfillOrderRequest) (*api.FulfillOrderResponse, error) {
	orderPart1CID := req.OrderPart1CID
	nodeID := s.NodeID()
	remoteIDDocCID := req.DocumentCID
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return nil, err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return nil, err
	}

	//Retrieve the order from IPFS
	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart1CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return nil, err
	}

	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, nodeID)
	if err != nil {
		return nil, err
	}

	//Generate the secret and store for later redemption
	seed, err := common.MakeRandomSeedAndStore(s.Store, s.Rng, order.Reference)
	if err != nil {
		return nil, err
	}

	//Generate the Public Key (Commitment) from the Seed/Secret
	commitmentPublicKey, err := cryptowallet.RedeemPublicKey(seed)
	if err != nil {
		return nil, err
	}

	//Create an order response in IPFS
	orderPart2CID, err := common.CreateAndStoreOrderPart2(s.Ipfs, s.Store, order, orderPart1CID, commitmentPublicKey, nodeID, recipientList)
	if err != nil {
		return nil, err
	}

	return &api.FulfillOrderResponse{
		OrderPart2CID: orderPart2CID,
	}, nil
}

// FulfillOrderSecret -
func (s *Service) FulfillOrderSecret(req *api.FulfillOrderSecretRequest) (*api.FulfillOrderSecretResponse, error) {
	//Initialise values from Request object
	orderPart3CID := req.OrderPart3CID
	nodeID := s.NodeID()
	remoteIDDocCID := req.SenderDocumentCID
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return nil, err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return nil, err
	}

	//Retrieve the order from IPFS
	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart3CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return nil, err
	}

	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, nodeID)
	if err != nil {
		return nil, err
	}

	//Retrieve the Seed
	seed, err := common.RetrieveSeed(s.Store, order.Reference)
	if err != nil {
		return nil, err
	}

	//Generate the Secert from the Seed
	commitmentPrivateKey, err := cryptowallet.RedeemSecret(seed)
	if err != nil {
		return nil, err
	}

	//Create an order response in IPFS
	orderPart4ID, err := common.CreateAndStoreOrderPart4(s.Ipfs, s.Store, order, commitmentPrivateKey, orderPart3CID, nodeID, recipientList)
	if err != nil {
		return nil, err
	}

	return &api.FulfillOrderSecretResponse{
		OrderPart4CID: orderPart4ID,
	}, nil
}
