package aircraftlib

import (
	"github.com/kumo-rpc/sharkgo"
)

var (
	// Make sure interface types satisfy TypeParam:
	_ shark.TypeParam[Echo] = Echo{}

	// ...and structs:
	_ shark.TypeParam[Zdate] = Zdate{}

	// ...and lists:
	_ shark.TypeParam[Echo_List]  = Echo_List{}
	_ shark.TypeParam[Zdate_List] = Zdate_List{}
)
