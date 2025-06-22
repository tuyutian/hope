package utils

import (
	"log"
	"testing"
)

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Printf("current uuid:%v\n", Uuid())
	}
}

func TestMd5(t *testing.T) {
	log.Printf("123456 md5:%v\n", Md5("123456"))
	log.Println("e10adc3949ba59abbe56e057f20f883e" == Md5("123456"))
}

func TestRandInt64(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(RandInt64(1000, 9999))
	}
}
