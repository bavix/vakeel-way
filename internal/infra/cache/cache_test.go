package cache_test

import (
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/bavix/vakeel-way/internal/infra/cache"
)

// CacheTestSuite represents the test suite for the cache functionality.
//
// It extends the suite.Suite struct from the testify package to provide
// additional testing functionality.
type CacheTestSuite struct {
	// Suite is the base test suite from the testify package.
	// It provides functionality for running tests, checking assertions,
	// and reporting test results.
	suite.Suite
	// cache is an instance of the Cache struct used for testing.
	// It is initialized in the SetupTest method.
	cache *cache.Cache[int, string]
}

// SetupTest initializes the cache for testing.
//
// The method is called before each test in the test suite.
// It creates a new instance of the Cache struct with a maximum size of 10 and
// an evict duration of 100 milliseconds.
func (suite *CacheTestSuite) SetupTest() {
	suite.cache = cache.NewCache(
		10,
		cache.WithEvictDuration[int, string](100*time.Microsecond),
	)
}

// AfterTest is called after each test in the test suite.
//
// The method is part of the suite.Suite interface from the testify package.
// It is used to clean up resources after a test has finished running.
//
// In this case, the method is used to set the cache instance to nil to prevent
// memory leaks.
func (suite *CacheTestSuite) AfterTest(_, _ string) {
	// Set the cache instance to nil to prevent memory leaks.
	suite.cache = nil
}

// TestCache_Get tests the Get method of the Cache struct.
//
// The test verifies that an item can be retrieved from the cache using the Get method.
// The test adds an item to the cache with a key of 1 and a value of "hello".
// It then retrieves the item from the cache using the Get method and checks if the item was retrieved successfully.
// If the item was retrieved successfully, the test passes.
// If the item was not retrieved successfully, the test fails.
//
// This test is part of the CacheTestSuite test suite.
func (suite *CacheTestSuite) TestCache_Get() {
	// Add an item to the cache with a key of 1 and a value of "hello".
	// The Add method is used to add an item to the cache.
	suite.cache.Add(1, "hello", time.Second)

	// Get an item from the cache using the key 1.
	// The Get method is used to retrieve an item from the cache.
	item, ok := suite.cache.Get(1)

	// Check if the item was retrieved successfully.
	// If the item was retrieved successfully, the test passes.
	// If the item was not retrieved successfully, the test fails.
	suite.True(ok, "Item was not retrieved successfully")

	// Check if the retrieved item is not nil.
	suite.NotNil(item, "Retrieved item is nil")

	// Check if the retrieved item has the correct value.
	// The Equal method is used to compare the retrieved item with the expected value.
	suite.Equal("hello", *item, "Retrieved item has incorrect value")
}

// TestCache_Expire tests the expiration of items in the cache.
// It demonstrates that items added to the cache with a time-to-live (TTL)
// expire after the specified duration.
//
// This test adds an item to the cache with a TTL of 100 microseconds.
// It then waits for 200 milliseconds to let the item expire.
// After the waiting period, it retrieves the item from the cache and checks
// if it is still retrievable. Since the item has expired, it should not be
// retrievable anymore.
//
// The test is part of the CacheTestSuite test suite.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (suite *CacheTestSuite) TestCache_Expire() {
	// Add an item to the cache with a TTL of 100 microseconds.
	suite.cache.Add(1, "hello", 100*time.Microsecond)

	// Retrieve the item from the cache.
	item, ok := suite.cache.Get(1)

	// Check if the item was retrieved successfully.
	suite.True(ok, "Item was not retrieved successfully")

	// Check if the retrieved item is not nil.
	suite.NotNil(item, "Item is nil")

	// Check if the retrieved item has the correct value.
	suite.Equal("hello", *item, "Item has incorrect value")

	// Wait for 200 milliseconds to let the item expire.
	time.Sleep(200 * time.Millisecond)

	// Retrieve the item from the cache again.
	item, ok = suite.cache.Get(1)

	// Check if the item was retrieved successfully.
	// Since the item has expired, it should not be retrievable anymore.
	// If the item was not retrieved successfully, the test passes.
	// If the item was retrieved successfully, the test fails.
	suite.False(ok, "Item was not retrieved successfully")

	// Check if the retrieved item is nil.
	// Since the item has expired, it should not be retrievable anymore.
	// If the retrieved item is nil, the test passes.
	// If the retrieved item is not nil, the test fails.
	suite.Nil(item, "Retrieved item is nil") // Check if the retrieved item is nil
}

