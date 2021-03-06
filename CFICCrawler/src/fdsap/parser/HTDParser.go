package parser

import "httpcontroller"

func (doc *HTMLDoc) HTD_Request(url string, file string, proxy *httpcontroller.Proxy) error {
	request := httpcontroller.Request{
		//Proxy:     proxy, //&httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"},
		Url:       url,
		File:      file,
		OverWrite: false,
	}

	if _, err := request.Get(); err != nil {
		return err
	}

	return nil
}
