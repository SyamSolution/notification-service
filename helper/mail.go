package helper

import (
	"bytes"
	"fmt"
	"github.com/SyamSolution/notification-service/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
	"os"
	"text/template"
)

func SendCreateTransactionMail(message model.DataMessage) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		fmt.Println("Error creating session", err)
	}

	svc := ses.New(sess)

	tmpl, err := template.New("email").Parse(emailTmplCreateTransaction)
	if err != nil {
		fmt.Println("Error parsing email template", err)
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, message); err != nil {
		fmt.Println("Error executing email template", err)
	}

	emailParams := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(message.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(htmlBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("New Transaction"),
			},
		},
		Source: aws.String("syams.arie@gmail.com"),
	}

	_, err = svc.SendEmail(emailParams)
	if err != nil {
		fmt.Println("Error sending email", err)
	}
}

func SendEmailWithPDF(message model.EmailPDFMessage) {
	pdfBase64, err := GeneratePDF(message)
	if err != nil {
		log.Println("Error generating PDF", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		fmt.Println("Error creating session", err)
	}
	svc := ses.New(sess)

	rawEmail := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"Content-Transfer-Encoding: 7bit\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: application/pdf; name=\"attachment.pdf\"\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"Content-Disposition: attachment; filename=\"attachment.pdf\"\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s--\r\n",
		"syams.arie@gmail.com", message.Email, "pdf", "boundary123", "boundary123", emailTmplCompleteTransaction, "boundary123", pdfBase64, "boundary123",
	)

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: []byte(rawEmail)},
	}

	_, err = svc.SendRawEmail(input)
	if err != nil {
		panic(err)
	}
}

var emailTmplCreateTransaction = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Template</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f7f7f7;
            margin: 0;
            padding: 0;
        }

        .container {
            width: 100%;
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        }

        .header {
            background-color: #007bff;
            color: #ffffff;
            padding: 20px;
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
        }

        .content {
            padding: 30px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #dddddd;
        }

        th {
            background-color: #f2f2f2;
        }

        td {
            background-color: #ffffff;
        }

        .footer {
            text-align: center;
            background-color: #007bff;
            color: #ffffff;
            padding: 15px;
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
        }

        .button {
            display: inline-block;
            background-color: #007bff;
            color: #ffffff;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 5px;
			text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Menunggu Pembayaran</h2>
        </div>
        <div class="content">
        	<p>Halo, <strong>{{.Name}}</strong></p>
            
            <p> Segera lakukan pembayaran pesananmu dengan detail sebagai berikut sebelum <strong> {{.DeadlineDate}} </strong>: </p>
            
            <table>
                <tr>
                    <td><strong>Order ID:</strong></td>
                    <td>{{.OrderID}}</td>
                </tr>
                <tr>
                    <td><strong>Date:</strong></td>
                    <td>{{.Date}}</td>
                </tr>
                <tr>
                    <td><strong>Total:</strong></td>
                    <td>{{.Total}}</td>
                </tr>
            </table>
            
            <br>
            
            <a href="{{.URL}}" class="button">Payment Here</a>
        </div>
        <div class="footer">
            <p>Copyright © 2024. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

var emailTmplCompleteTransaction = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Template</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f7f7f7;
            margin: 0;
            padding: 0;
        }

        .container {
            width: 100%;
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        }

        .header {
            background-color: #007bff;
            color: #ffffff;
            padding: 20px;
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
        }

        .content {
            padding: 30px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #dddddd;
        }

        th {
            background-color: #f2f2f2;
        }

        td {
            background-color: #ffffff;
        }

        .footer {
            text-align: center;
            background-color: #007bff;
            color: #ffffff;
            padding: 15px;
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
        }

        .button {
            display: inline-block;
            background-color: #007bff;
            color: #ffffff;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Pembayaran Berhasil</h2>
        </div>
        <div class="content">
            <p> Pembayaran anda telah berhasil</p>
        </div>
        <div class="footer">
            <p>Copyright © 2024. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
