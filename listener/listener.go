package listener

type Empty struct{}

var client = map[string]map[chan []byte]Empty{}

func NewClient(namespace string) chan []byte {
	ch := make(chan []byte, 100)
	clientSpace := client[namespace]
	if clientSpace == nil {
		clientSpace = map[chan []byte]Empty{}
		client[namespace] = clientSpace
	}
	clientSpace[ch] = Empty{}
	return ch
}

func CloseClient(namespace string, ch chan []byte) {
	delete(client[namespace], ch)
	close(ch)
}

func SendMessage(namespace string, message []byte) {
	for ch, _ := range client[namespace] {
		ch <- message
	}
}
