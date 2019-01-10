/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package eventmanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-evm/event"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm"
	"github.com/hyperledger/burrow/execution/exec"
)

type EventManager struct {
	stub       shim.ChaincodeStubInterface
	EventCache []event.Event
}

var _ evm.EventSink = &EventManager{}

type EventSink interface {
	Call(call *exec.CallEvent, exception *errors.Exception) error
	Log(log *exec.LogEvent) error
}

func NewEventManager(stub shim.ChaincodeStubInterface) *EventManager {
	return &EventManager{
		stub:       stub,
		EventCache: []event.Event{},
	}
}

// eventName is for fabric, typically the evm 8byte function hash
func (evmgr *EventManager) Flush(eventName string) error {
	if len(evmgr.EventCache) == 0 {
		return nil
	}
	payload, err := json.Marshal(evmgr.EventCache)
	if err != nil {
		return fmt.Errorf("Failed to marshal event messages: %s", err)
	}
	return evmgr.stub.SetEvent(eventName, payload)
}

// noop for now, need to figure out what it means (burrow or evm)
func (evmgr *EventManager) Call(call *exec.CallEvent, exception *errors.Exception) error {
	return nil
}

func (evmgr *EventManager) Log(log *exec.LogEvent) error {
	e := event.Event{Address: strings.ToLower(log.Address.String()), Data: strings.ToLower(log.Data.String())}

	topics := []string{}
	for _, topic := range log.Topics {
		t, err := topic.MarshalText()
		if err != nil {
			return fmt.Errorf("Failed to Marshal Topic: %s", err)
		}
		topics = append(topics, strings.ToLower(string(t)))
	}
	e.Topics = topics

	evmgr.EventCache = append(evmgr.EventCache, e)

	return nil
}
