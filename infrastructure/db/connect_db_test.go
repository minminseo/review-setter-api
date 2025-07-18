package db_test

import (
	"fmt"
	"testing"
)

func TestNewDB(t *testing.T) {
	testcases := []struct {
		name string
	}{{name: ""}, {name: ""}, {name: ""}, {name: ""}}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			fmt.Println("write your unit test!")
		})
	}
}
func TestGetQuery(t *testing.T) {
	testcases := []struct {
		name string
	}{{name: ""}, {name: ""}}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			fmt.Println("write your unit test!")
		})
	}
}
func TestWithQueries(t *testing.T) {
	testcases := []struct {
		name string
	}{{name: ""}}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			fmt.Println("write your unit test!")
		})
	}
}
