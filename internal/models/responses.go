package models

type ResponseBody struct {
	Link string
	OriginalUrl string
}


type RecordResponse struct {
	Shorturl string
	Longurl string
	CreatedAt string
}
type RecordsResponse struct {
	IP string
	Records []*RecordResponse
	Count int
}