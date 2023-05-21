package weshnet

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
	coreapi "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"go.uber.org/zap"

	ipfslog "berty.tech/go-ipfs-log"
	"berty.tech/go-ipfs-log/identityprovider"
	ipliface "berty.tech/go-ipfs-log/iface"
	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/basestore"
	"berty.tech/go-orbit-db/stores/operation"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/pushtypes"
	"berty.tech/weshnet/pkg/secretstore"
	"berty.tech/weshnet/pkg/tyber"
)

// FIXME: replace cache by a circular buffer to avoid an attack by RAM saturation
type MessageStore struct {
	basestore.BaseStore
	eventBus event.Bus
	emitters struct {
		groupMessage      event.Emitter
		groupCacheMessage event.Emitter
	}

	secretStore               secretstore.SecretStore
	currentDevicePublicKey    crypto.PubKey
	currentDevicePublicKeyRaw []byte
	group                     *protocoltypes.Group
	groupPublicKey            crypto.PubKey
	logger                    *zap.Logger

	deviceCaches   map[string]*groupCache
	muDeviceCaches sync.RWMutex

	messagesQueue *simpleMessageQueue

	ctx    context.Context
	cancel context.CancelFunc
}

func (m *MessageStore) setLogger(l *zap.Logger) {
	if l == nil {
		return
	}

	m.logger = l.With(logutil.PrivateString("group-id", fmt.Sprintf("%.6s", base64.StdEncoding.EncodeToString(m.group.PublicKey))))
}

func (m *MessageStore) openMessage(ctx context.Context, e ipfslog.Entry) (*protocoltypes.GroupMessageEvent, error) {
	if e == nil {
		return nil, errcode.ErrInvalidInput
	}

	op, err := operation.ParseOperation(e)
	if err != nil {
		m.logger.Error("unable to parse operation", zap.Error(err))
		return nil, err
	}

	env, headers, err := m.secretStore.OpenEnvelopeHeaders(op.GetValue(), m.group)
	if err != nil {
		return nil, errcode.ErrCryptoDecrypt.Wrap(err)
	}

	devicePublicKey, err := crypto.UnmarshalEd25519PublicKey(headers.DevicePK)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	if !m.secretStore.IsChainKeyKnownForDevice(ctx, m.groupPublicKey, devicePublicKey) {
		if err := m.addToMessageQueue(ctx, e); err != nil {
			m.logger.Error("unable to add message to cache", zap.Error(err))
		}

		return nil, fmt.Errorf("no secret for device")
	}

	return m.processMessage(ctx, &messageItem{
		op:      op,
		env:     env,
		headers: headers,
		hash:    e.GetHash(),
	})
}

type groupCache struct {
	self, hasKnownChainKey bool
	locker                 sync.Locker
	queue                  *priorityMessageQueue
}

func (m *MessageStore) CacheSizeForDevicePK(devicePK []byte) (size int, ok bool) {
	m.muDeviceCaches.RLock()
	var device *groupCache
	if device, ok = m.deviceCaches[string(devicePK)]; ok {
		size = device.queue.Size()
	}
	m.muDeviceCaches.RUnlock()
	return
}

func (m *MessageStore) ProcessMessageQueueForDevicePK(ctx context.Context, devicePK []byte) {
	m.muDeviceCaches.Lock()
	if device, ok := m.deviceCaches[string(devicePK)]; ok {
		devicePublicKey, errDevice := crypto.UnmarshalEd25519PublicKey(devicePK)

		if errDevice != nil {
			m.logger.Error("unable to process message, unmarshal of device pk failed", logutil.PrivateBinary("devicepk", devicePK))
		} else if device.hasKnownChainKey = m.secretStore.IsChainKeyKnownForDevice(ctx, m.groupPublicKey, devicePublicKey); !device.hasKnownChainKey {
			m.logger.Error("unable to process message, no secret found for device pk", logutil.PrivateBinary("devicepk", devicePK))
		} else if next := device.queue.Next(); next != nil {
			m.messagesQueue.Add(next)
		}
	}
	m.muDeviceCaches.Unlock()
}

