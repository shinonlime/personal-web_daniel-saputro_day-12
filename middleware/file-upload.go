package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileNames, fileHandler, err := r.FormFile("upload-image")
		if err != nil {
			fmt.Println("message : " + err.Error())
			json.NewEncoder(w).Encode("Error retrieving the file")
			return
		}

		defer fileNames.Close()

		fmt.Printf("Upload File : %+v\n", fileHandler.Filename)

		tempFile, err := ioutil.TempFile("uploads", "image-*"+fileHandler.Filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}

		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(fileNames)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		filename := data[8:]

		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
