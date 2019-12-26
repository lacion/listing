package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var (
	handlerLambda *httpadapter.HandlerAdapter
	newsletter    *NewsletterResource
	HtmlTemplate  *template.Template
	TextTemplate  *template.Template
)

const (
	contextSessionKey = "ctx_sess"
)

// Response is an alias of events.APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is an alias of events.APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is the main entry point to this lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handlerLambda.ProxyWithContext(ctx, req)
}

func main() {
	HtmlTemplate = template.Must(template.New("HtmlBody").Parse(HTMLBody))
	TextTemplate = template.Must(template.New("TextBody").Parse(TextBody))

	secret := os.Getenv("TOKEN_SECRET")
	apiToken := os.Getenv("API_TOKEN")
	subscribeRedirectUrl := os.Getenv("SUBSCRIBE_REDIRECT_URL")
	unsubscribeRedirectUrl := os.Getenv("UNSUBSCRIBE_REDIRECT_URL")
	confirmRedirectUrl := os.Getenv("CONFIRM_REDIRECT_URL")
	confirmUrl := os.Getenv("CONFIRM_URL")
	tableName := os.Getenv("SUBSCRIBERS_TABLE")
	supportedNewsletters := os.Getenv("SUPPORTED_NEWSLETTERS")
	emailFrom := os.Getenv("EMAIL_FROM")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		log.Fatalf("Failed to create AWS session. err=%v", err)
	}

	store := NewStore(tableName, sess)
	mailer := &SESMailer{
		svc:    ses.New(sess),
		sender: emailFrom,
		secret: secret,
	}

	router := http.NewServeMux()
	newsletter := &NewsletterResource{
		apiToken:               apiToken,
		secret:                 secret,
		subscribeRedirectURL:   subscribeRedirectUrl,
		unsubscribeRedirectURL: unsubscribeRedirectUrl,
		confirmRedirectURL:     confirmRedirectUrl,
		confirmURL:             confirmUrl,
		store:                  store,
		mailer:                 mailer,
		newsletters:            make(map[string]bool),
	}

	sn := strings.Split(supportedNewsletters, ";")
	newsletter.addNewsletters(sn)

	newsletter.setup(router)
	handlerLambda = httpadapter.New(router)

	lambda.Start(Handler)
}
