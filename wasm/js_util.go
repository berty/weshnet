//go:build js

package main

import (
	"fmt"
	"syscall/js"

	"errors"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func consoleLog(args ...any) {
	js.Global().Get("console").Call("log", args...)
}

func awaitRaw(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc, catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}

func await(awaitable js.Value) (js.Value, error) {
	success, failure := awaitRaw(awaitable)
	if len(failure) == 0 {
		if len(success) == 0 {
			return js.Undefined(), nil
		}
		return success[0], nil
	}
	return js.Undefined(), errors.New(failure[0].Get("message").String())
}

func promisify(cb func() ([]any, error)) js.Value {
	// Handler for the Promise: this is a JS function
	// It receives two arguments, which are JS functions themselves: resolve and reject
	handler := js.FuncOf(func(this js.Value, handlerArgs []js.Value) interface{} {
		resolve := handlerArgs[0]
		reject := handlerArgs[1]

		// Now that we have a way to return the response to JS, spawn a goroutine
		// This way, we don't block the event loop and avoid a deadlock
		go func() {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if ok {
						jsErr := js.Global().Get("Error").New(err.Error())
						reject.Invoke(jsErr)
					} else {
						reject.Invoke()
					}
				}
			}()

			// work
			ret, err := cb()

			// reject if error
			if err != nil {
				jsErr := js.Global().Get("Error").New(err.Error())
				reject.Invoke(jsErr)
			}

			// Resolve the Promise, passing anything back to JavaScript
			// This is done by invoking the "resolve" function passed to the handler
			resolve.Invoke(ret...)
		}()

		// The handler of a Promise doesn't return any value
		return nil
	})
	defer handler.Release()

	// Create and return the Promise object
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func heliaListenAddresses(helia js.Value) (fa []ma.Multiaddr, fb error) {
	defer func() {
		fmt.Println("listen addresses return", fa, fb)
	}()
	jsmaddrs := helia.Get("libp2p").Call("getMultiaddrs")
	l := jsmaddrs.Length()
	ret := make([]ma.Multiaddr, l)
	for i := 0; i < l; i++ {
		jsmaddr := jsmaddrs.Index(i)
		maddrStr := jsmaddr.Call("toString").String()
		maddr, err := ma.NewMultiaddr(maddrStr)
		if err != nil {
			return nil, err
		}
		ret[i] = maddr
	}
	return ret, nil
}

func heliaConnectedPeers(helia js.Value) ([]peer.ID, error) {
	peers := helia.Get("libp2p").Call("getPeers")
	ids := make([]peer.ID, peers.Length())
	for i := 0; i < peers.Length(); i++ {
		p := peers.Index(i)
		rawId := p.Call("toString").String()
		id, err := peer.Decode(rawId)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}
