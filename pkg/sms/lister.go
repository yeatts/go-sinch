package sms

import (
	"fmt"
	"net/http"

	"github.com/thezmc/go-sinch/pkg/interfaces"
	"github.com/thezmc/go-sinch/pkg/models"
)

type batchLister struct {
	client   *client
	request  *models.ListRequest
	response *models.ListResponse
	err      error
}

func (c *client) NewBatchLister() interfaces.SMSBatchLister {
	return &batchLister{client: c}
}

func (l *batchLister) FromPage(page int) interfaces.SMSBatchLister {
	l.request.Page = page
	return l
}

func (l *batchLister) WithPageSize(pageSize int) interfaces.SMSBatchLister {
	l.request.PageSize = pageSize
	return l
}

func (l *batchLister) From(from ...string) interfaces.SMSBatchLister {
	l.request.From = append(l.request.From, from...)
	return l
}
func (l *batchLister) WithStartDate(startDate string) interfaces.SMSBatchLister {
	l.request.StartDate = startDate
	return l
}

func (l *batchLister) WithEndDate(endDate string) interfaces.SMSBatchLister {
	l.request.EndDate = endDate
	return l
}
func (l *batchLister) WithClientReference(clientReference string) interfaces.SMSBatchLister {
	l.request.ClientReference = clientReference
	return l
}

func (l *batchLister) validate() error {
	if l.request.PageSize < 1 && l.request.PageSize > 100 {
		return fmt.Errorf("page size must be between 1 and 100")
	}
	return nil
}

func (l *batchLister) toQueryString() string {
	queryString := ""
	if l.request.Page > 0 {
		queryString += fmt.Sprintf("&page=%d", l.request.Page)
	}
	if l.request.PageSize > 0 {
		queryString += fmt.Sprintf("&page_size=%d", l.request.PageSize)
	}
	if l.request.From != nil {
		froms := ""
		for _, from := range l.request.From {
			froms += from + ","
		}
		queryString += fmt.Sprintf("&from=%s", l.request.From)
	}
	if l.request.StartDate != "" {
		queryString += fmt.Sprintf("&start_date=%s", l.request.StartDate)
	}
	if l.request.EndDate != "" {
		queryString += fmt.Sprintf("&end_date=%s", l.request.EndDate)
	}
	if l.request.ClientReference != "" {
		queryString += fmt.Sprintf("&client_reference=%s", l.request.ClientReference)
	}
	if len(queryString) > 0 && queryString[0] == '&' {
		queryString = queryString[1:]
	}
	return queryString
}

func (l *batchLister) toRequest() *http.Request {
	req := new(http.Request)
	req.Method = "GET"
	req.URL.Path = fmt.Sprintf("%s/%s/%s", l.client.baseURL, l.client.planID, ResouceName)
	req.URL.RawQuery = l.toQueryString()
	return req
}

func (l *batchLister) List() interfaces.SMSBatchLister {

	return l
}

func (l *batchLister) Error() error {
	return l.err
}
func (l *batchLister) Response() *models.ListResponse {
	return l.response
}
