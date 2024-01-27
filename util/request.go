package util

import (
	"net/http"
)

type Request struct {
	State State `json:"state"`
}

func NewRequest(r *http.Request) Request {
	return Request{
		State: State{
			Config: Config{
				CfAddMeta: "test",
			},
		},
	}
}

func FromMap(m map[string]any) Request {
	return MapToStruct[Request](m)
}

func (r Request) ToMap() map[string]any {
	return StructToMap(r)
}

type State struct {
	Config Config `json:"config"`
}

type Config struct {
	CfAddMeta string `json:"cf_add_meta"`
}
