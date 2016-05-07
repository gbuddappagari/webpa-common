package convey

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEncoding(t *testing.T) {
	assertions := assert.New(t)

	var testData = []struct {
		encoding *base64.Encoding
		value    string
		expected Payload
	}{
		{
			encoding: base64.StdEncoding,
			value:    "eyAicGFyYW1ldGVycyI6IFsgeyAibmFtZSI6ICJEZXZpY2UuRGV2aWNlSW5mby5XZWJwYS5YX0NPTUNBU1QtQ09NX0NJRCIsICJ2YWx1ZSI6ICIwIiwgImRhdGFUeXBlIjogMCB9LCB7ICJuYW1lIjogIkRldmljZS5EZXZpY2VJbmZvLldlYnBhLlhfQ09NQ0FTVC1DT01fQ01DIiwgInZhbHVlIjogIjI2OSIsICJkYXRhVHlwZSI6IDIgfSBdIH0K",
			expected: Payload{
				"parameters": []Payload{
					Payload{
						"name":     "Device.DeviceInfo.Webpa.X_COMCAST-COM_CID",
						"value":    "0",
						"dataType": 0,
					},
					Payload{
						"name":     "Device.DeviceInfo.Webpa.X_COMCAST-COM_CMC",
						"value":    "269",
						"dataType": 2,
					},
				},
			},
		},
	}

	for _, record := range testData {
		var actual Payload
		if err := actual.DecodeBase64(record.encoding, record.value); err != nil {
			t.Errorf("DecodeBase64 failed: %v", err)
		}

		expectedJson, err := json.Marshal(record.expected)
		if err != nil {
			t.Fatalf("Unable to marshal expected JSON: %v", err)
		}

		actualJson, err := json.Marshal(actual)
		if err != nil {
			t.Fatalf("Unable to marshal actual JSON: %v", err)
		}

		assertions.JSONEq(string(expectedJson), string(actualJson))
	}

	for _, record := range testData {
		// perform the reverse test: use the expected as our actual JSON
		actualEncoded, err := record.expected.EncodeBase64(record.encoding)
		if err != nil {
			t.Fatalf("Unable to encode: %v", err)
		}

		actualInput := bytes.NewBufferString(actualEncoded)
		actualDecoder := base64.NewDecoder(record.encoding, actualInput)
		actualDecoded, err := ioutil.ReadAll(actualDecoder)
		if err != nil {
			t.Fatalf("Unable to decode the output of EncodeBase64: %v", err)
		}

		expectedInput := bytes.NewBufferString(record.value)
		expectedDecoder := base64.NewDecoder(record.encoding, expectedInput)
		expectedDecoded, err := ioutil.ReadAll(expectedDecoder)
		if err != nil {
			t.Fatalf("Unable to decode expected value: %v", err)
		}

		assertions.JSONEq(string(expectedDecoded), string(actualDecoded))
	}
}
