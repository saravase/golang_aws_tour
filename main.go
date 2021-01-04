package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type AWS struct {
	AccessId  string
	SecretKey string
	Region    string
	Bucket    string
}

func NewAWS() *AWS {
	return &AWS{
		AccessId:  GetEnvWithKey("AWS_ACCESS_KEY_ID"),
		SecretKey: GetEnvWithKey("AWS_SECRET_ACCESS_KEY"),
		Region:    GetEnvWithKey("AWS_REGION"),
		Bucket:    GetEnvWithKey("AWS_S3_BUCKET"),
	}
}

func main() {

	// Load environment variables
	LoadEnv()

	_aws := NewAWS()

	// Create AWS session
	sess := _aws.ConnectAWS()

	// Initialize appication
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("sess", sess)
		return c.Next()
	})

	app.Post("/upload", _aws.HandlerFileUpload)

	app.Listen(":9090")
}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	log.Info("Environment variables loaded successfully")
}

func (_aws *AWS) ConnectAWS() *session.Session {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(_aws.Region),
		Credentials: credentials.NewStaticCredentials(
			_aws.AccessId,
			_aws.SecretKey,
			"", // a token will be created when the session it's used.
		),
	})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("AWS session created successfully")
	return sess
}

func (_aws *AWS) HandlerFileUpload(c *fiber.Ctx) error {

	sess := c.Locals("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	file, err := c.FormFile("photo")
	filename := "avatar/" + file.Filename
	if err != nil {
		log.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Failed to read file",
			"filename": filename,
		})
	}

	data, err := file.Open()
	if err != nil {
		log.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to open file",
			"filename": filename,
		})
	}

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(_aws.Bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   data,
	})

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to upload file",
			"uploader": up,
		})
	}

	filepath := "https://" + _aws.Bucket + "." + "s3-" + _aws.Region + ".amazonaws.com/" + filename
	log.WithFields(log.Fields{
		"filepath": filepath,
	}).Info("Uploaded File Path")

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"filepath": filepath,
	})
}