// TestCache_WithOnEvict tests the OnEvict callback function of the Cache struct.
//
// The test adds multiple items to the cache with different time-to-live (TTL) durations.
// The first item has a TTL of 100 milliseconds, and the second and third items have a TTL of 100 milliseconds.
// The fourth item has a TTL of 700 milliseconds.
//
// The test verifies that the OnEvict callback function is called when items are evicted
// from the cache. It checks that the OnEvict callback function is called only when items
// are evicted from the cache, and not when items are retrieved from the cache.
//
// Parameters:
//
//   - suite: The TestSuite instance that contains the cache under test.
//
// Returns:
//
//	None.
//
//nolint:funlen
func (suite *CacheTestSuite) TestCache_WithOnEvict() {
	// Create a mutex to synchronize access to the onEvictCount and keys variables.
	var mu sync.Mutex

	// Count the number of times the OnEvict callback function is called.
	var onEvictCount int

	// Store the keys of the items that were evicted.
	var keys []int

	// Define the OnEvict callback function.
	// The callback function is called whenever an item is evicted from the cache.
	// It increments the onEvictCount and appends the key of the evicted item to the keys slice.
	onEvict := func(k int, _ string) {
		// Lock the mutex to synchronize access to the onEvictCount and keys variables.
		mu.Lock()
		defer mu.Unlock()

		// Increment the onEvictCount.
		onEvictCount++

		// Append the key of the evicted item to the keys slice.
		keys = append(keys, k)
	}

	// Create a new cache with a maximum size of 10,
	// an evict duration of 100 milliseconds, and an OnEvict callback function.
	suite.cache = cache.NewCache(
		10,
		cache.WithEvictDuration[int, string](100*time.Millisecond),
		cache.WithOnEvict[int, string](onEvict),
	)

	// Add multiple items to the cache.
	suite.cache.Add(1, "hello", 100*time.Millisecond)
	suite.cache.Add(3, "hello", 100*time.Millisecond)
	suite.cache.Add(2, "hello", 100*time.Millisecond)
	suite.cache.Add(15, "world", 700*time.Millisecond)

	// Retrieve the item from the cache.
	item, ok := suite.cache.Get(1)

	// Check if the item was retrieved successfully.
	suite.True(ok, "Item with key 1 was not retrieved successfully")

	// Check if the retrieved item is not nil.
	suite.NotNil(item, "Retrieved item with key 1 is nil")

	// Check if the retrieved item has the correct value.
	suite.Equal("hello", *item, "Retrieved item with key 1 has incorrect value")

	// Wait for 2 seconds to let the items expire.
	time.Sleep(200 * time.Millisecond)

	// Check if the OnEvict callback function was called 3 times.
	suite.Equal(3, onEvictCount, "OnEvict callback function was not called 3 times")

	// Check if the OnEvict callback function was called for the item with key 1.
	suite.True(slices.Contains(keys, 1), "OnEvict callback function was not called for item with key 1")

	// Check if the OnEvict callback function was called for the item with key 2.
	suite.True(slices.Contains(keys, 2), "OnEvict callback function was not called for item with key 2")

	// Check if the OnEvict callback function was called for the item with key 3.
	suite.True(slices.Contains(keys, 3), "OnEvict callback function was not called for item with key 3")
}

// TestCache_OnEvict tests the OnEvict callback function of the Cache struct.
//
// This test verifies that the OnEvict callback function is called for each
// item in the cache when it expires. The test adds two items to the cache
// with different time-to-live (TTL) durations and waits for a short period
// of time to allow the items to expire.
//
// The test initializes an atomic counter to track the number of evicted items.
// It sets the OnEvict callback function of the cache to increment the counter
// by the key of the evicted item.
//
// After adding the items to the cache, the test waits for a short period of time
// to allow the items to expire.
//
// Finally, the test checks if the OnEvict callback function was called 4 times,
// and that the sum of the keys of the evicted items is equal to 4.
func (suite *CacheTestSuite) TestCache_OnEvict() {
	// Initialize an atomic counter to track the number of evicted items.
	var onEvictCount atomic.Int64

	// Set the OnEvict callback function of the cache.
	// The OnEvict callback function increments the counter by the key of the evicted item.
	suite.cache.OnEvict(func(k int, _ string) {
		// Increment the counter by the key of the evicted item.
		onEvictCount.Add(int64(k))
	})

	// Add two items to the cache with different TTLs.
	suite.cache.Add(1, "hello", 100*time.Microsecond)
	suite.cache.Add(3, "hello", 100*time.Microsecond)

	// Wait for a short period of time to allow the items to expire.
	time.Sleep(50 * time.Millisecond)

	// Check if the OnEvict callback function was called 4 times.
	// The sum of the keys of the evicted items should be equal to 4.
	suite.Equal(int64(4), onEvictCount.Load())
}

