package wellez

type Data struct {
	CompletionCost     []CompletionCost     `xml:"CompletionCost>row"`
	CompletionCostItem []CompletionCostItem `xml:"CompletionCostItem>row"`
	CostAllocation     []CostAllocation     `xml:"CostAllocation>row"`
	CostAllocationItem []CostAllocationItem `xml:"CostAllocationItem>row"`
	DailyOps           []DailyOps           `xml:"DailyOps>row"`
	DrillingCost       []DrillingCost       `xml:"DrillingCost>row"`
	DrillingCostItem   []DrillingCostItem   `xml:"DrillingCostItem>row"`
	FacilitiesCost     []FacilitiesCost     `xml:"FacilitiesCost>row"`
	FacilitiesCostItem []FacilitiesCostItem `xml:"FacilitiesCostItem>row"`
	JobDetails         []JobDetails         `xml:"JobDetails>row"`
	LocationInfo       []LocationInfo       `xml:"LocationInfo>row"`
	WellInfo           []WellInfo           `xml:"WellInfo>row"`
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
	UserDefined1       string  `json:"UserDefined_1" xml:"UserDefined_1,attr"`
	UserDefined2       string  `json:"UserDefined_2" xml:"UserDefined_2,attr"`
	UserDefined3       string  `json:"UserDefined_3" xml:"UserDefined_3,attr"`
	UserDefined4       string  `json:"UserDefined_4" xml:"UserDefined_4,attr"`
	UserDefined5       string  `json:"UserDefined_5" xml:"UserDefined_5,attr"`
	UserDefined6       string  `json:"UserDefined_6" xml:"UserDefined_6,attr"`
	UserDefined7       string  `json:"UserDefined_7" xml:"UserDefined_7,attr"`
	UserDefined8       string  `json:"UserDefined_8" xml:"UserDefined_8,attr"`
	UserDefined9       string  `json:"UserDefined_9" xml:"UserDefined_9,attr"`
	UserDefined10      string  `json:"UserDefined_10" xml:"UserDefined_10,attr"`
	UserDefined11      string  `json:"UserDefined_11" xml:"UserDefined_11,attr"`
	UserDefined12      string  `json:"UserDefined_12" xml:"UserDefined_12,attr"`
	UserDefined13      string  `json:"UserDefined_13" xml:"UserDefined_13,attr"`
	UserDefined14      string  `json:"UserDefined_14" xml:"UserDefined_14,attr"`
	UserDefined15      string  `json:"UserDefined_15" xml:"UserDefined_15,attr"`
	UserDefined16      string  `json:"UserDefined_16" xml:"UserDefined_16,attr"`
	UserDefined17      string  `json:"UserDefined_17" xml:"UserDefined_17,attr"`
	UserDefined18      string  `json:"UserDefined_18" xml:"UserDefined_18,attr"`
	UserDefined19      string  `json:"UserDefined_19" xml:"UserDefined_19,attr"`
	UserDefined20      string  `json:"UserDefined_20" xml:"UserDefined_20,attr"`
	UserDefined21      string  `json:"UserDefined_21" xml:"UserDefined_21,attr"`
	UserDefined22      string  `json:"UserDefined_22" xml:"UserDefined_22,attr"`
	UserDefined23      string  `json:"UserDefined_23" xml:"UserDefined_23,attr"`
	UserDefined24      string  `json:"UserDefined_24" xml:"UserDefined_24,attr"`
	UserDefined25      string  `json:"UserDefined_25" xml:"UserDefined_25,attr"`
	UserDefined26      string  `json:"UserDefined_26" xml:"UserDefined_26,attr"`
	UserDefined27      string  `json:"UserDefined_27" xml:"UserDefined_27,attr"`
	UserDefined28      string  `json:"UserDefined_28" xml:"UserDefined_28,attr"`
	UserDefined29      string  `json:"UserDefined_29" xml:"UserDefined_29,attr"`
	UserDefined30      string  `json:"UserDefined_30" xml:"UserDefined_30,attr"`
	NetRevenueInterest float64 `json:"NetRevenueInterest" xml:"NetRevenueInterest,attr"`
}

