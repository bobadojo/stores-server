package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bobadojo/go/pkg/stores/v1/storespb"
	"github.com/tidwall/rtree"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	var stores []*storespb.Store
	c := 0
	startTime := time.Now()
	for i := 0; i < 3000; i++ {
		filename := "data/us_post_offices_01-27-23-sample.csv"
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		stores = nil
		reader := csv.NewReader(file)
		reader.Read()
		for {
			fields, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			store := &storespb.Store{
				Name:  fields[0],
				Type:  fields[1],
				Title: fields[2],
				Location: &storespb.Location{
					Latitude:  atof(fields[10]),
					Longitude: atof(fields[11]),
				},
				Address: &storespb.Address{
					Street:     fields[3],
					City:       fields[4],
					State:      fields[5],
					ZipCode:    atoi(fields[6]),
					RegionCode: fields[13],
					County:     fields[14],
				},
				StoreHours: fields[9],
			}
			stores = append(stores, store)
			c++
		}
	}
	executionTime := time.Since(startTime)
	fmt.Printf("load of %d records ran in %+v\n", c, executionTime)

	b, err := proto.Marshal(&storespb.ListStoresResponse{Stores: stores})
	if err != nil {
		return err
	}
	err = os.WriteFile("stores.pb", b, 0644)
	if err != nil {
		return err
	}

	var tr rtree.RTree
	for _, store := range stores {
		for i := 0; i < 100000; i++ {
			point := [2]float64{
				float64(store.Location.Latitude) + float64(i)/10000.0,
				float64(store.Location.Longitude) + float64(i)/10000.0,
			}
			tr.Insert(point, point, store.Name+fmt.Sprintf("_%d", i))
		}
	}

	log.Printf("%+v", tr.Len())

	point := [2]float64{
		38.08973633943899,
		-96.8860452622175,
	}
	min := point
	max := point
	epsilon := 0.001
	min[0] -= epsilon
	min[1] -= epsilon
	max[0] += epsilon
	max[1] += epsilon

	startTime = time.Now()
	tr.Search(min, max, func(min, max [2]float64, data interface{}) bool {
		log.Printf("%+v", data)
		return true
	})
	executionTime = time.Since(startTime)
	fmt.Printf("query ran in %+v\n", executionTime)

	return nil
}

func atoi(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return int32(i)
}

func atof(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return float32(f)
}
