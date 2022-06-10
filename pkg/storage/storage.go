package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

type IStorage interface {
	Upload(bucket string, object string, file io.Reader) (string, error)
	SignedURL(bucket string, object string) (string, error)
	GetFile(bucket string, fileName string) ([]byte, error)
	MakePublic(bucket string, object string) (string, error)
}

// Storage is a storage client
type Storage struct {
	client *storage.Client
	creds  *jwt.Config
}

var storageURL = "https://storage.googleapis.com"

// New Creates a new storage client
func New() (*Storage, error) {

	c, err := storage.NewClient(context.Background())
	if err != nil {
		return &Storage{}, err
	}

	creds, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(creds)
	if err != nil {
		return nil, err
	}

	return &Storage{client: c, creds: conf}, nil
}

// Upload uploads file to bucket and returns url to uploaded file
func (c *Storage) Upload(bucket string, object string, file io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	wc := c.client.Bucket(bucket).Object(object).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", errors.Wrap(err, "error io.Copy")
	}
	if err := wc.Close(); err != nil {
		return "", errors.Wrap(err, "writer.Close")
	}

	_, err := c.client.Bucket(bucket).Object(object).Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	})
	return object, err
}

// SignedURL signs url and returns signed url
func (c *Storage) SignedURL(bucket string, object string) (string, error) {

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: c.creds.Email,
		PrivateKey:     c.creds.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(time.Minute * 86400),
	})

	if err != nil {
		return "", errors.Wrap(err, "error signing url")
	}

	return url, nil
}

func (c *Storage) GetFile(bucket string, fileName string) ([]byte, error) {
	rc, err := c.client.Bucket(bucket).Object(fileName).NewReader(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "readFile: unable to open file from bucket")
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "readFile: unable to read file from bucket")
	}
	return slurp, err
}

// MakePublic makes an object public
func (c Storage) MakePublic(bucket string, object string) (string, error) {

	acl := c.client.Bucket(bucket).Object(object).ACL()

	err := acl.Set(context.Background(), storage.AllUsers, storage.RoleReader)
	if err != nil {
		return "", errors.Wrap(err, "error set acl")
	}
	return storageURL + "/" + bucket + "/" + object, nil

}
