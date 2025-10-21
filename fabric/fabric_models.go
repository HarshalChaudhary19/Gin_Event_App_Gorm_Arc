package fabric

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Basically Password hashed form mein hoga store but retrive nhi hoga
	Age      int    `json:"age"`
}

type UserUpdate struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password,omitempty"`
}

type UserReturn struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type UserAssetHistory struct {
	TxID            string          `json:"txId"`
	Timestamp       FabricTimestamp `json:"timestamp"`
	TimestampString string          `json:"-"` //Optional hai ye
	IsDelete        bool            `json:"isDelete"`
	Value           UserReturn      `json:"value"`
}

type FabricTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int64 `json:"nanos"`
}
