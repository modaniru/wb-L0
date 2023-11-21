package storage

//Sync
type Cache struct{
	mapa map[string][]byte
}

func NewCache() *Cache{
	return &Cache{mapa: make(map[string][]byte)}
}

func (c *Cache) Get(key string) ([]byte, bool){
	value, ok :=  c.mapa[key]
	return value, ok
}

func (c *Cache) Put(key string, data []byte){
	c.mapa[key] = data
}