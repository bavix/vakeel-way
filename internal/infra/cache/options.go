package cache

import "time"

// Option is a function that can be used to configure a Cache instance.
//
// It takes a pointer to a Cache instance and returns nothing.
// The purpose of the Option function is to provide a way to configure
// a Cache instance with different options, such as setting the maximum
// number of items in the cache or setting the onEvict function.
//
// Parameters:
//   - c: A pointer to the Cache struct that will be customized.
//
// Returns:
//
//	None.
type Option[K comparable, V any] func(c *Cache[K, V])

// WithOnEvict returns an Option that sets the onEvict function for the cache.
//
// The onEvict function is called when an item is evicted from the cache. It takes
// the key of the evicted item as a parameter, and can be used to perform an action
// when an item is evicted, such as logging the eviction of an item.
//
// The onEvict function is called after the item has been removed from the cache, and
// is not called if the item is not evicted from the cache (e.g. if the cache is
// cleared).
//
// Parameters:
//   - onEvict: The function to be called when an item is evicted from the cache.
//
// Returns:
//   - An Option that sets the onEvict function for the cache.
func WithOnEvict[K comparable, V any](onEvict Fn[K, V]) Option[K, V] {
	// Return an Option function that sets the onEvict function in the cache.
	return func(c *Cache[K, V]) {
		// Set the onEvict function in the cache to the provided onEvict function.
		// The onEvict function is called when an item is evicted from the cache.
		// It takes the key of the evicted item as a parameter.
		// The purpose of the onEvict function is to provide a way to perform an action
		// when an item is evicted from the cache, such as logging the eviction of an item.
		//
		// Parameters:
		// - c: A pointer to the Cache struct that will be customized.
		//
		// Returns:
		// None.
		c.onEvict = onEvict
	}
}

// WithClock returns an Option that sets the clock for the cache.
//
// The clock is used to get the current time, which is used to determine when
// items in the cache have expired. The clock can be used to control the behavior
// of the cache, such as setting the initial time for the cache, or setting the
// maximum time-to-live for items in the cache.
//
// Parameters:
//   - clock: The clock to use for the cache. It must implement the Click interface.
//
// Returns:
//
//	An Option that sets the clock for the cache.
func WithClock[K comparable, V any](clock Click) Option[K, V] {
	return func(c *Cache[K, V]) {
		// Set the clock in the cache to the provided clock.
		// The clock is used to get the current time, which is used to determine when
		// items in the cache have expired.
		//
		// This is useful for testing, as it allows you to control the behavior of the cache
		// by setting the current time to a specific value.
		//
		// Parameters:
		// - c: A pointer to the Cache struct that will be customized.
		//
		// Returns:
		// None.
		c.clock = clock
	}
}

// WithEvictDuration is an option that sets the eviction duration for the cache.
//
// The eviction duration is the time after which an item is evicted from the cache.
// The parameter evictDuration specifies the duration after which an item is evicted
// from the cache. The default value for evictDuration is 1 minute.
//
// Parameters:
//   - evictDuration: The duration after which an item is evicted from the cache.
//
// Returns:
//
//	An option that sets the eviction duration for the cache.
func WithEvictDuration[K comparable, V any](evictDuration time.Duration) Option[K, V] {
	// The WithEvictDuration function returns an Option that sets the eviction duration
	// for the cache. The eviction duration is the time after which an item is evicted
	// from the cache. This function takes an evictDuration parameter, which specifies
	// the duration after which an item is evicted from the cache. The default value
	// for evictDuration is 1 minute.
	//
	// The returned Option is a function that takes a pointer to a Cache struct and
	// sets the eviction duration for the cache.
	return func(c *Cache[K, V]) {
		// Set the eviction duration in the cache to the provided evictDuration.
		// The eviction duration is the time after which an item is evicted from the cache.
		//
		// This is useful for controlling the behavior of the cache, such as setting
		// the maximum time-to-live for items in the cache.
		//
		// Parameters:
		// - c: A pointer to the Cache struct that will be customized.
		//
		// Returns:
		// None.
		c.evictDuration = evictDuration
	}
}
