package notifier

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Notifier struct {
	Proxy *url.URL
	token string
	param notifyRequestParams
}

type notifyRequestParams struct {
	message          string
	imageThumbnail   string
	imageFullsize    string
	imageFile        string
	stickerPackageId int
	stickerId        int
}

type NotifyResposnse struct {
	StatusCode int
	Header     map[string][]string
	Body       NotifyResponseBody
}

type NotifyResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (nf Notifier) NotifyMessage(token, message string) (NotifyResposnse, error) {

	nf.token = token
	nf.param.message = message

	return nf.sendNotification()

}

func (nf Notifier) NotifyImageFile(token, message, filePath string) (NotifyResposnse, error) {

	nf.token = token
	nf.param.message = message
	nf.param.imageFile = filePath

	return nf.sendNotification()

}

func (nf Notifier) NofityImageURL(token, message, URL, thumbnailURL string) (NotifyResposnse, error) {

	nf.token = token
	nf.param.message = message
	nf.param.imageFullsize = URL
	nf.param.imageThumbnail = thumbnailURL

	return nf.sendNotification()

}

func (nf Notifier) NotifySticker(token, message string, packageID, stickerID int) (NotifyResposnse, error) {

	nf.token = token
	nf.param.message = message
	nf.param.stickerPackageId = packageID
	nf.param.stickerId = stickerID

	return nf.sendNotification()

}

func (nf Notifier) sendNotification() (NotifyResposnse, error) {

	endpoint := "https://notify-api.line.me/api/notify"

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	if nf.param.message != "" {
		w.WriteField("message", nf.param.message)
	}

	if nf.param.imageFullsize != "" {
		w.WriteField("imageFullsize", nf.param.imageFullsize)
	}

	if nf.param.imageThumbnail != "" {
		w.WriteField("imageThumbnail", nf.param.imageThumbnail)
	}

	if nf.param.imageFile != "" {
		file, err := os.Open(nf.param.imageFile)
		if err != nil {
			return NotifyResposnse{}, err
		}

		fileWriter, err := w.CreateFormFile("imageFile", file.Name())
		io.Copy(fileWriter, file)
		if err != nil {
			return NotifyResposnse{}, err
		}
	}

	if nf.param.stickerPackageId != 0 {
		w.WriteField("stickerPackageId", strconv.Itoa(nf.param.stickerPackageId))
	}

	if nf.param.stickerId != 0 {
		w.WriteField("stickerId ", strconv.Itoa(nf.param.stickerId))
	}
	w.Close()

	req, err := http.NewRequest("POST", endpoint, b)
	if err != nil {
		return NotifyResposnse{}, err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+nf.token)

	http_client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(nf.Proxy),
		},
	}

	res, err := http_client.Do(req)
	if err != nil {
		return NotifyResposnse{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NotifyResposnse{}, err
	}

	ret := NotifyResposnse{}
	err = json.Unmarshal(body, &ret.Body)
	if err != nil {
		return NotifyResposnse{}, err
	}
	ret.Header = res.Header
	ret.StatusCode = res.StatusCode

	return ret, nil

}

func GetFileContentType(ouput *os.File) (string, error) {

	// to sniff the content type only the first
	// 512 bytes are used.

	buf := make([]byte, 512)

	_, err := ouput.Read(buf)

	if err != nil {
		return "", err
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}
