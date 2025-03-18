package service

import (
	"fmt"
	"github.com/qoentz/evedict/internal/api/dto"
	"github.com/qoentz/evedict/internal/eventfeed/polymarket"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"log"
	"strconv"
)

type MarketService struct {
	PolyMarketService *polymarket.Service
	AIService         llm.Service
}

func NewMarketService(polyMarketService *polymarket.Service, aiService llm.Service) *MarketService {
	return &MarketService{
		PolyMarketService: polyMarketService,
		AIService:         aiService,
	}
}

func (s *MarketService) GetMarketEvents(num int) ([]polymarket.Event, error) {
	events, err := s.PolyMarketService.FetchTopEvents()
	if err != nil {
		return nil, fmt.Errorf("error fetching events: %v", err)
	}

	var SMPEvents []polymarket.Event
	for _, e := range events {
		if len(e.Markets) == 1 {
			SMPEvents = append(SMPEvents, e)
		}
	}
	selectedIndexes, err := s.AIService.SelectIndexes(promptgen.SelectMarkets, struct {
		Events []polymarket.Event
	}{Events: SMPEvents}, num)
	if err != nil {
		return nil, fmt.Errorf("error selecting markets: %v", err)
	}

	var selectedMarkets []polymarket.Event
	for _, idx := range selectedIndexes {
		selectedMarkets = append(selectedMarkets, SMPEvents[idx])
	}

	return selectedMarkets, nil
}

func (s *MarketService) AttachMarketData(event polymarket.Event, forecast *dto.Forecast) {
	var firstMarket polymarket.Market
	if len(event.Markets) > 0 {
		firstMarket = event.Markets[0]
	}

	marketID, err := strconv.ParseInt(firstMarket.ID, 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid market ID %q, skipping market assignment", firstMarket.ID)
		return
	}

	forecast.Market = &dto.Market{
		ID:            marketID,
		Question:      firstMarket.Question,
		Outcomes:      firstMarket.Outcomes,
		OutcomePrices: firstMarket.OutcomePrices,
		Volume:        firstMarket.Volume,
		ImageURL:      event.Image,
	}
}
