// +build mobile

package crocmobile

import (
	"fmt"
	"io/ioutil"

	"github.com/schollz/croc/v8/src/croc"
	"github.com/schollz/croc/v8/src/models"
	"github.com/schollz/croc/v8/src/utils"
)

type CallbackText interface {
	OnResponse(sharedSecret string)
}

type SendCallbacks interface {
	OnSharedSecret(sharedSecret string)
	OnCrocError(error error)
}

func makeTempFileWithString(tempPath string, s string) (fnames string, err error) {
	f, err := ioutil.TempFile(tempPath, "croc-stdin-")
	if err != nil {
		return
	}

	_, err = f.WriteString(s)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}
	fnames = f.Name()
	return
}

func Send(tempPath string, fileOrText string, isText bool, callbacks SendCallbacks) bool {
	secret := utils.GetRandomName()

	var file string
	if !isText {
		file = fileOrText
	} else {
		filex, err := makeTempFileWithString(tempPath, fileOrText)
		if err != nil {
			callbacks.OnCrocError(err)
			return false
		}
		file = filex
	}

	cr, err := croc.New(croc.Options{
		IsSender:       true,
		SharedSecret:   secret,
		Debug:          true,
		RelayAddress:   models.DEFAULT_RELAY,
		RelayAddress6:  models.DEFAULT_RELAY6,
		RelayPorts:     []string{"9009", "9010", "9011", "9012", "9013"},
		RelayPassword:  models.DEFAULT_PASSPHRASE,
		Stdout:         true,
		NoPrompt:       true,
		NoMultiplexing: true,
		DisableLocal:   true,
		Ask:            false,
		SendingText:    isText,
		NoCompress:     false,
	})

	if err != nil {
		fmt.Println(err)
		callbacks.OnCrocError(err)
		return false
	}

	callbacks.OnSharedSecret(secret)

	err = cr.Send(croc.TransferOptions{
		PathToFiles:      []string{file},
		KeepPathInRemote: false,
	})

	if err != nil {
		fmt.Println(err)
		callbacks.OnCrocError(err)
		return false
	}

	return true
}

// TODO
func Receive(sharedSecret string) {
	cr, err := croc.New(croc.Options{
		SharedSecret:  sharedSecret,
		IsSender:      false,
		Debug:         true,
		NoPrompt:      true,
		RelayAddress:  models.DEFAULT_RELAY,
		RelayAddress6: models.DEFAULT_RELAY6,
		Stdout:        true,
		Ask:           false,
		RelayPassword: models.DEFAULT_PASSPHRASE,
	})


	if err != nil {
		fmt.Println(err)
		//callbacks.OnCrocError(err)
		return
	}

	cr.Receive()
}