func (m *MessageStore) processMessage(ctx context.Context, message *messageItem) (*protocoltypes.GroupMessageEvent, error) {
	// process message
	msg, err := m.secretStore.OpenEnvelopePayload(ctx, message.env, message.headers, m.groupPublicKey, m.currentDevicePublicKey, message.hash)
	if err != nil {
		return nil, err
	}

	err = m.secretStore.UpdateOutOfStoreGroupReferences(ctx, message.headers.DevicePK, message.headers.Counter, m.group)
	if err != nil {
		m.logger.Error("unable to update push group references", zap.Error(err))
	}

	entry := message.op.GetEntry()
	eventContext := newEventContext(entry.GetHash(), entry.GetNext(), m.group)
	return &protocoltypes.GroupMessageEvent{
		EventContext: eventContext,
		Headers:      message.headers,
		Message:      msg.GetPlaintext(),
	}, nil
}

func (m *MessageStore) processMessageLoop(ctx context.Context, tracer *messageMetricsTracer) {
	for {
		// wait for next message
		message, ok := m.messagesQueue.WaitForNewItem(ctx)
		if !ok {
			// context expired, return
			return
		}

		devicePublicKeyString := string(message.headers.DevicePK)

		// check if we have the material to decrypt the message
		m.muDeviceCaches.Lock()
		device, ok := m.deviceCaches[devicePublicKeyString]
		if !ok {
			devicePublicKey, err := crypto.UnmarshalEd25519PublicKey(message.headers.DevicePK)
			if err != nil {
				m.logger.Error("unable to process message, unmarshal of device pk failed", logutil.PrivateBinary("devicepk", message.headers.DevicePK))
				continue
			}

			hasSecret := m.secretStore.IsChainKeyKnownForDevice(ctx, m.groupPublicKey, devicePublicKey)
			device = &groupCache{
				self:             bytes.Equal(m.currentDevicePublicKeyRaw, message.headers.DevicePK),
				queue:            newPriorityMessageQueue("undecrypted", tracer),
				locker:           &sync.RWMutex{},
				hasKnownChainKey: hasSecret,
			}
			m.deviceCaches[devicePublicKeyString] = device
		}

		// check for device chain key, if unavailable add message to cache queue
		if !device.hasKnownChainKey {
			device.queue.Add(message)
			_ = m.emitters.groupCacheMessage.Emit(*message)

			m.muDeviceCaches.Unlock()
			continue
		}

		// process the message
		evt, err := m.processMessage(ctx, message)
		if err != nil {
			if errcode.Is(err, errcode.ErrCryptoDecryptPayload) {
				// @FIXME(gfanton): this should not happen
				m.logger.Warn("unable to open envelope, adding envelope to cache for later process", zap.Error(err))
			} else {
				m.logger.Error("unable to process message", zap.Error(err))
			}

			// add to device queue, will try to load it on next received message
			device.queue.Add(message)
			_ = m.emitters.groupCacheMessage.Emit(*message)

			// if failed to process message, for later process
			m.muDeviceCaches.Unlock()
			continue
		}

		// emit new message event
		if err := m.emitters.groupMessage.Emit(*evt); err != nil {
			m.logger.Warn("unable to emit group message event", zap.Error(err))
		}

		m.muDeviceCaches.Unlock()
	}
}

func (m *MessageStore) addToMessageQueue(ctx context.Context, e ipfslog.Entry) error {
	if e == nil {
		return errcode.ErrInvalidInput
	}

	op, err := operation.ParseOperation(e)
	if err != nil {
		return err
	}

	env, headers, err := m.secretStore.OpenEnvelopeHeaders(op.GetValue(), m.group)
	if err != nil {
		return errcode.ErrCryptoDecrypt.Wrap(err)
	}

	msg := &messageItem{
		hash:    e.GetHash(),
		env:     env,
		headers: headers,
		op:      op,
	}

	m.messagesQueue.Add(msg)

	return nil
}

