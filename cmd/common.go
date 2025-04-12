package main

import "io/ioutil"

func writeToFile(src []byte, outputName string) error {
	return ioutil.WriteFile(outputName, src, 0600)
}
