package xmetricstest

import (
	"testing"

	"github.com/Comcast/webpa-common/xmetrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testNewProviderDefault(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewProvider(nil))
}

func testNewProviderGoodConfiguration(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewProvider(
		&xmetrics.Options{
			Metrics: []xmetrics.Metric{
				{Name: "Injected", Type: "counter"},
			},
		},
		func() []xmetrics.Metric {
			return []xmetrics.Metric{
				{Name: "FromModule", Type: "gauge"},
			}
		},
	))
}

func testNewProviderBadConfiguration(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		NewProvider(nil, func() []xmetrics.Metric {
			return []xmetrics.Metric{
				{Name: "duplicate", Type: "counter"},
				{Name: "duplicate", Type: "counter"},
			}
		})
	})
}

func testNewProviderUnsupportedType(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		NewProvider(nil, func() []xmetrics.Metric {
			return []xmetrics.Metric{
				{Name: "unsupported", Type: "asdfasdfasdfasdf"},
			}
		})
	})
}

func TestNewProvider(t *testing.T) {
	t.Run("Default", testNewProviderDefault)
	t.Run("GoodConfiguration", testNewProviderGoodConfiguration)
	t.Run("BadConfiguration", testNewProviderBadConfiguration)
	t.Run("UnsupportedType", testNewProviderUnsupportedType)
}

func exampleProvider() Provider {
	return NewProvider(nil, func() []xmetrics.Metric {
		return []xmetrics.Metric{
			{Name: "counter", Type: "counter"},
			{Name: "gauge", Type: "gauge"},
			{Name: "histogram", Type: "histogram"},
		}
	})
}

func testProviderNewCounter(t *testing.T) {
	var (
		assert   = assert.New(t)
		provider = exampleProvider()
	)

	assert.Panics(func() {
		provider.NewCounter("gauge")
	})

	preconfigured := provider.NewCounter("counter")
	assert.NotNil(preconfigured)
	assert.Implements((*xmetrics.Valuer)(nil), preconfigured)
	assert.True(preconfigured == provider.NewCounter("counter"))

	adhoc := provider.NewCounter("adhoc")
	assert.NotNil(adhoc)
	assert.Implements((*xmetrics.Valuer)(nil), adhoc)
	assert.True(adhoc == provider.NewCounter("adhoc"))
	assert.True(preconfigured != adhoc)
}

func testProviderAssertValue(t *testing.T) {
	t.Run("DoesNotExist", func(t *testing.T) {
		var (
			assert   = assert.New(t)
			testingT = new(mockTestingT)
			provider = exampleProvider()
		)

		testingT.On("Errorf", mock.MatchedBy(func(string) bool { return true }), mock.MatchedBy(func([]interface{}) bool { return true })).Once()

		v, ok := provider.AssertValue(testingT, "doesnotexist", 1.0)
		assert.Nil(v)
		assert.False(ok)

		testingT.AssertExpectations(t)
	})

	t.Run("NonValuer", func(t *testing.T) {
		var (
			assert   = assert.New(t)
			testingT = new(mockTestingT)
			provider = exampleProvider()
		)

		testingT.On("Errorf", mock.MatchedBy(func(string) bool { return true }), mock.MatchedBy(func([]interface{}) bool { return true })).Once()

		v, ok := provider.AssertValue(testingT, "histogram", 1.0)
		assert.Nil(v)
		assert.False(ok)

		testingT.AssertExpectations(t)
	})

	t.Run("Preregistered", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			var (
				assert   = assert.New(t)
				testingT = new(mockTestingT)
				provider = exampleProvider()
			)

			testingT.On("Errorf", mock.MatchedBy(func(string) bool { return true }), mock.MatchedBy(func([]interface{}) bool { return true })).Once()

			v, ok := provider.AssertValue(testingT, "counter", 0.0)
			assert.NotNil(v)
			assert.Equal(0.0, v.Value())
			assert.True(ok)
			testingT.AssertNumberOfCalls(t, "Errorf", 0)

			provider.NewCounter("counter").Add(1.0)
			v, ok = provider.AssertValue(testingT, "counter", 1.0)
			assert.NotNil(v)
			assert.Equal(1.0, v.Value())
			assert.True(ok)
			testingT.AssertNumberOfCalls(t, "Errorf", 0)

			v, ok = provider.AssertValue(testingT, "counter", 0.0)
			assert.NotNil(v)
			assert.Equal(1.0, v.Value())
			assert.False(ok)

			testingT.AssertExpectations(t)
		})

		t.Run("Gauge", func(t *testing.T) {
			var (
				assert   = assert.New(t)
				testingT = new(mockTestingT)
				provider = exampleProvider()
			)

			testingT.On("Errorf", mock.MatchedBy(func(string) bool { return true }), mock.MatchedBy(func([]interface{}) bool { return true })).Once()

			v, ok := provider.AssertValue(testingT, "gauge", 0.0)
			assert.NotNil(v)
			assert.Equal(0.0, v.Value())
			assert.True(ok)
			testingT.AssertNumberOfCalls(t, "Errorf", 0)

			provider.NewGauge("gauge").Add(1.0)
			v, ok = provider.AssertValue(testingT, "gauge", 1.0)
			assert.NotNil(v)
			assert.Equal(1.0, v.Value())
			assert.True(ok)
			testingT.AssertNumberOfCalls(t, "Errorf", 0)

			v, ok = provider.AssertValue(testingT, "gauge", 0.0)
			assert.NotNil(v)
			assert.Equal(1.0, v.Value())
			assert.False(ok)

			testingT.AssertExpectations(t)
		})
	})

	t.Run("AdHoc", func(t *testing.T) {
		var (
			assert   = assert.New(t)
			testingT = new(mockTestingT)
			provider = exampleProvider()
		)

		testingT.On("Errorf", mock.MatchedBy(func(string) bool { return true }), mock.MatchedBy(func([]interface{}) bool { return true })).Once()

		c := provider.NewCounter("adhoc")

		v, ok := provider.AssertValue(testingT, "adhoc", 0.0)
		assert.NotNil(v)
		assert.Equal(0.0, v.Value())
		assert.True(ok)
		testingT.AssertNumberOfCalls(t, "Errorf", 0)

		c.Add(1.0)
		v, ok = provider.AssertValue(testingT, "adhoc", 1.0)
		assert.NotNil(v)
		assert.Equal(1.0, v.Value())
		assert.True(ok)
		testingT.AssertNumberOfCalls(t, "Errorf", 0)

		v, ok = provider.AssertValue(testingT, "adhoc", 0.0)
		assert.NotNil(v)
		assert.Equal(1.0, v.Value())
		assert.False(ok)

		testingT.AssertExpectations(t)
	})
}

func testProviderExpectValue(t *testing.T) {
}

func TestProvider(t *testing.T) {
	t.Run("AssertValue", testProviderAssertValue)
	t.Run("ExpectValue", testProviderExpectValue)
	t.Run("NewCounter", testProviderNewCounter)
}
