package client_http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Instance *http.Client
}

type Response struct {
	Body []byte
	Status string
	StatusCode int
}

type HeaderParameters struct {
	Key string
	Value string
}

func NewHttpClient(skipTLS bool) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 1000
	transport.MaxConnsPerHost = 1000
	transport.MaxIdleConns = 1000

	if skipTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{Transport: transport, Timeout: 120 * time.Second}
	return &Client{Instance: httpClient}

}

// GetResponseWithCredentials - Get response from url with credentials
func (c *Client) GetResponseWithCredentials(url, username, password string) (*Response, error) {
	// Get request for url
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't get request error [%v]", err)
	}

	// set Credentials
	request.SetBasicAuth(username, password)

	// Do request
	response, err := c.Instance.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't do request error [%v]", err)
	}

	// defer closing body
	defer Defer(func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				fmt.Printf("can't close body error [%v]", err)
			}
		}
	})

	// Read body response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body error [%v]", err)
	}

	// Create Result
	return &Response {
		Body: body,
		Status: response.Status,
		StatusCode: response.StatusCode,
	}, nil


}



// GetResponseWithPayloadAndAuth - Get response sending payload, authentication header
func (c *Client) GetResponseWithPayloadAndAuth(url, username, password string, payload []byte) (*Response, error){
	// Get request for url and payload
	request, err := http.NewRequest("GET", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error on build request [%s] - [%v]", url, err)
	}

	// set Authentication headers
	request.SetBasicAuth(username, password)

	// Do request
	response,err := c.Instance.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error on make request [%v] ", err)
	}

	// defer body closing
	defer response.Body.Close()
	//defer Defer(func() {
	//	if response.Body != nil {
	//		err := response.Body.Close()
	//		if err != nil {
	//			fmt.Printf("error closing response.body [%v]", err)
	//		}
	//	}
	//})

	// reading body result
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	// returning response
	return &Response{Body: body, Status: response.Status, StatusCode: response.StatusCode}, nil


}

// GetResponseWithPayloadAuthAndHeader - Get response sending payload, authentication header and headers
func (c *Client) GetResponseWithPayloadAuthAndHeader(url, username, password string, payload []byte, headers []HeaderParameters) (*Response, error){
	// Get request for url and payload
	request, err := http.NewRequest("GET", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error on build request [%s] - [%v]", url, err)
	}

	// set Authentication headers
	request.SetBasicAuth(username, password)

	// set additional headers
	for _, h := range headers {
		request.Header.Set(h.Key, h.Value)
	}

	// Do request
	response,err := c.Instance.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error on make request [%v] ", err)
	}

	// defer body closing
	defer Defer(func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				fmt.Printf("error closing response.body [%v]", err)
			}
		}
	})

	// reading body result
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	// returning response
	return &Response{Body: body, Status: response.Status, StatusCode: response.StatusCode}, nil


}

// GetResponseWithPayloadAndHeaders - Get response using url, payload and custom headers
func (c *Client) GetResponseWithPayloadAndHeaders(url string, payload []byte, headers []HeaderParameters) (*Response, error) {
	// create request
	request, err := http.NewRequest("GET", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request for url [%s] - [%v]", url, err)
	}

	// set additional headers
	for _, h := range headers {
		request.Header.Set(h.Key, h.Value)
	}


	// do request
	response, err := c.Instance.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error doing request [%v]", err)
	}

	// closing body
	defer Defer(func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				fmt.Printf("error closing response body [%v]", err)
			}
		}
	})

	// reading data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return nil, fmt.Errorf("error reading response body [%v]", err)
	}

	// return response
	return &Response{
		Body:       body,
		Status:     response.Status,
		StatusCode: response.StatusCode,
	}, nil

}

// GetResponse - execute a simple request on url
func (c *Client) GetResponse(url string) (*Response, error) {
	// creating request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for url [%s] =  [%v]", url, err)
	}

	// executing request
	response, err:= c.Instance.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing request for url [%s] =  [%v]", url, err)
	}

	// closing body response
	defer Defer(func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				fmt.Printf("error closing response body [%v]", err)
			}
		}
	})

	// reading body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return nil, fmt.Errorf("error reading response body [%v]", err)
	}


	// return response
	return &Response{
		Body:       body,
		Status:     response.Status,
		StatusCode: response.StatusCode,
	}, nil

}

func Defer(f func()) {
	defer f()
}

