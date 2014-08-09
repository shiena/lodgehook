package lodgehook

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func dumpPostForm(req *http.Request) (url.Values, error) {
	save := req.Body
	var err error

	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		req.Body = save
	}()

	err = req.ParseForm()
	if err != nil {
		return nil, err
	}

	return req.PostForm, nil
}
