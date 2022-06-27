package stryd

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.39 Safari/537.36"
)

type Token struct {
	UserName     string `json:"user_name"`
	Token        string `json:"token"`
	ID           string `json:"id"`
	Weight       int    `json:"weight"`
	Height       int    `json:"height"`
	Unit         string `json:"unit"`
	ProfileImage string `json:"profile_image"`
}

type Client struct {
	Email    string
	Password string
	Token    *Token
}

type Zone struct {
	PowerLow  float64 `json:"power_low"`
	PowerHigh float64 `json:"power_high"`
	SpeedLow  float64 `json:"speed_low"`
	SpeedHigh float64 `json:"speed_high"`
	Name      string  `json:"name"`
}

type GPSPoint struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}

type Activity struct {
	ID                   int64    `json:"id"`
	WorkoutID            string   `json:"workout_id"`
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Rpe                  int      `json:"rpe"`
	Feel                 string   `json:"feel"`
	Type                 string   `json:"type"`
	Source               string   `json:"source"`
	SurfaceType          string   `json:"surface_type"`
	Weight               int      `json:"weight"`
	Ftp                  float64  `json:"ftp"`
	Zones                []Zone   `json:"zones"`
	StartPoint           GPSPoint `json:"start_point"`
	EndPoint             GPSPoint `json:"end_point"`
	TimeStamp            int64    `json:"timestamp"`
	StartTime            int64    `json:"start_time"`
	MovingTime           int64    `json:"moving_time"`
	ElapsedTime          int64    `json:"elapsed_time"`
	ClockTime            int64    `json:"clock_time"`
	TimeZone             string   `json:"time_zone"`
	LapTimestampList     []int64  `json:"lap_timestamp_list"`
	Distance             float64  `json:"distance"`
	TotalElevationGain   float64  `json:"total_elevation_gain"`
	TotalElevationLoss   float64  `json:"total_elevation_loss"`
	MaxElevation         float64  `json:"max_elevation"`
	MinElevation         float64  `json:"min_elevation"`
	AverageSpeed         float64  `json:"average_speed"`
	MaxSpeed             float64  `json:"max_speed"`
	AverageCadence       float64  `json:"average_cadence"`
	MaxCadence           float64  `json:"max_cadence"`
	MinCadence           float64  `json:"min_cadence"`
	AverageStrideLength  float64  `json:"average_stride_length"`
	MaxStrideLegth       float64  `json:"max_stride_length"`
	MinStrideLength      float64  `json:"min_stride_length"`
	AverageGroundTime    float64  `json:"average_ground_time"`
	MaxGroundTime        float64  `json:"max_ground_time"`
	MinGroundTime        float64  `json:"min_ground_time"`
	AverageOscillation   float64  `json:"average_oscillation"`
	MaxOscillation       float64  `json:"max_oscillation"`
	MinOscillation       float64  `json:"min_oscillation"`
	AveragePower         float64  `json:"average_power"`
	MaxPower             float64  `json:"max_power"`
	AverageHeartRate     float64  `json:"average_heart_rate"`
	MaxHeartRate         float64  `json:"max_heart_rate"`
	AverageLegSpring     float64  `json:"average_leg_spring"`
	Calories             float64  `json:"calories"`
	Stress               float64  `json:"stress"`
	MaxVerticalStiffness float64  `json:"max_vertical_stiffness"`
	SecondsInZones       []int    `json:"seconds_in_zones"`
	Public               bool     `json:"public"`
	Deleted              bool     `json:"deleted"`
	Favorite             bool     `json:"favorite"`
	Flagged              bool     `json:"flagged"`
	Pending              bool     `json:"pending"`
	Excluded             bool     `json:"excluded"`
	Elevation            float64  `json:"elevation"`
	Temperature          float64  `json:"temperature"`
	Humidity             float64  `json:"humidity"`
	DewPoint             float64  `json:"dewPoint"`
	WindBearing          float64  `json:"windBearing"`
	WindSpeed            float64  `json:"windSpeed"`
	WindGust             float64  `json:"windGust"`
	Icon                 string   `json:"icon"`
	LapEvents            []int    `json:"lap_events"`
	StartEvents          []int    `json:"start_events"`
	StopEvents           []int    `json:"stop_events"`
	FilePath             string   `json:"file_path"`
	UpdatedTime          string   `json:"updated_time"`
	UserID               string   `json:"user_id"`
}

func NewClient(email string, password string) *Client {
	var client Client
	client.Email = email
	client.Password = password

	return &client
}

func (this *Client) Login() (*Token, error) {

	api := "https://www.stryd.com/b/email/signin"

	headers := map[string]string{
		"Accept":       "application/json",
		"Referer":      "https://www.stryd.com/signin",
		"User-Agent":   UA,
		"Content-Type": "application/json",
	}

	param := map[string]interface{}{
		"email":    this.Email,
		"password": this.Password,
	}

	respBody, err := HttpPost(api, headers, param)

	if err != nil {
		return nil, err
	}

	var token Token

	err = json.Unmarshal(respBody, &token)

	if err != nil {
		return nil, err
	}

	this.Token = &token

	return &token, nil
}

func (this *Client) Activities(updatedAfter int64, includeDeleted bool) ([]Activity, error) {

	api := fmt.Sprintf("https://www.stryd.com/b/api/v1/users/calendar?updated_after=%d&include_deleted=%v", updatedAfter, includeDeleted)
	log.Println(api)

	headers := map[string]string{
		"authority":     "www.stryd.com",
		"accept":        "application/json, text/plain, */*",
		"authorization": fmt.Sprintf("Bearer: %s", this.Token.Token),
		"User-Agent":    UA,
		"referer":       "https://www.stryd.com/powercenter/profile'",
	}

	respBody, err := HttpGet(api, headers, nil)

	if err != nil {
		return nil, err
	}

	type Resp struct {
		Activities []Activity `json:"activities"`
	}

	var resp Resp

	err = json.Unmarshal(respBody, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Activities, nil
}

func (this *Client) GetDownloadUrl(activityID int64) (string, error) {

	api := fmt.Sprintf("https://www.stryd.com/b/api/v1/activities/%d/fit", activityID)

	headers := map[string]string{
		"authority":     "www.stryd.com",
		"Accept":        "application/json, text/plain, */*'",
		"authorization": fmt.Sprintf("Bearer: %s", this.Token.Token),
		"origin":        "https://www.stryd.com",
		"Referer":       "https://www.stryd.com/powercenter/profile",
		"User-Agent":    UA,
	}

	respBody, err := HttpPost(api, headers, nil)

	if err != nil {
		return "", err
	}

	type Resp struct {
		URL string `json:"url"`
	}

	var resp Resp

	err = json.Unmarshal(respBody, &resp)

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
