package main

import (
	"log"
	"poke/models"
	"runtime"
	"snatch"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}
func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	fids, err := snatch.GetFid()
	if err != nil {
		log.Fatalln(err)
	}
	avfids, err := snatch.ChkFidAv(fids)
	if err != nil {
		log.Fatalln(err)
	}

	err = models.InsertFids(avfids)
	if err != nil {
		log.Fatalln(err)
	}

	models.DeleteAllTids()

	wg := sync.WaitGroup{}
	for _, v := range fids {
		go func(v string) {
			wg.Add(1)
			defer func() {
				wg.Done()
			}()

			tids, err := snatch.GetTid(v)
			if err != nil {
				log.Println(err)
				return
			}
			avtids, err := snatch.ChkTidAv(tids)
			if err != nil {
				log.Println(err)
				return
			}
			err = models.InsertTids(avtids)
			if err != nil {
				log.Println(err)
				return
			}
		}(v)

	}
	wg.Wait()

}
