package internal

//PrinterRequest struct to handle incoming Request as Json
type PrinterRequest struct {
	JobTitle    string `json:"jobTitle"`
	JobLocation string `json:"jobLocation"`
}

//Crawler interface
type Crawler interface {
	Scrape(PrinterRequest)
	ProcessData(key interface{}, value interface{})
}
