package modifierleasingbodyapprove

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/martian/parse"
)

func init() {
	parse.Register("bodyleasingapprove.Modifier", modifierFromJSON)
}

type XmlModifier struct {
	contentType string
}

type XmlModifierJSON struct {
	ContentType string               `json:"contentType"`
	Scope       []parse.ModifierType `json:"scope"`
}

type Request struct {
	XMLName           xml.Name       `xml:"Request" json:"-"`
	ApplicationId     string         `xml:"applicationId" json:"applicationId"`
	CreditScore       string         `xml:"creditScore" json:"creditScore"`
	PartiallyApproved bool           `xml:"partiallyApproved" json:"partiallyApproved"`
	ProductDesc       string         `xml:"productDesc" json:"productDesc"`
	RiskRating        string         `xml:"riskRating" json:"riskRating"`
	SafrapayId        string         `xml:"safrapayId" json:"safrapayId"`
	AdverseAction     *AdverseAction `xml:"adverseAction" json:"adverseAction"`
	LeasingProducts       [] *LeasingProducts    `xml:"leasingProducts" json:"leasingProducts"`
}

type LeasingProducts struct {
	ProductId      string  `xml:"productId" json:"productId"`
	Funding  float32 `xml:"funding" json:"funding"`
	Duration int32 `xml:"duration" json:"duration"`
	Quantity int32 `xml:"quantity" json:"quantity"`
	Manufacturer      string  `xml:"manufacturer" json:"manufacturer"`
	Mid      string  `xml:"mid" json:"mid"`
	Model      string  `xml:"model" json:"model"`
	TotalMonthlyPayment  float32 `xml:"totalMonthlyPayment" json:"totalMonthlyPayment"`
	CalculateTax  float32 `xml:"calculateTax" json:"calculateTax"`
	LineItemNumber  int32 `xml:"lineItemNumber" json:"lineItemNumber"`
}

type AdverseAction struct {
	Code        string `xml:"code" json:"code"`
	Source      string `xml:"source" json:"source"`
	Description string `xml:"description" json:"description"`
}

func (m *XmlModifier) ModifyRequest(req *http.Request) error {
	req.Header.Set("Content-Type", m.contentType)

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return err
	}

	var data Request
	xml.Unmarshal([]byte(body), &data)
	jsonData, _ := json.Marshal(data)

	req.ContentLength = int64(len(jsonData))
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
	req.URL.Path = strings.Replace(req.URL.Path, "safrapayId", data.SafrapayId, 1)

	return nil
}

func (m *XmlModifier) ModifyResponse(res *http.Response) error {
	
	if(res.StatusCode != 200){
		res.StatusCode = 200
	}	

	return nil
}

func XmlNewModifier(contentType string) *XmlModifier {
	return &XmlModifier{
		contentType: contentType,
	}
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &XmlModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(XmlNewModifier(msg.ContentType), msg.Scope)
}
