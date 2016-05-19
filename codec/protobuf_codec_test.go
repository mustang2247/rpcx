package codec

import (
	"testing"

	"github.com/smallnest/betterrpc"
)

type ProtoArith int

func (t *ProtoArith) Mul(args *ProtoArgs, reply *ProtoReply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *ProtoArith) Error(args *ProtoArgs, reply *ProtoReply) error {
	panic("ERROR")
}

func TestProtobufCodec(t *testing.T) {
	server := betterrpc.NewServer()
	server.ServerCodecFunc = NewProtobufServerCodec
	server.RegisterName(serviceName, new(ProtoArith))
	server.Start("tcp", "127.0.0.1:0")
	serverAddr := server.Address()

	s := &betterrpc.DirectClientSelector{Network: "tcp", Address: serverAddr}
	client := betterrpc.NewClient(s)
	client.ClientCodecFunc = NewProtobufClientCodec
	err := client.Start()
	if err != nil {
		t.Errorf("can't connect to %s because of %v \n", serverAddr, err)
	}

	args := &ProtoArgs{7, 8}
	var reply ProtoReply
	err = client.Call(serviceMethodName, args, &reply)
	if err != nil {
		t.Errorf("error for Arith: %d*%d, %v \n", args.A, args.B, err)
	} else {
		t.Logf("Arith: %d*%d=%d \n", args.A, args.B, reply.C)
	}

	client.Close()
}