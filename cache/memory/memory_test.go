package memory

func createTestCache() *MemoryCache {
	return New(20 * 1024 * 1024)
}
