package installbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
/* func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	name := context.GetInput("name").(string)
	version := context.GetInput("version").(string)
	url := context.GetInput("url").(string)
	msKey := fmt.Sprintf("%v%v", name, version)
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("microservicemanager.db")
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		_, err = tx.Get(msKey)
		return err
	})

	if err == nil {
		return false, errors.New("Microservice is already installed")
	}

	DownloadToFile(url, "./"+msKey)

	// add microservice
	microService := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		URL     string `json:"URL"`
		Status  string `json:"status"`
	}{
		Name:    name,
		Version: version,
		URL:     url,
		Status:  "installed",
	}
	jsonMicroService, err := json.Marshal(microService)
	if err != nil {
		return false, err
	}

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(msKey, string(jsonMicroService), nil)
		return err
	})
	if err != nil {
		return false, err
	}

	return true, nil
} */

func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	name := context.GetInput("name").(string)
	version := context.GetInput("version").(string)
	url := context.GetInput("url").(string)
	msKey := fmt.Sprintf("%v%v", name, version)

	// download bin
	DownloadToFile(url, "/home/pi/gateway/"+msKey)
	err = RunCMD("chmod +x /home/pi/gateway/" + msKey)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	// download conf
	DownloadToFile(url+".d", "/etc/systemd/system/"+msKey+".service")
	// daemon reload
	err = RunCMD("systemctl daemon-reload")
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	// service start
	err = RunCMD("systemctl restart" + msKey)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	// service enabled
	err = RunCMD("systemctl enable" + msKey)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	// save installed message
	microService := make(map[string]interface{})
	microService["key"] = msKey
	microService["name"] = name
	microService["version"] = version
	microService["url"] = url
	microService["enabled"] = true

	err = SaveInstalledService(microService)
	if err != nil {
		return false, err
	}
	return true, nil
}

func HTTPDownload(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))

	err := ioutil.WriteFile(dst, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func DownloadToFile(uri string, dst string) {
	fmt.Printf("DownloadToFile From: %s.\n", uri)
	if d, err := HTTPDownload(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)
		RunCMD(fmt.Sprintf("chmod +x %v ", dst))
		if WriteFile(dst, d) == nil {
			fmt.Printf("saved %s as %s\n", uri, dst)
		}
	}
}

func RunCMD(cmd string) error {
	command := exec.Command(cmd)

	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	command.Stdout = &out
	err := command.Run()
	return err
}

func SaveInstalledService(microService map[string]interface{}) error {
	db, err := buntdb.Open("/home/pi/gateway/gateway.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *buntdb.Tx) error {
		json, _ := json.Marshal(microService)
		_, _, err := tx.Set(microService["key"].(string), string(json), nil)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
