package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/hollgett/shortURL.git/internal/logger"
	"go.uber.org/zap"
)

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zw.Write(p)
}
func (cw *compressWriter) Close() error {
	err := cw.zw.Close()
	if err != nil {
		logger.LogInfo("close reader compress", zap.Error(err))
		return err
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zw *gzip.Reader
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zw.Read(p)
}

func (cr *compressReader) Close() error {
	err := cr.r.Close()
	if err != nil {
		logger.LogInfo("close read body", zap.Error(err))
		return err
	}
	err = cr.zw.Close()
	if err != nil {
		logger.LogInfo("close reader compress", zap.Error(err))
		return err
	}
	return nil
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aliasW := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gzipWriter := &compressWriter{
				ResponseWriter: w,
				zw:             gzip.NewWriter(w),
			}
			aliasW = gzipWriter
			defer gzipWriter.Close()
			w.Header().Set("Content-Encoding", "gzip")
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cr := &compressReader{
				r:  r.Body,
				zw: gzReader,
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(aliasW, r)
	})
}
