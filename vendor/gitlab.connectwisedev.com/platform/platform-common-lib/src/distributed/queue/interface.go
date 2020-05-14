package queue

// Interface - Interface to hold Queue data
type Interface interface {
	GetList(queueName string) ([]string, error)
	CreateItem(data []byte, queueName string) (string, error)
	GetItemData(queueName, itemName string) ([]byte, error)
	RemoveItem(queueName string, itemName string) error
}
