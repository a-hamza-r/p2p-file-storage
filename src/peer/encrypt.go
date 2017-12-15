package main

import (
    "fmt"
	"golang.org/x/crypto/blowfish"
    "log"
)


func Encrypt(src_file []byte,key []byte,size int64)[]byte{
    cipher, err := blowfish.NewCipher(key)
    if err != nil {
        fmt.Println(err.Error())
    }
    if size%8 != 0{
        // log.Fatal("size must be power of 2 greater than 8")
       pad := make([]byte, 8-(len(src_file)%8))
       src_file = append(src_file,pad...)
       size = int64(len(src_file))
	}
    encrypt := make([]byte, size)
    var i int64
	for i = 0; i < size; i+=8 {
	    cipher.Encrypt(encrypt[i:i+8], src_file[i:i+8])
	}
    return encrypt
}

func Decrypt(enc_file []byte,key []byte) []byte{
    cipher, err := blowfish.NewCipher(key)
    if err != nil {
        fmt.Println(err.Error())
    }
	var size int64 = int64(len(enc_file))
	if size%8 != 0{
		log.Fatal("size must be power of 2 greater than 8")
	}
    var decrypt [8]byte
    dec_file := make([]byte, size)
    var i int64
	for i = 0; i < size; i+=8 {
	    cipher.Decrypt(decrypt[0:], enc_file[i:])
	    copy(dec_file[i:i+8],decrypt[0:8])
	}
    return dec_file
}

