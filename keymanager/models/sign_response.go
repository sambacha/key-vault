package models

// SignResponse is the vault sign response model.
type SignResponse struct {
	Data SignatureModel `json:"data"`
}

// SignatureModel represents vault signature model.
type SignatureModel struct {
	Signature string `json:"signature"`
}
