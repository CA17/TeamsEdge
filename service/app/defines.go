package app

type JsonOptions struct {
	Id    interface{} `json:"id,string"`
	Value interface{} `json:"value"`
}

const (
	EdgeInformTask = "EdgeInformTask"
)
