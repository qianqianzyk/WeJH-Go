package yxyServices

import (
	"net/url"
	"wejh-go/app/apiException"
	"wejh-go/config/api/yxyApi"

	"github.com/mitchellh/mapstructure"
)

type BusInfoResp struct {
	List []struct {
		Name     string   `json:"name" mapstructure:"name"`
		Seats    int      `json:"seats" mapstructure:"seats"`
		Price    int      `json:"price" mapstructure:"price"`
		Stations []string `json:"stations" mapstructure:"stations"`
		BusTime  []struct {
			DepartureTime string `json:"departure_time" mapstructure:"departure_time"`
			RemainSeats   int    `json:"remain_seats" mapstructure:"remain_seats"`
			OrderedSeats  int    `json:"ordered_seats" mapstructure:"ordered_seats"`
		} `json:"bus_time" mapstructure:"bus_time"`
	} `json:"list" mapstructure:"list"`
}

func GetBusInfo(page, pageSize, search string) (*BusInfoResp, error) {
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

	var data BusInfoResp
	err = mapstructure.Decode(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