type LocationInfo struct {
	RowID             int64  `json:"RowID" xml:"RowID,attr"`
	Parentid          int64  `json:"Parentid" xml:"Parentid,attr"`
	AssetId           int64  `json:"asset_id" xml:"asset_id,attr"`
	TimeZoneDesc      string `json:"Time_Zone_Description" xml:"Time_Zone_Description,attr"`
	TimeZone          string `json:"Time_Zone" xml:"Time_Zone,attr"`
	Name              string `json:"Name" xml:"Name,attr"`
	PropertyNumber    string `json:"PropertyNumber" xml:"PropertyNumber,attr"`
	Lease             string `json:"Lease" xml:"Lease,attr"`
	CountyName        string `json:"CountyName" xml:"CountyName,attr"`
	StateName         string `json:"StateName" xml:"StateName,attr"`
	DistrictName      string `json:"DistrictName" xml:"DistrictName,attr"`
	SectionName       string `json:"SectionName" xml:"SectionName,attr"`
	TownshipName      string `json:"TownshipName" xml:"TownshipName,attr"`
	Range             string `json:"Range" xml:"Range,attr"`
	LatitudeLocation  string `json:"LatitudeLocation" xml:"LatitudeLocation,attr"`
	LongitudeLocation string `json:"LongitudeLocation" xml:"LongitudeLocation,attr"`
	GasPurchaser      string `json:"GasPurchaser" xml:"GasPurchaser,attr"`
	OilPurchaser      string `json:"OilPurchaser" xml:"OilPurchaser,attr"`
	LineNumber        int64  `json:"line_number" xml:"line_number,attr"`
	WaterPurchaser    string `json:"WaterPurchaser" xml:"WaterPurchaser,attr"`
	Deleted           string `json:"deleted" xml:"deleted,attr"`
	reuseEmailList    string `json:"reuseEmailList" xml:"reuseEmailList,attr"`
	Comments          string `json:"Comments" xml:"Comments,attr"`
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

type CostAllocation struct {
	RowID        int64  `json:"RowID" xml:"RowID,attr"`
	LocationID   int64  `json:"location_id" xml:"location_id,attr"`
	JobNumber    int64  `json:"job_number" xml:"job_number,attr"`
	ReportTypeID int64  `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	Comments     string `json:"Comments" xml:"Comments,attr"`
}

type CostAllocationItem struct {
	RowID        int64   `json:"RowID" xml:"RowID,attr"`
	LocationID   int64   `json:"location_id" xml:"location_id,attr"`
	JobNumber    int64   `json:"job_number" xml:"job_number,attr"`
	ReportTypeID int64   `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	LineNumber   int64   `json:"line_number" xml:"line_number,attr"`
	LocWellId    int64   `json:"LocWellId" xml:"LocWellId,attr"`
	LocWellJobId int64   `json:"LocWellJobId" xml:"LocWellJobId,attr"`
	AllocPrecent float64 `json:"AllocPercent" xml:"AllocPercent,attr"`
}

type DailyOps struct {
	WellID        int64   `json:"well_id" xml:"well_id,attr"`
	ReportDate    string  `json:"report_date" xml:"report_date,attr"`
	JobNumber     int64   `json:"job_number" xml:"job_number,attr"`
	WellStatus    string  `json:"WellStatus" xml:"WellStatus,attr"`
	UserDefined1  string  `json:"UserDefined_1" xml:"UserDefined_1,attr"`
	UserDefined2  string  `json:"UserDefined_2" xml:"UserDefined_2,attr"`
	UserDefined3  string  `json:"UserDefined_3" xml:"UserDefined_3,attr"`
	UserDefined4  string  `json:"UserDefined_4" xml:"UserDefined_4,attr"`
	UserDefined5  string  `json:"UserDefined_5" xml:"UserDefined_5,attr"`
	UserDefined6  string  `json:"UserDefined_6" xml:"UserDefined_6,attr"`
	UserDefined7  string  `json:"UserDefined_7" xml:"UserDefined_7,attr"`
	UserDefined8  string  `json:"UserDefined_8" xml:"UserDefined_8,attr"`
	UserDefined9  string  `json:"UserDefined_9" xml:"UserDefined_9,attr"`
	UserDefined10 string  `json:"UserDefined_10" xml:"UserDefined_10,attr"`
	UserDefined11 string  `json:"UserDefined_11" xml:"UserDefined_11,attr"`
	UserDefined12 string  `json:"UserDefined_12" xml:"UserDefined_12,attr"`
	UserDefined13 string  `json:"UserDefined_13" xml:"UserDefined_13,attr"`
	UserDefined14 string  `json:"UserDefined_14" xml:"UserDefined_14,attr"`
	UserDefined15 string  `json:"UserDefined_15" xml:"UserDefined_15,attr"`
	UserDefined16 string  `json:"UserDefined_16" xml:"UserDefined_16,attr"`
	UserDefined17 string  `json:"UserDefined_17" xml:"UserDefined_17,attr"`
	UserDefined18 string  `json:"UserDefined_18" xml:"UserDefined_18,attr"`
	UserDefined19 string  `json:"UserDefined_19" xml:"UserDefined_19,attr"`
	UserDefined20 string  `json:"UserDefined_20" xml:"UserDefined_20,attr"`
	UserDefined21 string  `json:"UserDefined_21" xml:"UserDefined_21,attr"`
	UserDefined22 string  `json:"UserDefined_22" xml:"UserDefined_22,attr"`
	UserDefined23 string  `json:"UserDefined_23" xml:"UserDefined_23,attr"`
	UserDefined24 string  `json:"UserDefined_24" xml:"UserDefined_24,attr"`
	UserDefined25 string  `json:"UserDefined_25" xml:"UserDefined_25,attr"`
	UserDefined26 string  `json:"UserDefined_26" xml:"UserDefined_26,attr"`
	UserDefined27 string  `json:"UserDefined_27" xml:"UserDefined_27,attr"`
	UserDefined28 string  `json:"UserDefined_28" xml:"UserDefined_28,attr"`
	UserDefined29 string  `json:"UserDefined_29" xml:"UserDefined_29,attr"`
	UserDefined30 string  `json:"UserDefined_30" xml:"UserDefined_30,attr"`
	TripGas       float64 `json:"TripGas" xml:"TripGas,attr"`
	TVD           float64 `json:"TVD" xml:"TVD,attr"`
	TMD           float64 `json:"TMD" xml:"TMD,attr"`
}

type DrillingCost struct {
	WellID       int64  `json:"well_id" xml:"well_id,attr"`
	ReportDate   string `json:"report_date" xml:"report_date,attr"`
	JobNumber    int64  `json:"job_number" xml:"job_number,attr"`
	ReportTypeID int64  `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	Comment      string `json:"Comment" xml:"Comment,attr"`
	LocationID   int64  `json:"location_id" xml:"location_id,attr"`
	RowID        int64  `json:"RowID" xml:"RowID,attr"`
}

type DrillingCostItem struct {
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

type FacilitiesCost struct {
	WellID       int64  `json:"well_id" xml:"well_id,attr"`
	ReportDate   string `json:"report_date" xml:"report_date,attr"`
	JobNumber    int64  `json:"job_number" xml:"job_number,attr"`
	ReportTypeID int64  `json:"ReportTypeID" xml:"ReportTypeID,attr"`
	Comment      string `json:"Comment" xml:"Comment,attr"`
	LocationID   int64  `json:"location_id" xml:"location_id,attr"`
	RowID        int64  `json:"RowID" xml:"RowID,attr"`
}

type FacilitiesCostItem struct {
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

type JobDetails struct {
	WellID          int64   `json:"well_id" xml:"well_id,attr"`
	JobStatus       int64   `json:"job_status" xml:"job_status,attr"`
	JobNumber       int64   `json:"job_number" xml:"job_number,attr"`
	JobDescription  string  `json:"job_description" xml:"jobDescription,attr"`
	JobType         string  `json:"JobType" xml:"JobType,attr"`
	WorkingInterest float64 `json:"WorkingInterest" xml:"WorkingInterest,attr"`
	WellType        string  `json:"WellType" xml:"WellType,attr"`
	WelLEngineer    string  `json:"WellEngineer" xml:"WellEngineer,attr"`
	WellBore        string  `json:"WellBore" xml:"WellBore,attr"`
}
