package funnelServices

import (
	"net/url"
	"wejh-go/app/apiExpection"
	"wejh-go/app/models"
	"wejh-go/config/api/funnelApi"
)

func GetCurrentBorrow(u *models.User) (interface{}, error) {
	if u.LibPassword == "" {
		return nil, apiExpection.NoThatPasswordOrWrong
	}
	form := url.Values{}
	form.Add("username", u.StudentID)
	form.Add("password", u.LibPassword)

	return FetchHandleOfPost(form, funnelApi.LibraryCurrent)
}

func GetHistoryBorrow(u *models.User) (interface{}, error) {
	if u.LibPassword == "" {
		return nil, apiExpection.NoThatPasswordOrWrong
	}
	form := url.Values{}
	form.Add("username", u.StudentID)
	form.Add("password", u.LibPassword)

	return FetchHandleOfPost(form, funnelApi.LibraryHistory)
}
