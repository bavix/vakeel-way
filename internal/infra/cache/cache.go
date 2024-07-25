package cache

import (
	"sync"
	"time"
)

// Fn is a function type that represents a callback function.
//
// It is used as the OnEvict parameter of the Cache struct. The OnEvict parameter
// is a function that is called when an item is evicted from the cache. The function
// takes the key and value of the evicted item as parameters.
//
// OnEvict functions are useful for performing actions when an item is evicted from
// the cache. For example, you can log the eviction of an item or perform additional
// cleanup operations.
//
// Parameters:
//   - key: The key of the evicted item.
//   - value: The value of the evicted item.
//
// Example:
//
//	// Log the eviction of an item
//	func itemEvicted(key string, value int) {
//	    log.Printf("Item with key %s and value %d evicted from cache", key, value)
//	}
//
//	cache := cache.NewCache[string, int](10, cache.WithOnEvict(itemEvicted))
type Fn[K comparable, V any] func(key K, value V)

// Cache is a thread-safe cache implementation that stores key-value pairs with a time-to-live (TTL)
// for each item. It is implemented as a map where the keys are strings and the values are pointers
// to item structs. The cache has a maximum size, which is specified by the maxSize parameter.
// The clock parameter is an interface that provides the current time. The onEvict parameter is a
// function that is called when an item is evicted from the cache. The evictDuration parameter
// specifies the duration after which an item is evicted from the cache.
type Cache[K comparable, V any] struct {
	// items is a map that stores the key-value pairs of the cache. The keys are the strings that
	// are used to identify the items in the cache, and the values are pointers to item structs,
	// which contain the value associated with the key and the time-to-live (TTL) of the item.
	items map[K]*item[V]

	// clock is an interface that provides the current time. It is used to get the current time
	// and calculate the expiration time of the items in the cache. The current time is used to
	// calculate the expiration time of the items in the cache. The expiration time of an item is
	// the sum of the current time and the TTL of the item.
	clock Click

	// onEvict is a function that is called when an item is evicted from the cache. It takes a
	// string parameter, which is the key of the evicted item. The purpose of the onEvict
	// parameter is to provide a way to perform an action when an item is evicted from the cache.
	// For example, the onEvict parameter can be used to log the eviction of an item.
	onEvict Fn[K, V]

	// evictDuration is the duration after which an item is evicted from the cache. It specifies
	// the time interval after which an item is considered expired and is evicted from the
	// cache. The evictDuration is used to calculate the expiration time of the items in the cache.
	// The expiration time of an item is the sum of the current time and the TTL of the item.
	evictDuration time.Duration

	// mu is a sync.RWMutex that is used to synchronize access to the cache. It is used to ensure
	// that only one goroutine can modify the cache at a time. The mu is used to protect the
	// cache from concurrent modifications. The mu is used to ensure that the cache is accessed
	// and modified in a thread-safe way.
	mu sync.RWMutex
}

// item is a struct that represents an item stored in the cache.
//
// It contains the key used to identify the item in the cache, the value associated with the key,
// and the time-to-live (TTL) of the item.
type item[V any] struct {
	// Value is the value associated with the key.
	//
	// It is the value that is stored in the cache and can be retrieved using the key.
	Value V

	// TTL is the time-to-live (TTL) of the item. The item will be automatically
	// removed from the cache after the TTL has expired.
	//
	// It represents the duration after which the item should be removed from the cache.
	// The TTL is calculated based on the current time and the evictDuration parameter of the Cache struct.
	TTL time.Time
}

// NewCacheWithOptions creates a new Cache instance with the specified minimum capacity and optional configurations.
//
// The function initializes a new Cache instance with the specified minimum capacity and default values.
// It takes the minimumCapacity parameter, which specifies the minimum capacity of the cache.
// The options parameter is a variadic list of Option functions that can be used to configure the cache.
//
// Parameters:
//   - minimumCapacity: The minimum capacity of the cache.
//   - options: Optional configurations for the cache.
//
// Returns:
//   - A pointer to the initialized Cache instance.
//
//nolint:exhaustruct
func NewCache[K comparable, V any](minimumCapacity int, options ...Option[K, V]) *Cache[K, V] {
	// Initialize a new Cache instance with the specified minimum capacity and default values.
	cache := &Cache[K, V]{
		// Create a map to store the items in the cache. The map has a minimum capacity specified by the minimumCapacity parameter.
		items: make(map[K]*item[V], minimumCapacity),
		// Use the default clock implementation.
		clock: clock{},
		// Set the default onEvict function to do nothing.
		onEvict: func(K, V) {},
		// Set the default evict duration to 1 minute.
		evictDuration: time.Minute,
	}

	// Apply any optional configurations provided through the options parameter.
	for _, option := range options {
		// Apply the configuration to the cache.
		option(cache)
	}

	// Start a cleanup goroutine for the cache.
	// The cleanup goroutine periodically removes expired items from the cache.
	go cache.cleanup()

	// Return the initialized Cache instance.
	return cache
}

// Get retrieves the value from the cache associated with the given key.
//
// It returns a pointer to the value and a boolean indicating whether the key was found in the cache.
// If the key is not found, the pointer will be nil and the boolean will be false.
//
// Parameters:
//   - key: The key used to identify the item in the cache.
//
// Returns:
//   - value: A pointer to the value associated with the key.
//   - found: A boolean indicating whether the key was found in the cache.
//
// This function locks the cache for read access to prevent concurrent modifications.
// It retrieves the value associated with the given key from the cache.
// If the key is found in the cache, it returns a pointer to the value and true.
// If the key is not found in the cache, it returns nil and false.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	// Lock the cache for read access.
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if the key exists in the cache.
	if v, ok := c.items[key]; ok {
		// Return a pointer to the value and indicate that the key was found.
		return &v.Value, true
	}

	// Return nil and false if the key was not found.
	return nil, false
}

