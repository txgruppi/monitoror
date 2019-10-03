//+build faker

package usecase

import (
	"math/rand"
	"time"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	httpUsecase struct {
	}
)

func NewHttpUsecase() http.Usecase {
	return &httpUsecase{}
}

// HttpAny only check status code
func (hu *httpUsecase) HttpAny(params *models.HttpAnyParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpAnyTileType, params.Url, params)
}

// HttpRaw check status code and content
func (hu *httpUsecase) HttpRaw(params *models.HttpRawParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpRawTileType, params.Url, params)
}

func (hu *httpUsecase) HttpJson(params *models.HttpJsonParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpJsonTileType, params.Url, params)
}

func (hu *httpUsecase) HttpYaml(params *models.HttpYamlParams) (tile *Tile, err error) {
	return hu.httpAll(http.HttpYamlTileType, params.Url, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in models.params
func (hu *httpUsecase) httpAll(tileType TileType, url string, params models.FakerParamsProvider) (tile *Tile, err error) {
	tile = NewTile(tileType)
	tile.Label = url

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.GetStatus(), randomStatus()).(TileStatus)
	if tile.Status == SuccessStatus && tileType != http.HttpAnyTileType {
		if len(params.GetValues()) != 0 {
			tile.Values = params.GetValues()
		} else if params.GetMessage() != "" {
			tile.Message = params.GetMessage()
		} else {
			if rand.Intn(2) == 0 {
				tile.Values = []float64{float64(rand.Intn(10000))}
			} else {
				tile.Message = "random message"
			}
		}

	}

	if tile.Status == FailedStatus {
		tile.Message = nonempty.String(params.GetMessage(), "random error message")
	}

	return
}

func randomStatus() TileStatus {
	if rand.Intn(2) == 0 {
		return SuccessStatus
	} else {
		return FailedStatus
	}
}
