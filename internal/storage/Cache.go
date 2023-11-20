package storage

//Sync
type Cache struct{
	mapa map[string][]byte
}

func NewCache() *Cache{
	return &Cache{mapa: make(map[string][]byte)}
}

func (c *Cache) Get(key string) []byte{
	return c.mapa[key]
}

func (c *Cache) Put(key string, data []byte){
	c.mapa[key] = data
}