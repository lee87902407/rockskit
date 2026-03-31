package rockskit

type OpType string

const (
	OpPut    OpType = "put"
	OpDelete OpType = "delete"
)

type KVOp struct {
	Type  OpType
	Key   []byte
	Value []byte
}

type KeyValue struct {
	Key   []byte
	Value []byte
}
