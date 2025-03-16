package yxyServices

import (
	"net/url"
	"wejh-go/app/apiException"
	"wejh-go/config/api/yxyApi"

	"github.com/mitchellh/mapstructure"
)

type authResp struct {
	Token string `json:"token" mapstructure:"token"`
}

type busResp struct {
	List []struct {
		ID       string `json:"id" mapstructure:"id"`
		Name     string `json:"name" mapstructure:"name"`
		Seats    int    `json:"seats" mapstructure:"seats"`
		Price    int    `json:"price" mapstructure:"price"`
		Stitions []struct {
			ID   string `json:"id" mapstructure:"id"`
			Name string `json:"station_name" mapstructure:"station_name"`
			Seq  int    `json:"station_seq" mapstructure:"station_seq"`
		} `json:"stations" mapstructure:"stations"`
		BusTime []struct {
			ID            string `json:"id" mapstructure:"id"`
			DepartureTime string `json:"departure_time" mapstructure:"departure_time"`
			RemainSeats   int    `json:"remain_seats" mapstructure:"remain_seats"`
			OrderedSeats  int    `json:"ordered_seats" mapstructure:"ordered_seats"`
		}
	} `json:"list" mapstructure:"list"`
}

type busQrCodeResp struct {
	QrCode string `json:"qrcode" mapstructure:"qrcode"`
}

type BusMessagw struct {
	List []struct {
		ID      string `json:"id"`
		MsgType string `json:"msg_type"`
		MsgID   string `json:"msg_id"`
		Title   string `json:"title"`
		IsRead  int    `json:"is_read"`
		Content string `json:"content"`
		HTML    string `json:"html"`
		Img     string `json:"img"`
		Author  string `json:"author"`
	} `json:"list" mapstructure:"list"`
}

func BusAuth(uid string) (*string, error) {
	params := url.Values{}
	Url, err := url.Parse(string(yxyApi.BusAuth))
	if err != nil {
		return nil, err
	}
	params.Set("uid", uid)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := FetchHandleOfGet(yxyApi.YxyApi(urlPath))
	if err != nil {
		return nil, err
	}
	var data authResp
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data.Token, nil
}

func BusInfo(page, pageSize, search string) (*busResp, error) {
	params := url.Values{}
	Url, err := url.Parse(string(yxyApi.BusInfo))
	if err != nil {
		return nil, err
	}
	params.Set("page", page)
	params.Set("page_size", pageSize)
	params.Set("search", search)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := FetchHandleOfGet(yxyApi.YxyApi(urlPath))
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, apiException.ServerError
	}

	var data busResp
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func BusRecords(token, page, pageSize, status string) (*busResp, error) {
	params := url.Values{}
	Url, err := url.Parse(string(yxyApi.BusRecords))
	if err != nil {
		return nil, err
	}
	params.Set("token", token)
	params.Set("page", page)
	params.Set("page_size", pageSize)
	params.Set("status", status)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := FetchHandleOfGet(yxyApi.YxyApi(urlPath))
	if err != nil {
		return nil, err
	}

	if resp.Code == 110001 {
		return nil, apiException.YxySessionExpired
	} else if resp.Code != 0 {
		return nil, apiException.ServerError
	}

	var data busResp
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func BusQrcode(token string) (*busQrCodeResp, error) {
	params := url.Values{}
	Url, err := url.Parse(string(yxyApi.BusQrCode))
	if err != nil {
		return nil, err
	}
	params.Set("token", token)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := FetchHandleOfGet(yxyApi.YxyApi(urlPath))
	if err != nil {
		return nil, err
	}

	if resp.Code == 110001 {
		return nil, apiException.YxySessionExpired
	} else if resp.Code != 0 {
		return nil, apiException.ServerError
	}

	var data busQrCodeResp
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func BusMessage(token, page, pageSize string) (*BusMessagw, error) {
	params := url.Values{}
	Url, err := url.Parse(string(yxyApi.BusMessage))
	if err != nil {
		return nil, err
	}
	params.Set("token", token)
	params.Set("page", page)
	params.Set("page_size", pageSize)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := FetchHandleOfGet(yxyApi.YxyApi(urlPath))
	if err != nil {
		return nil, err
	}

	if resp.Code == 110001 {
		return nil, apiException.YxySessionExpired
	} else if resp.Code != 0 {
		return nil, apiException.ServerError
	}

	var data BusMessagw
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
