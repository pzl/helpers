/*
2-step authentication code generator, given the time-based keys.
Functions similarly to Google Authenticator

Copyright (c) 2014, 2018 Dan Panzarella

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type OTP interface {
	Hash(string) string
}

type TOTP struct {
	interval int
}

func (t *TOTP) Hash(secret string) string {
	timecode := int(time.Now().Unix()) / 30
	return hash(timecode, secret)
}

type HOTP struct {
	name  string
	count int
}

func (h *HOTP) Hash(secret string) string {
	path := stepFilePath()

	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error incrementing HOTP counter for", h.name)
	} else {
		lines := strings.Split(string(f), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, h.name+"=") {
				lines[i] = strings.Replace(line, ":"+strconv.Itoa(h.count), ":"+strconv.Itoa(h.count+1), 1)
				output := strings.Join(lines, "\n")
				err = ioutil.WriteFile(path, []byte(output), 0600)
			}
		}
	}

	return hash(h.count, secret)
}

func hmacsha(i int, secret string) []byte {
	data, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		panic(err)
	}
	h := hmac.New(sha1.New, data)
	message := make([]byte, 8)
	binary.BigEndian.PutUint64(message, uint64(i))
	h.Write(message)
	return h.Sum(nil)
}

func hash(i int, secret string) string {
	digest := hmacsha(i, secret)
	offset := int(digest[19]) & 0x0F
	code := (int(digest[offset])&0x7F)<<24 |
		(int(digest[offset+1])&0xFF)<<16 |
		(int(digest[offset+2])&0xFF)<<8 |
		(int(digest[offset+3]) & 0xFF)

	final := int(code) % 1000000 // digit control. 6 zeros = 6 digits

	return fmt.Sprintf("%d", final)
}

func userHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func stepFilePath() string {
	return filepath.Join(userHome(), ".config", ".2steps")
}

type Provider struct {
	name   string
	secret string
	hash   OTP
}

func (p *Provider) Hash() string {
	return p.hash.Hash(p.secret)
}

func config() map[string]Provider {
	path := stepFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("secrets file does not exist", path)
		os.Exit(1)
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "err opening file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	cfg := map[string]Provider{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.SplitN(scanner.Text(), "=", 2)
		name := s[0]
		secret := s[1]
		var hasher OTP
		if strings.ContainsRune(secret, ':') {
			resplit := strings.SplitN(secret, ":", 2)
			secret = resplit[0]
			count, _ := strconv.Atoi(resplit[1])
			hasher = &HOTP{name: name, count: count}
		} else {
			hasher = &TOTP{}
		}
		cfg[s[0]] = Provider{name: name, secret: secret, hash: hasher}
	}

	return cfg
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing argument: provider")
		os.Exit(1)
	}

	cfg := config()

	if val, exists := cfg[os.Args[1]]; exists {
		fmt.Println((&val).Hash())
	} else {
		fmt.Fprintln(os.Stderr, "provider "+os.Args[1]+" not found")
		os.Exit(1)
	}

}
