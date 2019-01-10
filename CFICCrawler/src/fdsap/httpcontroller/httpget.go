package httpcontroller

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "golang.org/x/net/html"
    "github.com/spf13/viper"
    "io"
    "path"

    "fdsap/utility"
)

var logger = utility.GetLogger()

type Proxy struct {
    Protocol string
    Host string
    Port string
}

type Request struct {
    //Proxy *Proxy // request by proxy

    Url string  // Url to request
    File string // Optional field, file to write the request result
    Root *html.Node // Nodes to buffer the request result
    OverWrite bool // Overwrite file if already exist.
}

func (r *Request)isFileExist(name string) bool{
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// Request to URL, and default return the response by html.Node;
// If file passed, the request result will be write to file only.
func (r *Request)Get() (*html.Node, error){
    if len(r.Url) > 0 {
        url_i := url.URL{}
        var client *http.Client

        // If proxy gave, use proxy to do the request
        if viper.GetStringMap("proxy")["host"] != nil && viper.GetStringMap("proxy")["host"] != "" {
        //if r.Proxy != nil {
            //url_proxy, _ := url_i.Parse(r.Proxy.Protocol + "://" + r.Proxy.Host + ":" + r.Proxy.Port)
            url_proxy, _ := url_i.Parse(viper.GetStringMapString("proxy")["host"])
            transport := &http.Transport{Proxy: http.ProxyURL(url_proxy)}
            client = &http.Client{Transport: transport}
        } else {
            client = &http.Client{}
        }

        // If file name passed, write the result to file.
        if len(r.File) > 0 {
            if r.isFileExist(r.File) {
                if r.OverWrite {
                    os.Remove(r.File)
                } else {
                    logger.Warningf("File %s is already exist, skip the request. Set OverWrite to true to overwite the request result.",
                                            r.File)
                    return nil, nil
                }
            }
            os.MkdirAll(path.Dir(r.File), 0777)

            logger.Infof("Requesting %s", r.Url)

            resp, err := client.Get(r.Url)
            if err != nil {
                fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

                return nil, err
            }
            defer resp.Body.Close()

            file, err:= os.OpenFile(r.File, os.O_RDWR | os.O_CREATE, 0777)
            if err != nil {
                fmt.Fprintf(os.Stderr, "WARN: Open file %s failed, %s\n", r.File, err)
                return nil, err
            }
            defer file.Close()

            io.Copy(file, resp.Body)
        } else {
            resp, err := client.Get(r.Url)
            if err != nil {
                fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

                return nil, err
            }
            defer resp.Body.Close()

            root, err := html.Parse(resp.Body)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: parse response from url\n")

                return nil, err
            }

            return root, nil
        }
    }

    return nil, nil
}