// FIXME: use iterator instead to reduce resource usage (require go-ipfs-log improvements)
func (m *MessageStore) ListEvents(ctx context.Context, since, until []byte, reverse bool) (<-chan *protocoltypes.GroupMessageEvent, error) {
	entries, err := getEntriesInRange(m.OpLog().GetEntries().Reverse().Slice(), since, until)
	if err != nil {
		return nil, err
	}

	out := make(chan *protocoltypes.GroupMessageEvent)

	go func() {
		iterateOverEntries(
			entries,
			reverse,
			func(entry ipliface.IPFSLogEntry) {
				message, err := m.openMessage(ctx, entry)
				if err != nil {
					m.logger.Error("unable to open message", zap.Error(err))
				} else {
					out <- message
					m.logger.Info("message store - sent 1 event from log history")
				}
			},
		)

		close(out)
	}()

	return out, nil
}

func (m *MessageStore) AddMessage(ctx context.Context, payload []byte) (operation.Operation, error) {
	ctx, newTrace := tyber.ContextWithTraceID(ctx)

	if newTrace {
		m.logger.Debug("Sending message to group "+base64.RawURLEncoding.EncodeToString(m.group.PublicKey), tyber.FormatTraceLogFields(ctx)...)
	}

	m.logger.Debug(
		fmt.Sprintf("Adding message to store with payload of %d bytes", len(payload)),
		tyber.FormatStepLogFields(
			ctx,
			[]tyber.Detail{
				{Name: "Payload", Description: string(payload)},
			},
		)...,
	)

	return messageStoreAddMessage(ctx, m.group, m, payload)
}

func messageStoreAddMessage(ctx context.Context, g *protocoltypes.Group, m *MessageStore, payload []byte) (operation.Operation, error) {
	msg, err := (&protocoltypes.EncryptedMessage{
		Plaintext:        payload,
		ProtocolMetadata: &protocoltypes.ProtocolMetadata{},
	}).Marshal()
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	sealedEnvelope, err := m.secretStore.SealEnvelope(ctx, g, msg)
	if err != nil {
		return nil, errcode.ErrCryptoEncrypt.Wrap(err)
	}
	m.logger.Debug(
		"Message sealed successfully in secretbox envelope",
		tyber.FormatStepLogFields(
			ctx,
			[]tyber.Detail{
				{Name: "Cleartext size", Description: fmt.Sprintf("%d bytes", len(msg))},
				{Name: "Cyphertext size", Description: fmt.Sprintf("%d bytes", len(sealedEnvelope))},
			},
		)...,
	)

	op := operation.NewOperation(nil, "ADD", sealedEnvelope)

	e, err := m.AddOperation(ctx, op, nil)
	if err != nil {
		return nil, errcode.ErrOrbitDBAppend.Wrap(err)
	}
	m.logger.Debug(
		"Envelope added to orbit-DB log successfully",
		tyber.FormatStepLogFields(ctx, []tyber.Detail{})...,
	)

	op, err = operation.ParseOperation(e)
	if err != nil {
		return nil, errcode.ErrOrbitDBDeserialization.Wrap(err)
	}

	m.logger.Debug(
		"Operation parsed by orbit-DB successfully",
		tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "CID", Description: op.GetEntry().GetHash().String()}})...,
	)

	return op, nil
}

