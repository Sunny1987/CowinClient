package main

//vaccines: COVISHIELD,
//Karnataka :

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

//***********************************************VARIABLES********************************************
//All important variables listed here mostly URLs
var (
	Client     *resty.Client = resty.New()
	ConfirmOTP string        = "https://cdn-api.co-vin.in/api/v2/auth/public/confirmOTP"
	GetOTP     string        = "https://cdn-api.co-vin.in/api/v2/auth/public/generateOTP"
	GETSTATES  string        = "https://cdn-api.co-vin.in/api/v2/admin/location/states"
	GetDist    string        = "https://cdn-api.co-vin.in/api/v2/admin/location/districts/"
	GetCenters string        = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByDistrict"
	Token      string        = ""
	Bearer     string        = "Bearer "
	FindByPin  string        = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByPin"
)

//****************************************************************************************************

//*****************************************JSON STRUCTURES********************************************
//States is the structure for states response
type States struct {
	States []State `json:"states"`
	TTL    int     `json:"ttl"`
}

//State is the individual state structure
type State struct {
	StateID   int    `json:"state_id"`
	StateName string `json:"state_name"`
}

type Districts struct {
	Districts []District `json:"districts"`
	TTL       int        `json:"ttl"`
}

type District struct {
	DistrictID   int    `json:"district_id"`
	DistrictName string `json:"district_name"`
}

type Centers struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	CenterID               int      `json:"center_id"`
	Name                   string   `json:"name"`
	Address                string   `json:"address"`
	StateName              string   `json:"state_name"`
	DistrictName           string   `json:"district_name"`
	BlockName              string   `json:"block_name"`
	Pincode                int      `json:"pincode"`
	From                   string   `json:"from"`
	To                     string   `json:"to"`
	Lat                    int      `json:"lat"`
	Long                   int      `json:"long"`
	FeeType                string   `json:"fee_type"`
	SessionID              string   `json:"session_id"`
	Date                   string   `json:"date"`
	AvailableCapacityDose1 int      `json:"available_capacity_dose1"`
	AvailableCapacityDose2 int      `json:"available_capacity_dose2"`
	AvailableCapacity      int      `json:"available_capacity"`
	Fee                    string   `json:"fee"`
	MinAgeLimit            int      `json:"min_age_limit"`
	Vaccine                string   `json:"vaccine"`
	Slots                  []string `json:"slots"`
}

type CentersByPIN struct {
	Sessions []CenterByPIN `json:"sessions"`
}

type CenterByPIN struct {
	CenterID               int      `json:"center_id"`
	Name                   string   `json:"name"`
	NameL                  string   `json:"name_l"`
	Address                string   `json:"address"`
	AddressL               string   `json:"address_l"`
	StateName              string   `json:"state_name"`
	StateNameL             string   `json:"state_name_l"`
	DistrictName           string   `json:"district_name"`
	DistrictNameL          string   `json:"district_name_l"`
	BlockName              string   `json:"block_name"`
	BlockNameL             string   `json:"block_name_l"`
	Pincode                string   `json:"pincode"`
	Lat                    float64  `json:"lat"`
	Long                   float64  `json:"long"`
	From                   string   `json:"from"`
	To                     string   `json:"to"`
	FeeType                string   `json:"fee_type"`
	Fee                    string   `json:"fee"`
	SessionID              string   `json:"session_id"`
	Date                   string   `json:"date"`
	AvailableCapacity      int      `json:"available_capacity"`
	AvailableCapacityDose1 int      `json:"available_capacity_dose1"`
	AvailableCapacityDose2 int      `json:"available_capacity_dose2"`
	MinAgeLimit            int      `json:"min_age_limit"`
	Vaccine                string   `json:"vaccine"`
	Slots                  []string `json:"slots"`
}

type GetOTPResp struct {
	TxnId string
}

type ConfirmMyOTPReq struct {
	Otp   string
	TxnId string
}

type ConfirmOTPRespToken struct {
	Token        string
	IsNewAccount string
}

//***************************************************************************************************

//*****************************************STRUCTURES METHODS****************************************
func (re *GetOTPResp) getOTP(mobile string) {
	ot := map[string]string{"mobile": mobile}
	resp, err := Client.R().SetHeader("Accept", "application/json").SetBody(ot).Post(GetOTP)

	if err != nil {
		log.Printf("Error fetching otp: %v", err)
	}

	json.Unmarshal(resp.Body(), re)

}

