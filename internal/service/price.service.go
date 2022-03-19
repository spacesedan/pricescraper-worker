package service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"pricescraper-worker/internal/dto"
	"pricescraper-worker/internal/models"
	"pricescraper-worker/internal/repo"
	"pricescraper-worker/internal/utils"
	"sync"
	"time"
)

type PriceService interface {
	CheckIfOnReservoir(slug string) bool
	GetReservoirPriceMap(contractAddress, slug string) *models.ReservoirPriceMap
	GetOpenseaCollectionFloor(collection models.Collection) models.OpenSeaStats
	ConsumerOpensea(ch <-chan dto.OpenseaTask)
	ConsumerReservoir(ch <-chan dto.ReservoirTask)
	HandleCollections(colls models.ChunkedCollection)
	UpdateMetaFloorPrice(collectionName string)
	UpdateTraitFloorPrices(collectionName string)
}

type priceService struct {
	dao repo.DAO
}

func NewPriceService(dao repo.DAO) PriceService {
	return &priceService{
		dao: dao,
	}
}

func (p *priceService) CheckIfOnReservoir(slug string) bool {
	uri := "https://api.reservoir.tools/collections/" + slug + "/v1"

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.Set("Content-Type", "application/json")
	req.SetRequestURI(uri)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		log.Fatal(err)
	}

	var reservoirData models.ReservoirCollectionResponse

	body := res.Body()

	err := json.Unmarshal(body, &reservoirData)
	if err != nil {
		log.Println(err)
	}

	if reservoirData.Collection.ID == "" {
		return false
	} else {
		return true
	}

}

func (p *priceService) GetReservoirPriceMap(contractAddress, slug string) *models.ReservoirPriceMap {
	uri := "https://api.reservoir.tools/tokens/floor/v1?contract=" + contractAddress
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.Set("Content-Type", "application/json")
	req.SetRequestURI(uri)

	res := fasthttp.AcquireResponse()

	if err := fasthttp.Do(req, res); err != nil {
		log.Fatal(err)
	}

	var reservoirRes *models.ReservoirPriceMap
	body := res.Body()

	fmt.Println(uri, slug)

	err := json.Unmarshal(body, &reservoirRes)
	if err != nil {
		fmt.Println(uri, slug)
		//log.Println(err)
	}

	return reservoirRes
}

func (p *priceService) GetOpenseaCollectionFloor(collection models.Collection) models.OpenSeaStats {
	uri := "https://api.opensea.io/api/v1/collection/" + collection.OSSlug + "/stats"

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.Set("Content-Type", "application/json")
	req.SetRequestURI(uri)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		log.Fatalln(err)
	}

	var stats models.OpenSeaStats
	body := res.Body()

	err := json.Unmarshal(body, &stats)
	if err != nil {
		log.Println("Poop", err)
	}

	return stats
}

func (p *priceService) HandleOpenSea(collections []dto.OpenseaTask) {
	in := make(chan dto.OpenseaTask)
	out := make(chan dto.OpenseaTask)

	for i := 0; i < len(collections); i++ {
		go utils.Worker(in, out)
	}

	go utils.Producer(in, collections)
	p.ConsumerOpensea(out)
}

func (p *priceService) HandleReservoir(tasks []dto.ReservoirTask) {
	in := make(chan dto.ReservoirTask)
	out := make(chan dto.ReservoirTask)

	for i := 0; i <= len(tasks); i++ {
		go utils.Worker(in, out)
	}
	go utils.Producer(in, tasks)
	p.ConsumerReservoir(out)

}

func (p *priceService) ConsumerOpensea(ch <-chan dto.OpenseaTask) {
	for {
		select {
		case msg := <-ch:
			p.dao.NewMetaQuery().UpdateFloorPrice(msg.Collection.Collection, msg.Stats.FloorPrice)
		case <-time.After(300 * time.Millisecond):
			return
		}
	}
}

func (p *priceService) ConsumerReservoir(ch <-chan dto.ReservoirTask) {
	for {
		select {
		case msg := <-ch:
			p.dao.NewCollectionQuery().RemoveOldPrices(msg.Collection.Collection)
			p.dao.NewCollectionQuery().StorePriceMap(msg.PriceMap, msg.Collection.Collection)
			p.UpdateMetaFloorPrice(msg.Collection.Collection)
			p.UpdateTraitFloorPrices(msg.Collection.Collection)
		case <-time.After(300 * time.Millisecond):
			return
		}
	}
}

func (p *priceService) UpdateMetaFloorPrice(collectionName string) {
	floorPrice := p.dao.NewCollectionQuery().FindLowestFloor(collectionName)
	p.dao.NewMetaQuery().UpdateFloorPrice(collectionName, floorPrice)
}

func (p *priceService) UpdateTraitFloorPrices(collectionName string) {
	traits, err := p.dao.NewTraitQuery().GetAllTraits(collectionName)
	if err != nil {
		log.Fatal(err)
	}
	length := len(traits)

	jobs := make(chan int, length)
	tokens := make(chan dto.CollectionWithTraitCombo, length)

	var wg sync.WaitGroup
	wg.Add(length)
	// use a worker pool to find all the tokens the token with the lowest price that has
	// that specific tait at the same time
	for i := 0; i < length; i++ {
		defer wg.Done()
		ttype := traits[i].ID.TraitType
		tvalue := traits[i].ID.Value
		go p.dao.NewCollectionQuery().FindWithTraitCombo(collectionName, ttype, tvalue, jobs, tokens)
	}

	for j := 0; j < length; j++ {
		// Send jobs to the job channel so that it begins on querying the db
		jobs <- j
	}
	go func() {
		wg.Wait()
		defer close(jobs)
	}()

	p.dao.NewTraitQuery().UpdateTraitPrices(collectionName, length, tokens)

}

func (p *priceService) HandleCollections(colls models.ChunkedCollection) {
	var onReservoir []dto.ReservoirTask
	var notOnReservoir []dto.OpenseaTask

	for _, collection := range colls.Collections {
		exists := p.CheckIfOnReservoir(collection.OSSlug)
		switch exists {
		case true:
			rPriceMap := p.GetReservoirPriceMap(collection.ContractAddress, collection.Collection)
			rTask := dto.ReservoirTask{
				PriceMap:   *rPriceMap,
				Collection: collection,
			}
			onReservoir = append(onReservoir, rTask)
		default:
			osRes := p.GetOpenseaCollectionFloor(collection)
			osTask := dto.OpenseaTask{
				Stats:      osRes.Stats,
				Collection: collection,
			}
			notOnReservoir = append(notOnReservoir, osTask)
		}

	}

	fmt.Println("length on reservoir", len(onReservoir))
	fmt.Println("length not on reservoir", len(notOnReservoir))
	p.HandleReservoir(onReservoir)
	p.HandleOpenSea(notOnReservoir)
}
