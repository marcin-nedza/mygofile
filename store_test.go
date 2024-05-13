package main

import (
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "test"
	pathkey := CASPathtransformFunc(key)
	fmt.Println(pathkey)
	expectedPath := "a94a8/fe5cc/b19ba/61c4c/0873d/391e9/87982/fbbd3"
	expectedFilename := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	if pathkey.Filename !=expectedFilename{
		t.Errorf("got [%s] expected [%s]",pathkey.Filename,expectedFilename)
	}
	if pathkey.Pathname !=expectedPath{
		t.Errorf("got [%s] expected [%s]",pathkey.Pathname,expectedPath)
	}
}
