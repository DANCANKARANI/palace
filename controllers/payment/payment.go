package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

//mpesa is an application that will ve making transaction
type Mpesa struct{
	ConsumerKey string
	ConsumerSecret string
	BaseURL string
	Client	*http.Client
}

//MpesaOpts stores all the configuration keys we need to set up a Mpesa app
type MpesaOpts struct {
	ConsumerKey string
	ConsumerSecret string
	BaseURL string
}

//NewMpesa sets up and returns an instance of Mpesa
func NewMpesa(m *MpesaOpts) *Mpesa{
	client := &http.Client{
		Timeout: 30 *time.Second,
	}

	return &Mpesa{
		ConsumerKey: m.ConsumerKey,
		ConsumerSecret: m.ConsumerSecret,
		BaseURL: m.BaseURL,
		Client: client,
	}
}

//makeRequest Performs all the http request for the specific app
func (m *Mpesa) MakeRequest(req *http.Request)([]byte,error){
	resp, err := m.Client.Do(req)
	if err != nil{
		return nil, err
	}

	defer func(Body io.ReadCloser){
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	return body, nil
}

//MpesaAccessTokenResponse is the response sent back by safarico  when we make a request to generate a token
type MpesaAccessTokenResponse struct{
	AccessToken string `json:"access_token"`
	ExpiresIn string	`json:"expires_in"`
	RequestID	string	`json:"request_id"`
	ErrorCode   string	`json:"error_code"`
	ErrorMessage	string `json:"error_message"`
}

//generateAccessToken sends a http request to generate new access token
func (m *Mpesa) GenerateAccessToken() (*MpesaAccessTokenResponse,error){
	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials",m.BaseURL)

	req, err := http.NewRequest(http.MethodGet,url,nil)
	if err != nil{
		return nil, errors.New("error sending request:"+err.Error())
	}

	req.SetBasicAuth(m.ConsumerKey,m.ConsumerSecret)
	req.Header.Set("Content-type","application/json")

	resp, err := m.MakeRequest(req)
	if err != nil{
		return nil, errors.New("error making request:"+err.Error())
	}

	accessTokenResponse := new(MpesaAccessTokenResponse)
	if err := json.Unmarshal(resp, accessTokenResponse); err != nil {
		return nil, errors.New("error unmarshalling: " + err.Error())
	}

	return accessTokenResponse,nil
}

//STKPushRequestBody is the body with the parameters to be used to initiate an STK push Request request
type STKPushRequestBody struct{
	BusinessShortCode		string	`json:"business_short_code"`
	Password				string	`json:"password"`
	Timestamp				string	`json:"timestamp"`
	TransactionType			string	`json:"transaction_type"`
	Amount					string	`json:"amount"`
	PartyA					string	`json:"party_a"`
	PartyB					string	`json:"party_b"`
	PhoneNumber				string	`json:"phone_number"`
	CallBackURL				string	`json:"call_back_url"`
	AccountReference		string	`json:"account_reference"`
	TransactionDesc			string	`json:"transaction_desc"`
}

//STKPushRequestResponse is the response sent back after an STK push request
type STKPushRequestResponse struct{
	MerchantRequestID		string 	`json:"merchant_request_id"`
	CheckOutRequestID		string	`json:"check_out_request_id"`
	ResponseCode			string	`json:"response_code"`
	ResponseDescription		string	`json:"response_description"`
	CustomerMessage			string	`json:"customer_message"`
	RequestID				string	`json:"request_id"`
	ErrorCode				string	`json:"error_code"`
	ErrorMessage			string	`json:"error_message"`
}

//InitiateSTKPushRequest makes a http request perfomming an STK push request
func (m *Mpesa) InitiateSTKPushRequest(body *STKPushRequestBody)(*STKPushRequestResponse,error){
	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest",m.BaseURL)

	requestBody, err := json.Marshal(body)
	if err != nil{
		return nil, errors.New("error parsong request body:"+err.Error())
	}

	req, err := http.NewRequest(http.MethodPost,url,bytes.NewBuffer(requestBody))
	if err != nil{
		return nil,errors.New("error sending request:"+err.Error())
	}

	accessTokenResponse , err := m.GenerateAccessToken()
	if err != nil{
		return nil,errors.New("error generating access token:"+err.Error())
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Authorization",fmt.Sprintf("Bearer %s",accessTokenResponse.AccessToken))

	resp, err := m.MakeRequest(req)
	if err != nil{
		return nil, errors.New("error making request:"+err.Error())
	}
	log.Printf("Raw STK Push Response: %s", string(resp))

	stkPushResponse := new(STKPushRequestResponse)
	if err := json.Unmarshal(resp,&stkPushResponse); err != nil{
		return nil, errors.New("error unmarshalling the body:"+err.Error())
	}

	return stkPushResponse,nil
}