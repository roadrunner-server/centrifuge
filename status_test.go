package centrifuge

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginStatusNilPool(t *testing.T) {
	p := &Plugin{}

	st, err := p.Status()
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, st.Code)
}

func TestPluginReadyNilPool(t *testing.T) {
	p := &Plugin{}

	st, err := p.Ready()
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, st.Code)
}
