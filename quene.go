package httpbot

func NewQuene(size int) *ResponseQuene {
	return &ResponseQuene{
		nodes: make([]ResponseReader, size),
		size:  size,
	}
}

type ResponseQuene struct {
	nodes []ResponseReader
	size  int
	head  int
	tail  int
	count int
}

func (q *ResponseQuene) Push(cb ResponseReader) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]ResponseReader, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = cb
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

func (q *ResponseQuene) Pop() ResponseReader {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}
