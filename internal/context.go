package internal

import (
	"context"
	"github.com/getsentry/sentry-go"

	"github.com/rs/zerolog"
)

type ctx string

var (
	ctxData ctx = "syncv3_data"
)

// logging metadata for a single request
type data struct {
	userID               string
	deviceID             string
	since                int64
	next                 int64
	numRooms             int
	txnID                string
	numToDeviceEvents    int
	numGlobalAccountData int
	numChangedDevices    int
	numLeftDevices       int
}

// prepare a request context so it can contain syncv3 info
func RequestContext(ctx context.Context) context.Context {
	d := &data{
		since:    -1,
		next:     -1,
		numRooms: -1,
	}
	return context.WithValue(ctx, ctxData, d)
}

// add the user ID to this request context. Need to have called RequestContext first.
func SetRequestContextUserID(ctx context.Context, userID, deviceID string) {
	d := ctx.Value(ctxData)
	if d == nil {
		return
	}
	da := d.(*data)
	da.userID = userID
	da.deviceID = deviceID
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{Username: userID})
		})
	}
}

func SetRequestContextResponseInfo(
	ctx context.Context, since, next int64, numRooms int, txnID string, numToDeviceEvents, numGlobalAccountData int,
	numChangedDevices, numLeftDevices int,
) {
	d := ctx.Value(ctxData)
	if d == nil {
		return
	}
	da := d.(*data)
	da.since = since
	da.next = next
	da.numRooms = numRooms
	da.txnID = txnID
	da.numToDeviceEvents = numToDeviceEvents
	da.numGlobalAccountData = numGlobalAccountData
	da.numChangedDevices = numChangedDevices
	da.numLeftDevices = numLeftDevices
}

func DecorateLogger(ctx context.Context, l *zerolog.Event) *zerolog.Event {
	d := ctx.Value(ctxData)
	if d == nil {
		return l
	}
	da := d.(*data)
	if da.userID != "" {
		l = l.Str("u", da.userID)
	}
	if da.deviceID != "" {
		l = l.Str("dev", da.deviceID)
	}
	if da.since >= 0 {
		l = l.Int64("p", da.since)
	}
	if da.next >= 0 {
		l = l.Int64("q", da.next)
	}
	if da.txnID != "" {
		l = l.Str("t", da.txnID)
	}
	if da.numRooms >= 0 {
		l = l.Int("r", da.numRooms)
	}
	if da.numToDeviceEvents > 0 {
		l = l.Int("d", da.numToDeviceEvents)
	}
	if da.numGlobalAccountData > 0 {
		l = l.Int("ag", da.numGlobalAccountData)
	}
	if da.numChangedDevices > 0 {
		l = l.Int("dl-c", da.numChangedDevices)
	}
	if da.numLeftDevices > 0 {
		l = l.Int("dl-l", da.numLeftDevices)
	}
	return l
}
