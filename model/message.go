package model

import (
	"encoding/xml"
)

type ReceiveMessageResponse struct {
	XMLName              xml.Name             `xml:"ReceiveMessageResponse"`
	ReceiveMessageResult ReceiveMessageResult `xml:"ReceiveMessageResult"`
	ResponseMetadata     ResponseMetadata     `xml:"ResponseMetadata"`
}

type ReceiveMessageResult struct {
	XMLName xml.Name `xml:"ReceiveMessageResult"`
	Message *Message `xml:"Message"`
}

type Response struct {
	RequestId     string `json:"-"`
	MessageId     string `json:"-"`
	ReceiptHandle string `json:"-"`
	MD5OfBody     string `json:"-"`
	UserID        string `json:"user_id"`
	AppVersion    string `json:"app_version"`
	DeviceType    string `json:"device_type"`
	IP            string `json:"ip"`
	Locale        string `json:"locale"`
	DeviceID      string `json:"device_id"`
}

type Message struct {
	XMLName       xml.Name `xml:"Message"`
	MessageId     string   `xml:"MessageId"`
	ReceiptHandle string   `xml:"ReceiptHandle"`
	MD5OfBody     string   `xml:"MD5OfBody"`
	Body          string   `xml:"Body"`
}

type ResponseMetadata struct {
	XMLName   xml.Name `xml:"ResponseMetadata"`
	RequestId string   `xml:"RequestId"`
}

// SetData sets the data fields of the Response struct based on the SQS message response.
func (res *Response) SetData(sqsMessageResponse ReceiveMessageResponse) {
	res.RequestId = sqsMessageResponse.ResponseMetadata.RequestId
	res.MessageId = sqsMessageResponse.ReceiveMessageResult.Message.MessageId
	res.ReceiptHandle = sqsMessageResponse.ReceiveMessageResult.Message.ReceiptHandle
	res.MD5OfBody = sqsMessageResponse.ReceiveMessageResult.Message.MD5OfBody
}

// MaskBody encrypts the DeviceID and IP fields of the Response struct if they are not empty.
func (res *Response) MaskBody() error {
	if res.DeviceID != "" {
		encryptedDeviceId, err := encrypt(res.DeviceID)
		if err != nil {
			return err
		}

		res.DeviceID = *encryptedDeviceId
	}

	if res.IP != "" {
		encryptedIP, err := encrypt(res.IP)
		if err != nil {
			return err
		}

		res.IP = *encryptedIP
	}

	return nil
}
