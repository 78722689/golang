package htmlparser

import "httpcontroller"

func (doc *HTMLDoc)HTD_Request(url string, file string) error{
	request := httpcontroller.Request{
		//Proxy:&httpcontroller.Proxy{"HTTP", "10.144.1.10", "8080"},
		Url : url,
		File : file,
	}

	if _, err := request.Get(); err != nil {
		return err
	}

	return nil
}