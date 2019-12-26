package main

import (
	"encoding/json"
	"testing"
)

const sesNotificationJson = `{
      "notificationType":"Complaint",
      "complaint":{
         "userAgent":"Comcast Feedback Loop (V0.01)",
         "complainedRecipients":[
            {
               "emailAddress":"recipient1@example.com"
            }
         ],
         "complaintFeedbackType":"abuse",
         "arrivalDate":"2009-12-03T04:24:21.000-05:00",
         "timestamp":"2012-05-25T14:59:38.623-07:00",
         "feedbackId":"000001378603177f-18c07c78-fa81-4a58-9dd1-fedc3cb8f49a-000000"
      },
      "mail":{
         "timestamp":"2012-05-25T14:59:38.623-07:00",
     "messageId":"000001378603177f-7a5433e7-8edb-42ae-af10-f0181f34d6ee-000000",
         "source":"email_1337983178623@amazon.com",
         "sourceArn": "arn:aws:sns:us-east-1:XXXXXXXXXXXX:ses-notifications",
         "sendingAccountId":"XXXXXXXXXXXX",
         "destination":[
            "recipient1@example.com",
            "recipient2@example.com",
            "recipient3@example.com",
            "recipient4@example.com"
         ]
      }
   }`

func TestParseSesNotification(t *testing.T) {
	var sesMessage SesMessage
	err := json.Unmarshal([]byte(sesNotificationJson), &sesMessage)
	if err != nil {
		t.Fatal(err)
	}
}