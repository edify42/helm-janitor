package container

// Keypair is a simple key + value combo
type Keypair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// KeyArray is an array of the above
type KeyArray []Keypair
