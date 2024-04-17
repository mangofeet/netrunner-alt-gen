package art

import (
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type Drawer interface {
	Draw(ctx *canvas.Context, card *nrdb.Printing) error
}
