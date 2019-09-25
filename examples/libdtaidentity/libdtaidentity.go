package main

import (
	"C"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
)

//export createIdentity
func createIdentity(cname *C.char) *C.char {
    ipfsNode, err := ipfs.NewMemoryConnector()
    if err != nil {
        panic(err);
    }
    name := C.GoString(cname);
    id,_ := identity.CreateIdentity(name,ipfsNode,nil);
    return C.CString(id);
}

func main() {}