// Add adds a new item to the cache with the given key, value, and time-to-live (TTL).
//
// If the key already exists in the cache, its value and TTL are updated.
//
// Parameters:
//   - key: The key used to identify the item in the cache.
//   - value: The value associated with the key.
//   - ttl: The time-to-live (TTL) of the item. The item will be automatically removed
//     from the cache after the TTL has expired.
//
// This function locks the cache for write access, creates a new cache item with the given key, value, and TTL,
// and adds it to the cache. If the key already exists in the cache, its value and TTL are updated.
func (c *Cache[K, V]) Add(key K, value V, ttl time.Duration) {
	// Lock the cache for write access.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a new item with the given key, value, and TTL.
	// The item struct contains the value associated with the key and the time-to-live (TTL) of the item.
	item := &item[V]{
		// The value associated with the key.
		Value: value, // Associate the value with the key.
		// The time-to-live (TTL) of the item. The item will be automatically removed
		// from the cache after the TTL has expired.
		TTL: c.clock.Now().Add(ttl), // Set the time-to-live (TTL) for the item.
	}

	// Add the new item to the cache with the given key.
	// If the key already exists in the cache, its value and TTL are updated.
	c.items[key] = item // Add or update the item in the cache.
}

// OnEvict sets a callback function that will be called when an item is evicted
// from the cache. The callback function takes the key of the evicted item as
// a parameter.
//
// Parameters:
//   - fn: The callback function that will be called when an item is evicted from
//     the cache. The callback function takes the key of the evicted item as a
//     parameter.
//
// This function sets the onEvict callback function for the cache. The onEvict
// callback function is called when an item is evicted from the cache. It takes
// the key of the evicted item as a parameter. The onEvict function is useful for
// performing actions when an item is evicted from the cache, such as logging the
// eviction of an item.
//
// The onEvict function is executed after the default onEvict function, which
// removes the item from the cache. If an onEvict function is already set, the new
// function is composed with the old onEvict function. The composed function
// calls the provided callback function and then calls the old onEvict function.
//
// Parameters:
//   - fn: The callback function that will be called when an item is evicted from
//     the cache. The callback function takes the key of the evicted item as a
//     parameter.
func (c *Cache[K, V]) OnEvict(fn func(K, V)) {
	// Lock the cache for write access.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save the old onEvict function.
	old := c.onEvict

	// Create a new onEvict function that calls the provided callback function
	// and then calls the old onEvict function.
	c.onEvict = func(key K, value V) {
		// Call the provided callback function with the key of the evicted item.
		fn(key, value)

		// Call the old onEvict function with the key of the evicted item.
		if old != nil {
			old(key, value)
		}
	}
}

// cleanup is a goroutine that periodically removes expired items from the cache.
//
// The cleanup process is triggered by a ticker, which ticks every evictDuration.
// The removeExpiredItems function is called periodically to remove the expired
// items from the cache.
//
// The cleanup goroutine runs indefinitely until the program is terminated. It
// is started when the Cache instance is created.
//
// cleanup is a goroutine, meaning it runs concurrently with other goroutines in
// the program. It is started when the Cache instance is created and it runs until
// the program is terminated.
func (c *Cache[K, V]) cleanup() {
	// Create a ticker that ticks every evictDuration.
	// The ticker is used to schedule the cleanup process.
	ticker := time.NewTicker(c.evictDuration)

	// Ensure that the ticker is stopped even if the function returns early.
	defer ticker.Stop()

	// Run the cleanup loop until the program is terminated.
	// The loop runs indefinitely until the program is terminated.
	for range ticker.C {
		// Remove expired items from the cache.
		// This function is called periodically by the cleanup goroutine.
		c.removeExpiredItems()
	}
}

// removeExpiredItems removes the expired items from the cache.
//
// This function is called periodically by the cleanup goroutine to remove the expired items from the cache.
// It locks the cache for write access and then iterates over each item in the cache. For each item, it checks
// if the item has expired by comparing the item's TTL (time-to-live) with the current time. If an item has
// expired, it calls the onEvict function with the key of the expired item. The onEvict function is a
// callback function that is called when an item is evicted from the cache. It is used to perform an action
// when an item is evicted, such as logging the eviction of an item. Finally, it removes the expired item from
// the cache.
func (c *Cache[K, V]) removeExpiredItems() {
	// Lock the cache for write access.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Iterate over each item in the cache.
	for k := range c.items {
		// Get the item from the cache. If the item is nil, skip it.
		item := c.items[k]
		if item == nil {
			continue
		}

		// Check if the item has expired.
		// An item is considered expired if its TTL (time-to-live) is before the current time.
		if item.TTL.Before(c.clock.Now()) {
			// Call the onEvict function with the key of the expired item.
			// The onEvict function is a callback function that is called when an item is evicted from the cache.
			// It is used to perform an action when an item is evicted, such as logging the eviction of an item.
			if c.onEvict != nil {
				c.onEvict(k, item.Value)
			}

			// Remove the expired item from the cache.
			// The delete function removes the item with the given key from the cache.
			// It does not return any value.
			delete(c.items, k)
		}
	}
}
