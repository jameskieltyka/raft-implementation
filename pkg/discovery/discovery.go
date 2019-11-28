package discovery

type NodeLister interface {
	GetNodes(name string) []string
}
