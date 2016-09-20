package wellez

type Data struct {
	WellInfo           []WellInfo           `xml:"WellInfo>row"`
	CompletionCost     []CompletionCost     `xml:"CompletionCost>row"`
	CompletionCostItem []CompletionCostItem `xml:"CompletionCostItem>row"`
}

type WellInfo struct {
	WellID             int64   `json:"well_id" xml:"well_id,attr"`
	WellName           string  `json:"WellName" xml:"WellName,attr"`
	ClientID           int64   `json:"client_id" xml:"client_id,attr"`
	AssetID            int64   `json:"asset_id" xml:"asset_id,attr"`
	State              string  `json:"state" xml:"state,attr"`
	LongitudeLocation  string  `json:"LongitudeLocation" xml:"LongitudeLocation,attr"`
	LatitudeLocation   string  `json:"LatitudeLocation" xml:"LatitudeLocation,attr"`
	Lease              string  `json:"Lease" xml:"Lease,attr"`
	GUID               string  `json:"GUID" xml:"GUID,attr"`
	FieldName          string  `json:"FieldName" xml:"FieldName,attr"`
	District           string  `json:"District" xml:"District,attr"`
	County             string  `json:"County" xml:"County,attr"`
	Country            string  `json:"Country" xml:"Country,attr"`
	CommonWellName     string  `json:"CommonWellName" xml:"CommonWellName,attr"`
	Comment            string  `json:"Comment" xml:"Comment,attr"`
	AssetName          string  `json:"AssetName" xml:"AssetName,attr"`
	Area               string  `json:"Area" xml:"Area,attr"`
	APINo              string  `json:"APINo" xml:"APINo,attr"`
	Unit               string  `json:"Unit" xml:"Unit,attr"`
	StatePermitNo      string  `json:"StatePermitNo" xml:"StatePermiitNo,attr"`
	Section            string  `json:"Section" xml:"Section,attr"`
	Township           string  `json:"Township" xml:"Township,attr"`
	Range              string  `json:"Range" xml:"Range,attr"`
	NetRevenueInterest float64 `json:"NetRevenueInterest" xml:"NetRevenueInterest,attr"`
}

type CompletionCost struct {
	WellID       int64  `json:"well_id" xml:"well_id,attr"`
	ReportDate   string `json:"report_date" xml:"report_date,attr"`
	JobNumber    int64  `json:"job_number" xml:"job_number,attr"`
	ReportTypeID int64  `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	Comment      string `json:"Comment" xml:"Comment,attr"`
	LocationID   int64  `json:"location_id" xml:"location_id,attr"`
	RowID        int64  `json:"RowID" xml:"RowID,attr"`
}

type CompletionCostItem struct {
	WellID       int64   `json:"well_id" xml:"well_id,attr"`
	ReportDate   string  `json:"report_date" xml:"report_date,attr"`
	LineNumber   int64   `json:"line_number" xml:"line_number,attr"`
	JobNumber    int64   `json:"job_number" xml:"job_number,attr"`
	Vendor       string  `json:"Vendor" xml:"Vendor,attr"`
	ReportTypeID int64   `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	Remarks      string  `json:"Remarks" xml:"Remarks,attr"`
	ItemCode     string  `json:"ItemCode" xml:"ItemCode,attr"`
	ExpenseDesc  string  `json:"ExpenseDescription" xml:"ExpenseDescription,attr"`
	Cost         float64 `json:"Cost" xml:"Cost,attr"`
	LocationID   int64   `json:"location_id" xml:"location_id,attr"`
	RowID        int64   `json:"RowID" xml:"RowID,attr"`
}
