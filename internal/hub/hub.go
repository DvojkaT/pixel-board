package hub

import (
	"board/internal/canvas"
	"board/internal/color"
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

type CanvasStateMsg struct {
	Type   string          `json:"type"`
	Pixels [][]color.Color `json:"pixels"`
}

type Hub struct {
	node    *centrifuge.Node
	canvas  *canvas.Canvas
	limiter *ratelimit.Limiter
}

func NewHub(canvas *canvas.Canvas, limiter *ratelimit.Limiter) (*Hub, error) {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create centrifuge node: %w", err)
	}
	return &Hub{
		node:    node,
		canvas:  canvas,
		limiter: limiter,
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
			msg := CanvasStateMsg{Type: "canvas_state", Pixels: snapshot}
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
			callback(centrifuge.PublishReply{}, nil)
		})
	})

	return h.node.Run()
}
