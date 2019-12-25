package data

type Data struct {
	Result     interface{} `json: result`
	Infomation Infomation  `json: infomation`
}

type Infomation struct {
	Profile string `json: profile`
	Region  string `json: region`
}

type Infomations []Infomation
