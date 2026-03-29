package hub

import (
	"board/internal/canvas"
	"board/internal/color"
	"board/internal/cursor"
	"board/internal/name"
	"board/internal/ratelimit"
	"context"
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
)

type PaintMsg struct {
	Type  string `json:"type"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

type PaintBroadcast struct {
	Type  string      `json:"type"`
	X     int         `json:"x"`
	Y     int         `json:"y"`
	Color color.Color `json:"color"`
}

type CursorMsg struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type CursorBroadcast struct {
	Type   string      `json:"type"`
	UserID string      `json:"user_id"`
	Name   name.Name   `json:"name"`
	Color  color.Color `json:"color"`
	X      float64     `json:"x"`
	Y      float64     `json:"y"`
}

type CanvasStateMsg struct {
	ID     string          `json:"user_id"`
	Type   string          `json:"type"`
	Pixels [][]color.Color `json:"pixels"`
	Name   name.Name       `json:"name"`
	Color  color.Color     `json:"color"`
}

type Hub struct {
	node        *centrifuge.Node
	canvas      *canvas.Canvas
	limiter     *ratelimit.Limiter
	cursorStore *cursor.Store
}

func NewHub(canvas *canvas.Canvas, limiter *ratelimit.Limiter) (*Hub, error) {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create centrifuge node: %w", err)
	}
	return &Hub{
		node:        node,
		canvas:      canvas,
		limiter:     limiter,
		cursorStore: cursor.NewStore(),
	}, nil
}

func (h *Hub) Node() *centrifuge.Node {
	return h.node
}

func (h *Hub) Run() error {
	h.node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID: uuid.New().String(),
			},
		}, nil
	})

	h.node.OnConnect(func(client *centrifuge.Client) {
		client.OnSubscribe(func(event centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			snapshot := h.canvas.Snapshot()
			msg := CanvasStateMsg{ID: client.UserID(), Type: "canvas_state", Pixels: snapshot, Name: name.Random(), Color: color.Random()}

			h.cursorStore.Set(client.UserID(), cursor.NewCursor(0, 0, msg.Name, msg.Color))

			data, err := json.Marshal(msg)
			if err != nil {
				callback(centrifuge.SubscribeReply{}, centrifuge.ErrorInternal)
				return
			}
			if err := client.Send(data); err != nil {
				callback(centrifuge.SubscribeReply{}, centrifuge.ErrorInternal)
				return
			}
			callback(centrifuge.SubscribeReply{}, nil)
		})

		client.OnPublish(func(event centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			var base struct {
				Type string `json:"type"`
			}

			if err := json.Unmarshal(event.Data, &base); err != nil {
				callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
				return
			}

			switch base.Type {
			case "pixel_paint":
				paintMessage := PaintMsg{}
				if err := json.Unmarshal(event.Data, &paintMessage); err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorBadRequest)
					return
				}
				ok := h.limiter.TryPaint(client.UserID())
				if !ok {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorTooManyRequests)
					return
				}
				if err := h.canvas.Paint(paintMessage.X, paintMessage.Y, color.Color(paintMessage.Color)); err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}
				broadcast := PaintBroadcast{
					Type:  "pixel_paint",
					X:     paintMessage.X,
					Y:     paintMessage.Y,
					Color: color.Color(paintMessage.Color),
				}

				data, err := json.Marshal(broadcast)
				if err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}

				if _, err := h.node.Publish("canvas:main", data); err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}

			case "cursor_move":
				cursorMessage := CursorMsg{}
				if err := json.Unmarshal(event.Data, &cursorMessage); err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}
				h.cursorStore.Update(client.UserID(), cursorMessage.X, cursorMessage.Y)
				cursor, err := h.cursorStore.Get(client.UserID())
				if err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}

				broadcast := CursorBroadcast{
					Type:   "cursor_move",
					UserID: client.UserID(),
					Name:   cursor.Name,
					Color:  cursor.Color,
					X:      cursor.X,
					Y:      cursor.Y,
				}

				data, err := json.Marshal(broadcast)
				if err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}
				if _, err := h.node.Publish("canvas:main", data); err != nil {
					callback(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
					return
				}
			}

			callback(centrifuge.PublishReply{Result: &centrifuge.PublishResult{}}, nil)
		})
	})

	return h.node.Run()
}
