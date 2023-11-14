package utils

import "log"

func IsError(err error, str1 string) {
	if err != nil {
		log.Printf("%s:%s\n", str1, err.Error())
	}
}

//func IsErrorFloat(err error, str1 string) {
//	if err != nil {
//		log.Fatalf("%s:%s\n", str1, err.Error())
//	}
//}
