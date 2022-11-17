package httpapi

import (
	"html"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	otlog "github.com/opentracing/opentracing-go/log"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/handlerutil"
	"github.com/sourcegraph/sourcegraph/internal/authz"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/featureflag"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	streamhttp "github.com/sourcegraph/sourcegraph/internal/search/streaming/http"
	"github.com/sourcegraph/sourcegraph/internal/trace"
)

// serveStreamBlame returns a HTTP handler that streams back the results of running
// git blame with the --incremental flag. It will stream back to the client the most
// recent hunks first and will gradually reach the oldests, or not if we timeout
// before that.
//
//	http://localhost:3080/github.com/gorilla/mux/-/stream-blame/mux.go
func serveStreamBlame(logger log.Logger, db database.DB, gitserverClient gitserver.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flags := featureflag.FromContext(r.Context())
		if !flags.GetBoolOr("enable-streaming-git-blame", false) {
			w.WriteHeader(404)
			return
		}
		tr, ctx := trace.New(r.Context(), "blame.Stream", "")
		defer tr.Finish()
		r = r.WithContext(ctx)

		if _, ok := mux.Vars(r)["Repo"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
		}

		repo, commitID, err := handlerutil.GetRepoAndRev(r.Context(), logger, db, mux.Vars(r))

		requestedPath := mux.Vars(r)["Path"]

		streamWriter, err := streamhttp.NewWriter(w)
		if err != nil {
			tr.SetError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Log events to trace
		streamWriter.StatHook = func(stat streamhttp.WriterStat) {
			fields := []otlog.Field{
				otlog.String("streamhttp.Event", stat.Event),
				otlog.Int("bytes", stat.Bytes),
				otlog.Int64("duration_ms", stat.Duration.Milliseconds()),
			}
			if stat.Error != nil {
				fields = append(fields, otlog.Error(stat.Error))
			}
			tr.LogFields(fields...)
		}

		if strings.HasPrefix(requestedPath, "/") {
			requestedPath = strings.TrimLeft(requestedPath, "/")
		}

		hunkReader, err := gitserverClient.StreamBlameFile(r.Context(), authz.DefaultSubRepoPermsChecker, repo.Name, requestedPath, &gitserver.BlameOptions{
			NewestCommit: commitID,
		})
		if err != nil {
			tr.SetError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			hunk, done, err := hunkReader.Read()
			if err != nil {
				tr.SetError(err)
				http.Error(w, html.EscapeString(err.Error()), http.StatusInternalServerError)
				return
			}
			if done {
				streamWriter.Event("done", map[string]any{})
				return
			}
			if err := streamWriter.Event("hunk", hunk); err != nil {
				tr.SetError(err)
				http.Error(w, html.EscapeString(err.Error()), http.StatusInternalServerError)
				return
			}
		}
	}
}