func (conf *ConfirmOTPRespToken) ConfirmMyOTP(txnId string, otpV string) string {

	//log.Println(otpV)
	//log.Println(txnId)
	h := sha256.New()
	h.Write([]byte(otpV))
	otpData := hex.EncodeToString(h.Sum(nil))
	// otpData := BytesToString(h.Sum(nil))
	//log.Printf("otp: %s", otpData)
	otp := map[string]interface{}{"otp": otpData, "txnId": txnId}
	//log.Println(otp)

	resp, err := Client.R().SetHeader("Accept", "application/json").SetBody(otp).Post(ConfirmOTP)
	if err != nil {
		log.Printf("Error fetching otp: %v", err)
	}
	//log.Println(resp)
	json.Unmarshal(resp.Body(), conf)
	//log.Printf("re: %v", conf.Token)
	Token = conf.Token
	Bearer = Bearer + Token
	//log.Printf("Bearer: %v", Bearer)
	return Bearer
}

func (pin *CentersByPIN) getCentersByPIN(pinCode string, date string, token string) {
	log.Printf("pincode : %s\n", pinCode)
	log.Printf("date : %s\n", date)
	resp, err := Client.R().SetQueryParams(map[string]string{
		"pincode": pinCode,
		"date":    date,
	}).SetHeader("Accept", "application/json").SetAuthToken(token).Get(FindByPin)
	if err != nil {
		log.Printf("ERROR in getting response: %v", err)
	}
	//log.Println(resp.Body())
	json.Unmarshal(resp.Body(), pin)
}

func (sts *States) getStates() {
	resp, err := Client.R().SetHeader("Accept", "application/json").SetAuthToken(Bearer).Get(GETSTATES)
	if err != nil {
		log.Printf("ERROR in getting response: %v", err)
	}
	json.Unmarshal(resp.Body(), sts)

}

func (dis *Districts) getDistricts(stateId string) {
	GetDist = GetDist + stateId
	resp, err := Client.R().SetHeader("Accept", "application/json").SetAuthToken(Bearer).Get(GetDist)
	if err != nil {
		log.Printf("ERROR in getting response: %v", err)
	}
	json.Unmarshal(resp.Body(), dis)

}

func (c *Centers) getCenters(district_id string, date string, vaccine string) {
	resp, err := Client.R().SetQueryParams(map[string]string{
		"district_id": district_id,
		"date":        date,
		"vaccine":     vaccine,
	}).SetHeader("Accept", "application/json").SetAuthToken(Bearer).Get(GetCenters)
	if err != nil {
		log.Printf("ERROR in getting response: %v", err)
	}
	json.Unmarshal(resp.Body(), c)
}

//***************************************************************************************************

//*****************************************MAIN******************************************************
func main() {
	log.Println("Choose options")
	log.Println("*********************************************")
	log.Println("1. No OTP info")
	log.Println("2. OTP based info")
	log.Println("*********************************************")
	var ch string
	fmt.Print(">>")
	fmt.Scanf("%s\n", &ch)
	log.Printf("The choice: %v", ch)

	switch ch {
	case "1":
		CaseA()
	case "2":
		CaseB()
	}

}

//***************************************************************************************************