// TestCache_ProlongLife tests the functionality of prolonging the life of cache items.
//
// This test adds two items to the cache with different time-to-live (TTL) durations.
// The first item has a TTL of 300 milliseconds, and the second item has a TTL of 1 millisecond.
//
// The test retrieves the first item from the cache, checks if it was retrieved successfully,
// and asserts that it has the expected value.
//
// Next, the test retrieves the second item from the cache, checks if it was retrieved successfully,
// and asserts that it has the expected value.
//
// After that, the test waits for the TTL of the first item to expire.
// After the waiting period, it retrieves the first item from the cache again,
// adds a new item with an updated TTL, and checks if the first item was retrieved successfully.
//
// Next, the test retrieves the second item from the cache after expiration time,
// and checks if it was not found.
//
// Finally, the test waits for the TTL of the first item to expire again
// and retrieves the first item from the cache again.
// It checks if the first item was retrieved successfully and has the expected value.
func (suite *CacheTestSuite) TestCache_ProlongLife() {
	// Add two items to the cache with different TTLs.
	suite.cache.Add(1, "hello", 300*time.Millisecond) // Add an item with key 1 and a TTL of 300 milliseconds.
	suite.cache.Add(2, "world", time.Millisecond)     // Add an item with key 2 and a TTL of 1 millisecond.

	// Retrieve the first item from the cache.
	item, ok := suite.cache.Get(1)

	// Assert that the item was retrieved successfully.
	suite.True(ok, "Item with key 1 was not retrieved successfully")
	suite.NotNil(item, "Retrieved item with key 1 is nil")
	suite.Equal("hello", *item, "Retrieved item with key 1 has incorrect value")

	// Retrieve the second item from the cache.
	item, ok = suite.cache.Get(2)

	// Assert that the item was retrieved successfully.
	suite.True(ok, "Item with key 2 was not retrieved successfully")
	suite.NotNil(item, "Retrieved item with key 2 is nil")
	suite.Equal("world", *item, "Retrieved item with key 2 has incorrect value")

	// Wait for the TTL of the first item to expire.
	time.Sleep(200 * time.Millisecond) // Wait for 200 milliseconds to let the first item expire.

	// Retrieve the first item from the cache after expiration time,
	// add a new item with an updated TTL, and check if the first item was retrieved successfully.
	item, ok = suite.cache.Get(1)                     // Retrieve the item with key 1.
	suite.cache.Add(1, "hello", 300*time.Millisecond) // Add a new item with key 1 and an updated TTL of 300 milliseconds.
	suite.True(ok, "Item with key 1 was not retrieved successfully after expiration time")
	suite.NotNil(item, "Retrieved item with key 1 after expiration time is nil")
	suite.Equal("hello", *item, "Retrieved item with key 1 after expiration time has incorrect value")

	// Retrieve the second item from the cache after expiration time and check if it was not found.
	item, ok = suite.cache.Get(2) // Retrieve the item with key 2.
	suite.False(ok, "Item with key 2 was retrieved successfully after expiration time")
	suite.Nil(item, "Retrieved item with key 2 after expiration time is not nil")

	// Wait for the TTL of the first item to expire again.
	time.Sleep(200 * time.Millisecond) // Wait for 200 milliseconds to let the first item expire again.

	// Retrieve the first item from the cache after the second expiration time and check if it was retrieved successfully.
	item, ok = suite.cache.Get(1) // Retrieve the item with key 1.
	suite.True(ok, "Item with key 1 was not retrieved successfully after second expiration time")
	suite.NotNil(item, "Retrieved item with key 1 after second expiration time is nil")
	suite.Equal("hello", *item, "Retrieved item with key 1 after second expiration time has incorrect value")
}

// TestCacheTestSuite runs the CacheTestSuite test suite.
//
// This test suite contains multiple test cases that test the functionality of the Cache struct.
// It is used to verify that the Cache struct operates as expected.
//
// Parameters:
//   - t: The testing.T parameter used to run the test suite.
//     This parameter is used to run the test suite and report the results.
//
// Returns:
//
//	None
func TestCacheTestSuite(t *testing.T) {
	// Mark the test as parallelizable.
	// This allows the test to run concurrently with other tests,
	// which can speed up the test execution time.
	t.Parallel()

	// Run the CacheTestSuite test suite using the suite.Run function.
	// The suite.Run function takes a testing.T parameter and a pointer to the test suite.
	//
	// The function runs the test suite and reports the results.
	// The test suite contains multiple test cases that test the functionality of the Cache struct.
	// It is used to verify that the Cache struct operates as expected.
	//
	// Parameters:
	//   - t: The testing.T parameter used to run the test suite.
	//        This parameter is used to run the test suite and report the results.
	//   - suite: A pointer to the CacheTestSuite instance.
	//        This parameter is used to run the test cases defined in the CacheTestSuite.
	//
	// Returns:
	//   None
	suite.Run(t, new(CacheTestSuite))
}
