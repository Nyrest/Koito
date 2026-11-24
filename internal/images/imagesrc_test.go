package images

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gabehf/koito/queue"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAlbumImageLookupWithoutMBID(t *testing.T) {
	requests := make(chan url.Values, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.URL.Query()
		_, _ = io.WriteString(w, `{}`)
	}))
	t.Cleanup(server.Close)

	requestQueue := queue.NewRequestQueue(5, 5)
	t.Cleanup(requestQueue.Shutdown)
	lastFM := &LastFMClient{baseUrl: server.URL, requestQueue: requestQueue}
	subsonic := &SubsonicClient{url: server.URL, requestQueue: requestQueue}

	_, err := lastFM.GetAlbumImage(t.Context(), uuid.Nil, "Test Artist", "Test Album")
	require.ErrorContains(t, err, "no suitable image found")
	_, err = subsonic.GetAlbumImage(t.Context(), uuid.Nil, "Test Artist", "Test Album")
	require.ErrorContains(t, err, "failed to get album art")

	require.Len(t, requests, 2)
	lastFMQuery := <-requests
	require.Equal(t, "Test Artist", lastFMQuery.Get("artist"))
	require.Equal(t, "Test Album", lastFMQuery.Get("album"))
	require.NotContains(t, lastFMQuery, "mbid")
	subsonicQuery := <-requests
	require.Equal(t, "Test Album", subsonicQuery.Get("query"))
}