//*****************************************FUNCTIONS*************************************************
func NoOTP(dStateName string, dDistrictName string, dDate string, dVac string) {
	log.Println("*********Calling non OTP based Center Information*********")

	// dStateName := flag.String("s", "Karnataka", "The required state..Default: Karnataka")
	// dDistrictName := flag.String("d", "BBMP", "The desired district..Default: BBMP")
	// dDate := flag.String("dt", "24-05-2021", "The desired date..Default: 24-05-2021")
	// dVac := flag.String("v", "COVISHIELD", "The desired vaccine..Default: COVISHIELD")
	// flag.Parse()

	var states States
	//get the states data
	states.getStates()

	//provide the value of the states
	for _, state := range states.States {

		//filter for state
		if state.StateName == dStateName {
			log.Printf("State: %s", state.StateName)
			var districts Districts
			log.Printf("State id : %s", strconv.Itoa(state.StateID))

			//get the list of districts
			districts.getDistricts(strconv.Itoa(state.StateID))
			for _, district := range districts.Districts {

				// filter for district
				if district.DistrictName == dDistrictName {
					log.Printf("District ID: %d", district.DistrictID)

					var centers Centers
					//get the centers for district , date, vaccine
					centers.getCenters(strconv.Itoa(district.DistrictID), dDate, dVac)
					if len(centers.Sessions) == 0 {
						log.Println("No Vaccines available for the location")
					} else {
						for _, session := range centers.Sessions {
							log.Println("**************************************************************************")
							log.Printf("Date: %v", session.Date)
							log.Printf("Address: %v", session.Address)
							log.Printf("Name: %v", session.Name)
							log.Printf("State: %v", session.StateName)
							log.Printf("District: %v", session.DistrictName)
							log.Printf("Block: %v", session.BlockName)
							log.Printf("Block: %v", session.BlockName)
							log.Printf("PIN: %v", session.Pincode)
							log.Printf("From: %v", session.From)
							log.Printf("To: %v", session.To)
							log.Printf("Fee Type: %v", session.FeeType)
							log.Printf("Fees: %v", session.Fee)
							log.Printf("Dosage 1: %v", session.AvailableCapacityDose1)
							log.Printf("Dosage 2: %v", session.AvailableCapacityDose2)
							log.Printf("Age Limit %v", session.MinAgeLimit)
							log.Printf("Vaccine %v", session.Vaccine)
							log.Printf("Slots %v", session.Slots)
							log.Println("**************************************************************************")
						}
					}
				}

			}
		}

	}
}

func CaseA() {
	fmt.Println("Enter State:")
	var st string
	fmt.Scanf("%s\n", &st)
	fmt.Println("Enter District:")
	var dt string
	fmt.Scanf("%s\n", &dt)
	fmt.Println("Enter Vaccine:")
	var vc string
	fmt.Scanf("%s\n", &vc)
	fmt.Println("Enter Date (dd-mm-yyyy):")
	var dDate string
	fmt.Scanf("%s\n", &dDate)
	fmt.Println("******************************************")
	fmt.Println()
	NoOTP(st, dt, dDate, vc)
}

func CaseB() {
	log.Println("*********Authentication Started*****************")
	var otp GetOTPResp
	var mob string
	var otpG string
	fmt.Println("Enter mobile:")
	fmt.Scanf("%s\n", &mob)
	otp.getOTP(mob)
	var conf ConfirmOTPRespToken
	fmt.Println("Enter OTP:")
	fmt.Scanf("%s\n", &otpG)
	bt := conf.ConfirmMyOTP(otp.TxnId, otpG)
	btarr := strings.Split(bt, " ")
	if len(btarr) > 1 {
		fmt.Println("***********Authentication Passed*******")
		var pin string
		var date string
		fmt.Println("Enter PIN for desired location:")
		fmt.Scanf("%s\n", &pin)
		fmt.Println("Enter Date (dd-mm-yyyy):")
		fmt.Scanf("%s\n", &date)
		var pinCenters CentersByPIN
		pinCenters.getCentersByPIN(pin, date, bt)
		displayCenters(&pinCenters)
	} else {
		fmt.Println("***********Authentication Failed*******")
	}

}

func displayCenters(pin *CentersByPIN) {
	if len(pin.Sessions) == 0 {
		log.Println("No Vaccines available for the location")
	} else {
		for _, session := range pin.Sessions {
			log.Println("**************************************************************************")
			log.Printf("Date: %v", session.Date)
			log.Printf("Address: %v", session.Address)
			log.Printf("Name: %v", session.Name)
			log.Printf("State: %v", session.StateName)
			log.Printf("District: %v", session.DistrictName)
			log.Printf("Block: %v", session.BlockName)
			log.Printf("Block: %v", session.BlockName)
			log.Printf("PIN: %v", session.Pincode)
			log.Printf("From: %v", session.From)
			log.Printf("To: %v", session.To)
			log.Printf("Fee Type: %v", session.FeeType)
			log.Printf("Fees: %v", session.Fee)
			log.Printf("Dosage 1: %v", session.AvailableCapacityDose1)
			log.Printf("Dosage 2: %v", session.AvailableCapacityDose2)
			log.Printf("Age Limit %v", session.MinAgeLimit)
			log.Printf("Vaccine %v", session.Vaccine)
			log.Printf("Slots %v", session.Slots)
			log.Println("**************************************************************************")
		}
	}

}
