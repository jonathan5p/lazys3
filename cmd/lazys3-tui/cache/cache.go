package cache

type Cache[O any] struct {
	data map[string]O
}

func NewCache[O any]() *Cache[O] {
	return &Cache[O]{data: make(map[string]O)}
}

func (c *Cache[O]) Get(key string) (O, bool) {
	obj, ok := c.data[key]
	return obj, ok
}

func (c *Cache[O]) Set(key string, obj O) {
	c.data[key] = obj
}
