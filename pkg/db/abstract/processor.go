package abstract

type DbProcessor interface {
	Connect() error
	Disconnect() error
	Exec(q string) error
	Insert(m map[string]string) error
	Query(q [][]string) struct{}
	IsConnected() bool
}
