package sign

import "encoding/json"

type CryptoMaterial struct {
	Type       string `json:"type"`
	Pin        string `json:"pin"`
	Identifier string `json:"identifier"`
	SlotId     *int   `json:"slotId,omitempty"`
	TokenLabel string `json:"tokenLabel,omitempty"`
}

type SignRequest struct {
	Bundle            json.RawMessage `json:"bundle"`
	Provenance        json.RawMessage `json:"provenance"`
	CertChain         []string        `json:"certChain"`
	Timestamp         int64           `json:"timestamp"`
	Strategy          string          `json:"strategy"`
	PolicyId          string          `json:"policyId"`
	CryptoMaterial    CryptoMaterial  `json:"cryptoMaterial"`
	OperationalConfig json.RawMessage `json:"operationalConfig"`
}

type Options struct {
	BundlePath     string
	ProvenancePath string
	CertChainPath  string
	ConfigPath     string
	Timestamp      int64
	Strategy       string
	PolicyId       string
	CryptoType     string
	Pin            string
	Identifier     string
	SlotId         int
	SlotIdSet      bool
	TokenLabel     string
	ServiceURL     string
	OutputPath     string
	TimeoutSeconds int
}
