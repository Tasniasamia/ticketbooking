package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	sslSandboxURL = "https://sandbox.sslcommerz.com/gwprocess/v4/api.php"
	sslLiveURL    = "https://securepay.sslcommerz.com/gwprocess/v4/api.php"
	sslSandboxVal = "https://sandbox.sslcommerz.com/validator/api/validationserverAPI.php"
	sslLiveVal    = "https://securepay.sslcommerz.com/validator/api/validationserverAPI.php"
)

type sslcommerzClient struct {
	storeID, storePass string
	isSandbox          bool
}

func newSSLCommerzClient(id, pass string, sandbox bool) *sslcommerzClient {
	return &sslcommerzClient{storeID: id, storePass: pass, isSandbox: sandbox}
}

type sslInitParams struct {
	TotalAmount, Currency, TranID, SuccessURL, FailURL, CancelURL, IPNURL string
	CusName, CusEmail, CusPhone, CusAdd1, CusCity, CusCountry, CusPostcode string
	ProductName, ProductCategory, ProductProfile                           string
	ValueA, ValueB, ValueC, ValueD                                         string
}

type sslInitResponse struct {
	Status, FailedReason, SessionKey, GatewayPageURL, RedirectGatewayURL string
}

func (c *sslcommerzClient) Initiate(p sslInitParams) (*sslInitResponse, error) {
	ep := sslSandboxURL
	if !c.isSandbox {
		ep = sslLiveURL
	}
	form := url.Values{}
	form.Set("store_id", c.storeID)
	form.Set("store_passwd", c.storePass)
	form.Set("total_amount", p.TotalAmount)
	form.Set("currency", p.Currency)
	form.Set("tran_id", p.TranID)
	form.Set("success_url", p.SuccessURL)
	form.Set("fail_url", p.FailURL)
	form.Set("cancel_url", p.CancelURL)
	if p.IPNURL != "" {
		form.Set("ipn_url", p.IPNURL)
	}
	form.Set("cus_name", p.CusName)
	form.Set("cus_email", p.CusEmail)
	form.Set("cus_phone", p.CusPhone)
	form.Set("cus_add1", p.CusAdd1)
	form.Set("cus_city", p.CusCity)
	form.Set("cus_country", p.CusCountry)
	if p.CusPostcode != "" {
		form.Set("cus_postcode", p.CusPostcode)
	}
	form.Set("product_name", p.ProductName)
	form.Set("product_category", p.ProductCategory)
	form.Set("product_profile", p.ProductProfile)
	form.Set("shipping_method", "NO")
	form.Set("num_of_item", "1")
	form.Set("value_a", p.ValueA)
	form.Set("value_b", p.ValueB)
	form.Set("value_c", p.ValueC)
	form.Set("value_d", p.ValueD)
	resp, err := http.Post(ep, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result sslInitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("sslcommerz parse: %w body=%s", err, string(body))
	}
	if !strings.EqualFold(result.Status, "SUCCESS") {
		r := result.FailedReason
		if r == "" {
			r = string(body)
		}
		return nil, fmt.Errorf("sslcommerz failed: %s", r)
	}
	if result.GatewayPageURL == "" {
		result.GatewayPageURL = result.RedirectGatewayURL
	}
	return &result, nil
}

type sslValidateResponse struct {
	Status, TranID, ValID, Amount, Currency, BankTranID string
}

func (c *sslcommerzClient) Validate(valID string) (*sslValidateResponse, error) {
	ep := sslSandboxVal
	if !c.isSandbox {
		ep = sslLiveVal
	}
	q := url.Values{}
	q.Set("val_id", valID)
	q.Set("store_id", c.storeID)
	q.Set("store_passwd", c.storePass)
	q.Set("format", "json")
	resp, err := http.Get(ep + "?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result sslValidateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
