package database

import (
	"testing"

	"github.com/donnigundala/dg-core/foundation"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	app := foundation.New(".")
	config := Config{
		Driver:   "sqlite",
		FilePath: ":memory:",
	}

	provider := &DatabaseServiceProvider{Config: config}
	err := provider.Register(app)
	assert.NoError(t, err)

	// Test Resolve
	db, err := Resolve(app)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB())
}

func TestResolve_Error(t *testing.T) {
	app := foundation.New(".")

	// Test Resolve without registration
	db, err := Resolve(app)
	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to resolve database")
}

func TestMustResolve(t *testing.T) {
	app := foundation.New(".")
	config := Config{
		Driver:   "sqlite",
		FilePath: ":memory:",
	}

	provider := &DatabaseServiceProvider{Config: config}
	err := provider.Register(app)
	assert.NoError(t, err)

	// Test MustResolve
	db := MustResolve(app)
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB())
}

func TestMustResolve_Panic(t *testing.T) {
	app := foundation.New(".")

	// Test MustResolve panics without registration
	assert.Panics(t, func() {
		MustResolve(app)
	})
}

func TestResolveConnection(t *testing.T) {
	app := foundation.New(".")
	config := Config{
		Driver:   "sqlite",
		FilePath: ":memory:",
	}

	// Add named connection
	config.Connections = map[string]ConnectionConfig{
		"analytics": {
			Driver:   "sqlite",
			FilePath: ":memory:",
		},
	}

	provider := &DatabaseServiceProvider{Config: config}
	err := provider.Register(app)
	assert.NoError(t, err)

	// Test ResolveConnection
	analyticsDB, err := ResolveConnection(app, "analytics")
	assert.NoError(t, err)
	assert.NotNil(t, analyticsDB)
}

func TestResolveConnection_Error(t *testing.T) {
	app := foundation.New(".")

	// Test ResolveConnection without registration
	db, err := ResolveConnection(app, "analytics")
	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to resolve connection")
}

func TestMustResolveConnection(t *testing.T) {
	app := foundation.New(".")
	config := Config{
		Driver:   "sqlite",
		FilePath: ":memory:",
	}

	config.Connections = map[string]ConnectionConfig{
		"analytics": {
			Driver:   "sqlite",
			FilePath: ":memory:",
		},
	}

	provider := &DatabaseServiceProvider{Config: config}
	err := provider.Register(app)
	assert.NoError(t, err)

	// Test MustResolveConnection
	analyticsDB := MustResolveConnection(app, "analytics")
	assert.NotNil(t, analyticsDB)
}

func TestMustResolveConnection_Panic(t *testing.T) {
	app := foundation.New(".")

	// Test MustResolveConnection panics without registration
	assert.Panics(t, func() {
		MustResolveConnection(app, "analytics")
	})
}

func TestInjectable(t *testing.T) {
	app := foundation.New(".")
	config := Config{
		Driver:   "sqlite",
		FilePath: ":memory:",
	}

	config.Connections = map[string]ConnectionConfig{
		"analytics": {
			Driver:   "sqlite",
			FilePath: ":memory:",
		},
	}

	provider := &DatabaseServiceProvider{Config: config}
	err := provider.Register(app)
	assert.NoError(t, err)

	// Test Injectable
	inject := NewInjectable(app)
	assert.NotNil(t, inject)

	// Test DB()
	db := inject.DB()
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB())

	// Test Connection()
	analyticsDB := inject.Connection("analytics")
	assert.NotNil(t, analyticsDB)

	// Test TryConnection() - existing
	tryDB := inject.TryConnection("analytics")
	assert.NotNil(t, tryDB)

	// Test TryConnection() - non-existing
	nilDB := inject.TryConnection("nonexistent")
	assert.Nil(t, nilDB)
}

func TestInjectable_Panic(t *testing.T) {
	app := foundation.New(".")

	inject := NewInjectable(app)

	// Test DB() panics without registration
	assert.Panics(t, func() {
		inject.DB()
	})

	// Test Connection() panics without registration
	assert.Panics(t, func() {
		inject.Connection("analytics")
	})
}
