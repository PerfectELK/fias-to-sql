package abstract

type DbProcessor interface {
	Connect() error
	Disconnect() error
	Use(q string) error
	Exec(q string) error
	Insert(m map[string]string) error
	Table(t string) DbProcessor
	Select(s []string) DbProcessor
	Where(q [][]string) DbProcessor
	Get() map[string]string
	IsConnected() bool
}
