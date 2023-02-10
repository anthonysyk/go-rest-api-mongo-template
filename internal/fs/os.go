package fs

import (
	"fmt"
	"log"
	"os"
)

var _ Client = OSClient{}

type OSClient struct {
	root string
}

func NewOSClient(root string) (*OSClient, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Printf("folder %s do not exists, creating it ...", root)
		err := os.MkdirAll(root, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &OSClient{root: root}, nil
}

func (lc OSClient) getFullPath(id string) string {
	return fmt.Sprintf("%s/%s", lc.root, id)
}

func (lc OSClient) Write(id string, data []byte) error {
	err := os.WriteFile(lc.getFullPath(id), data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (lc OSClient) Read(id string) ([]byte, error) {
	content, err := os.ReadFile(lc.getFullPath(id))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (lc OSClient) Remove(id string) error {
	return os.Remove(lc.getFullPath(id))
}
