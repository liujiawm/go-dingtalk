package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)



const oapi = "oapi.dingtalk.com"
const accessToken string = "d49b526f51985336ac96fa7219d7c11bf88f60532127b0c4a02a4dd7746f2896"

var webhookUrl = url.URL{
	Scheme: "https",
	Host:   oapi,
	Path:   "robot/send",
}

var secret = "SEC9f37a7277d57a57e25815e44f4b652997aedea390dcd840e2442203d5fb02768"
//const webhookUrl string = "https://oapi.dingtalk.com/robot/send?access_token=d49b526f51985336ac96fa7219d7c11bf88f60532127b0c4a02a4dd7746f2896"




func getUrl(timestamp,accessToken string, secret string)(string, error){
	dtu := webhookUrl
	value := url.Values{}
	value.Set("access_token", accessToken)

	// 如果不需要加签
	if secret == "" {
		dtu.RawQuery = value.Encode()
		return dtu.String(), nil
	}

	// 加签
	sign, err := sign(timestamp, secret)
	if err != nil {
		dtu.RawQuery = value.Encode()
		return dtu.String(), err
	}

	value.Set("timestamp", timestamp)
	value.Set("sign", sign)
	dtu.RawQuery = value.Encode()
	return dtu.String(), nil
}


func sign(timestamp,secret string)(string, error){
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	// return url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil))),nil
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func PostData(){
	var timestamp = strconv.FormatInt(time.Now().Unix()*1000, 10) // 该值要求毫秒


	postData := make(map[string]interface{})
	postData["msgtype"] = "text"
	postData["text"] = map[string]string{"content":"我就是我, 是不一样的烟火"}
	bytesData, _ := json.Marshal(postData)

	posturl,err := getUrl(timestamp,accessToken,secret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp, _ := http.Post(posturl, "application/json",bytes.NewReader(bytesData))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
