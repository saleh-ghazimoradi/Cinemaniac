package service

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
)

func background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slg.Logger.Error(fmt.Sprintf("%v", err))
			}
		}()
		fn()
	}()
}
