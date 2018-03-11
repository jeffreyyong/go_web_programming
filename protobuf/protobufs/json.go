package main

import (
	"encoding/json"
	"fmt"
	"web/protobuf_grpc/protofiles"
)

func main() {
	p := &protofiles.Person{
		Id:    1234,
		Name:  "Jeffrey Yong",
		Email: "jeffreyyong10@gmail.com",
		Phones: []*protofiles.Person_PhoneNumber{
			{Number: "1234", Type: protofiles.Person_MOBILE},
		},
	}
	body, _ := json.Marshal(p)
	fmt.Println(string(body))
}

// Benefits of protobuf: Intended for two backend systems to communicate with each other with less overhead
// Since the size of the binary is less than text, protocol marshaled data is of less size than JSON
