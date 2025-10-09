package happening

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveConnectionURL(t *testing.T) {
	databaseName := deriveDatabaseName("postgresql://postgres:secret@localhost:6666/dbname?sslmode=disable")
	assert.Equal(t, "dbname", databaseName, "database name not parsed correctly")
	databaseName = deriveDatabaseName("postgresql://postgres:secret@localhost:6666/%2Fvar%2Flib%2Fpostgresql/dbname?sslmode=disable")
	assert.Equal(t, "dbname", databaseName, "database name not parsed correctly")
}
