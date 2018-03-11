package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"

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

	p1 := &protofiles.Person{}
	body, _ := proto.Marshal(p)
	proto.Unmarshal(body, p1)
	fmt.Println("Original struct loaded from proto file:", p, "\n")
	fmt.Println("Marshaled proto data: ", body, "\n")
	fmt.Println("Unmarshaled struct: ", p1)
}
