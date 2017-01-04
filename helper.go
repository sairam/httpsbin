package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString generates random alphabets
// source: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// RandStringBytesMaskImprSrc
func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func writeCompressedFileIO(file io.WriteCloser, data []byte) error {
	defer file.Close()

	fileWriter := gzip.NewWriter(file)
	fileWriter.Write(data)
	fileWriter.Close()

	return nil
}

func readCompressedFileIO(file io.ReadCloser) ([]byte, error) {
	defer file.Close()

	fileReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	data, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