func constructorFactoryGroupMessage(s *WeshOrbitDB, logger *zap.Logger) iface.StoreConstructor {
	metricsTracer := NewMessageMetricsTracer(s.prometheusRegister)
	return func(ipfs coreapi.CoreAPI, identity *identityprovider.Identity, addr address.Address, options *iface.NewStoreOptions) (iface.Store, error) {
		g, err := s.getGroupFromOptions(options)
		if err != nil {
			return nil, errcode.ErrInvalidInput.Wrap(err)
		}

		groupPublicKey, err := g.GetPubKey()
		if err != nil {
			return nil, errcode.ErrDeserialization.Wrap(err)
		}

		if options.EventBus == nil {
			options.EventBus = s.EventBus()
		}

		replication := false

		store := &MessageStore{
			// cmessage:       make(chan *messageItem),
			eventBus:       options.EventBus,
			secretStore:    s.secretStore,
			messagesQueue:  newMessageQueue("cache", metricsTracer),
			group:          g,
			groupPublicKey: groupPublicKey,
			logger:         logger,
			deviceCaches:   make(map[string]*groupCache),
		}

		if s.replicationMode {
			replication = true
		} else {
			currentMemberDevice, err := s.secretStore.GetOwnMemberDeviceForGroup(g)

			if err != nil {
				if errcode.Is(err, errcode.ErrInvalidInput) {
					replication = true
				} else {
					return nil, errcode.ErrOrbitDBInit.Wrap(err)
				}
			} else {
				store.currentDevicePublicKey = currentMemberDevice.Device()
				store.currentDevicePublicKeyRaw, err = store.currentDevicePublicKey.Raw()
				if err != nil {
					return nil, errcode.ErrOrbitDBInit.Wrap(err)
				}
			}
		}

		store.ctx, store.cancel = context.WithCancel(context.Background())

		go func() {
			store.processMessageLoop(store.ctx, metricsTracer)
			logger.Debug("store message process loop ended", zap.Error(store.ctx.Err()))
		}()

		if store.emitters.groupMessage, err = store.eventBus.Emitter(new(protocoltypes.GroupMessageEvent)); err != nil {
			store.cancel()
			return nil, errcode.ErrOrbitDBInit.Wrap(err)
		}

		// for debug/test purpose
		if store.emitters.groupCacheMessage, err = store.eventBus.Emitter(new(messageItem)); err != nil {
			store.cancel()
			return nil, errcode.ErrOrbitDBInit.Wrap(err)
		}

		options.Index = basestore.NewNoopIndex

		if err := store.InitBaseStore(ipfs, identity, addr, options); err != nil {
			store.cancel()
			return nil, errcode.ErrOrbitDBInit.Wrap(err)
		}

		if replication {
			return store, nil
		}

		chSub, err := store.EventBus().Subscribe([]interface{}{
			new(stores.EventWrite),
			new(stores.EventReplicated),
		}, eventbus.Name("weshnet/store-message"), eventbus.BufSize(128))
		if err != nil {
			return nil, fmt.Errorf("unable to subscribe to store events")
		}

		go func(ctx context.Context) {
			defer chSub.Close()
			for {
				var e interface{}
				select {
				case e = <-chSub.Out():
				case <-ctx.Done():
					return
				}

				var entries []ipfslog.Entry

				switch evt := e.(type) {
				case stores.EventWrite:
					entries = []ipfslog.Entry{evt.Entry}

				case stores.EventReplicated:
					entries = evt.Entries
				}

				for _, entry := range entries {
					ctx = tyber.ContextWithConstantTraceID(ctx, "msgrcvd-"+entry.GetHash().String())
					store.logger.Debug("Received message store event", tyber.FormatTraceLogFields(ctx)...)
					store.logger.Debug(
						"Message store event",
						tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "RawEvent", Description: fmt.Sprint(e)}})...,
					)

					if err := store.addToMessageQueue(ctx, entry); err != nil {
						logger.Error("unable to add message to queue", zap.Error(err))
					}
				}
			}
		}(store.ctx)

		return store, nil
	}
}

func (m *MessageStore) GetMessageByCID(c cid.Cid) (operation.Operation, error) {
	logEntry, ok := m.OpLog().Get(c)
	if !ok {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("unable to find message entry"))
	}

	op, err := operation.ParseOperation(logEntry)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	return op, nil
}

func (m *MessageStore) GetOutOfStoreMessageEnvelope(ctx context.Context, c cid.Cid) (*pushtypes.OutOfStoreMessageEnvelope, error) {
	op, err := m.GetMessageByCID(c)
	if err != nil {
		return nil, errcode.ErrInvalidInput.Wrap(err)
	}

	env, headers, err := m.secretStore.OpenEnvelopeHeaders(op.GetValue(), m.group)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	sealedMessageEnvelope, err := m.secretStore.SealOutOfStoreMessageEnvelope(c, env, headers, m.group)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	return sealedMessageEnvelope, nil
}

func (m *MessageStore) Close() error {
	m.cancel()
	return m.BaseStore.Close()
}
