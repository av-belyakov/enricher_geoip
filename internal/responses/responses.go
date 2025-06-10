package responses

type Response struct {
	FoundInformation []DetailedInformation `json:"found_information"`
	TaskId           string                `json:"task_id"`
	Source           string                `json:"source"`
	Error            string                `json:"error"`
}

type DetailedInformation struct {
	IpRange struct {
		IpFirst string `json:"ip_first"`
		IpLast  string `json:"ip_last"`
	} `json:"ip_range"`
	IpAddr    string `json:"ip_address"`
	Code      string `json:"code"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Subnet    string `json:"subnet"`
	UpdatedAt string `json:"updated_at"`
	Error     string `json:"error"`
}

type ResponseGeoIPDataBase struct {
	IpLocations             []IpAddrLocation `json:"ip_locations"`
	InternetProtocolVersion string           `json:"address_version"`
}

type IpAddrLocation struct {
	Asns        []any  `json:"asns"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	IpRange     struct {
		IpFirst string `json:"ip_first"`
		IpLast  string `json:"ip_last"`
	} `json:"ip_range"`
	Isp       string  `json:"isp"`
	Region    string  `json:"region"`
	Source    string  `json:"source"`
	Subnet    string  `json:"subnet"`
	UpdatedAt string  `json:"updated_at"`
	Rating    float64 `json:"rating"`
}
