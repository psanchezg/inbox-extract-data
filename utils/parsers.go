package utils

import (
	"encoding/base64"
	"log"
	"regexp"
	"strings"

	"google.golang.org/api/gmail/v1"
)

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func GetParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	if paramsMap["Seg"] != "" && paramsMap["Min"] == "" {
		paramsMap["Min"] = "0"
	}
	return paramsMap
}

func GetMsgBody(msg *gmail.Message) (string, error) {
	// Obtener el cuerpo del mensaje
	var body string
	if msg.Payload.Body.Data != "" {
		bodyBytes, err := base64.StdEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			log.Fatalf("Error al decodificar el cuerpo del mensaje: %v", err)
			return "", err
		}
		body = strings.ReplaceAll(string(bodyBytes), "\r", "")
	} else {
		// Si el cuerpo está en multipart, iterar por las partes para encontrar el cuerpo deseado
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/plain" || part.MimeType == "text/html" {
				bodyBytes, err := base64.StdEncoding.DecodeString(part.Body.Data)
				if err != nil {
					log.Fatalf("Error al decodificar el cuerpo del mensaje: %v", err)
					return "", err
				}
				body = strings.ReplaceAll(string(bodyBytes), "\r", "")
				break
			}
		}
	}
	return body, nil
}
