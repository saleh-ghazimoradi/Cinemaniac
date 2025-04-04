package service

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"sync"
)

var WG sync.WaitGroup

func background(fn func()) {
	WG.Add(1)
	go func() {
		defer WG.Done()
		defer func() {
			if err := recover(); err != nil {
				slg.Logger.Error(fmt.Sprintf("%v", err))
			}
		}()
		fn()
	}()
}
