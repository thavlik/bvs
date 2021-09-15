package mongodb_storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type MyStruct struct {
	Foo int `bson:"foo"`
}

func TestBSONMarshal(t *testing.T) {
	in := &MyStruct{Foo: 12345}
	ser, err := bson.Marshal(in)
	require.NoError(t, err)
	out := &MyStruct{}
	bson.Unmarshal(ser, out)
	require.Equal(t, in.Foo, out.Foo)
}
