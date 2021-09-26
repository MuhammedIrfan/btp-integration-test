/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package base

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/reactivex/rxgo/v2"

	"github.com/icon-project/btp/common/codec"
	"github.com/icon-project/btp/common/log"
)

/*---------------Constants-------------------------------*/
const (
	DefaultGetRelayResultInterval = time.Second
	DefaultRelayReSendInterval    = time.Second
)

/*---------------Types-----------------------------------*/
type sender struct {
	client      Client
	source      BtpAddress
	destination BtpAddress
	wallet      Wallet
	logger      log.Logger
	options     struct {
		StepLimit int64
	}
}

/*------------- Public functions --------------------------*/

func NewSender(source, destination BtpAddress, wallet Wallet, endpoint string, options map[string]interface{}, logger log.Logger, client Client) Sender {
	sender := &sender{
		client:      client,
		source:      source,
		destination: destination,
		wallet:      wallet,
		logger:      logger,
	}

	byteData, err := json.Marshal(options)
	if err != nil {
		logger.Panicf("fail to marshal opt:%#v err:%+v", options, err)
	}

	if err = json.Unmarshal(byteData, &sender.options); err != nil {
		logger.Panicf("fail to unmarshal opt:%#v err:%+v", options, err)
	}

	return sender
}

func (s *sender) MonitorLoop(height int64, cb MonitorCallback, scb func()) error {
	blockRequest := s.client.GetBlockRequest(height)

	return s.client.MonitorSenderBlock(blockRequest,
		func(observable rxgo.Observable) error {
			result := observable.Observe()

			for item := range result {
				notification := item.V.(*BlockNotification)

				if notificationHeight, err := s.client.GetBlockNotificationHeight(notification); err != nil {
					return err
				} else {
					return cb(notificationHeight)
				}
			}
			return nil
		},
		func() {
			scb()
		})
}

func (s *sender) GetResult(p GetResultParam) (TransactionResult, error) {
	for {
		transactionResult, err := s.client.GetTransactionResult(&p)
		if err != nil {
			switch err {
			case ErrGetResultFailByPending:
				<-time.After(DefaultGetRelayResultInterval)
				continue
			}
		}
		return transactionResult, err
	}
}

func (s *sender) GetStatus() (*BMCLinkStatus, error) {
	return s.client.GetBMCLinkStatus(s.wallet, s.destination, s.source)
}

func (s *sender) UpdateSegment(bp *BlockProof, segment *Segment) error {
	transactioParams, err := s.client.GetTransactionParams(segment)
	if err != nil {
		return err
	}

	message, previous, err := s.client.GetRelayMethodParams(&transactioParams)
	if err != nil {
		return err
	}

	clientRelayMessage := &RelayMessageClient{}
	if err := s.client.UnmarshalFromSegment(message, clientRelayMessage); err != nil {
		return err
	}

	if clientRelayMessage.BlockProof, err = codec.RLP.MarshalToBytes(bp); err != nil {
		return err
	}
	segment.TransactionParam, err = s.client.BMCRelayMethodTransactionParam(s.wallet, s.destination, s.source, previous, clientRelayMessage, s.options.StepLimit)

	return err
}

func (s *sender) Relay(segment *Segment) (GetResultParam, error) {
	transactionParams, err := s.client.GetTransactionParams(segment)

	if err != nil {
		return nil, err
	}
	transactionHashParam := new(TransactionHashParam)

SignLoop:
	for {
		if err := s.client.SignTransaction(s.wallet, &transactionParams); err != nil {
			return nil, err
		}

	SendLoop:
		for {
			transactionHash, err := s.client.SendTransaction(&transactionParams)
			if transactionHash != nil {
				if err := s.client.AssignHash(transactionHashParam, transactionHash); err != nil {
					return nil, err
				}
			}

			if err != nil {
				switch err {
				case ErrSendFailByOverflow:
					<-time.After(DefaultRelayReSendInterval)
					continue SendLoop

				case ErrSendDuplicateTransaction:
					s.logger.Debugf("DuplicateTransactionError txh:%v", transactionHash)
					return transactionHashParam, nil

				case ErrSendFailByExpired:
					continue SignLoop
				}
				return nil, err
			}
			return transactionHashParam, nil
		}
	}
}

func (s *sender) Segment(relayMessage *RelayMessage, height int64) ([]*Segment, error) {
	segments := make([]*Segment, 0)
	var err error
	clientRelayMessage := &RelayMessageClient{
		BlockUpdates:  make([][]byte, 0),
		ReceiptProofs: make([][]byte, 0),
	}
	size := 0
	//TODO rm.BlockUpdates[len(rm.BlockUpdates)-1].Height <= s.bmcStatus.Verifier.Height
	//	using only rm.BlockProof
	for _, blockUpdates := range relayMessage.BlockUpdates {
		if blockUpdates.Height <= height {
			continue
		}

		blockUpdateSize := len(blockUpdates.Proof)
		if s.client.IsTransactionOverLimit(blockUpdateSize) {
			return nil, fmt.Errorf("invalid BlockUpdate.Proof size")
		}

		size += blockUpdateSize
		if s.client.IsTransactionOverLimit(size) {
			segment := &Segment{
				Height:              clientRelayMessage.height,
				NumberOfBlockUpdate: clientRelayMessage.numberOfBlockUpdate,
			}

			if segment.TransactionParam, err = s.client.BMCRelayMethodTransactionParam(s.wallet, s.destination, s.source, relayMessage.From.String(), clientRelayMessage, s.options.StepLimit); err != nil {
				return nil, err
			}

			segments = append(segments, segment)
			clientRelayMessage = &RelayMessageClient{
				BlockUpdates:  make([][]byte, 0),
				ReceiptProofs: make([][]byte, 0),
			}
			size = blockUpdateSize
		}

		clientRelayMessage.BlockUpdates = append(clientRelayMessage.BlockUpdates, blockUpdates.Proof)
		clientRelayMessage.height = blockUpdates.Height
		clientRelayMessage.numberOfBlockUpdate += 1
	}

	var blockProof []byte
	if blockProof, err = codec.RLP.MarshalToBytes(relayMessage.BlockProof); err != nil {
		return nil, err
	}

	if s.client.IsTransactionOverLimit(len(blockProof)) {
		return nil, fmt.Errorf("invalid BlockProof size")
	}

	var byteData []byte
	for _, receiptProof := range relayMessage.ReceiptProofs {
		if s.client.IsTransactionOverLimit(len(receiptProof.Proof)) {
			return nil, fmt.Errorf("invalid ReceiptProof.Proof size")
		}

		if len(clientRelayMessage.BlockUpdates) == 0 {
			size += len(blockProof)
			clientRelayMessage.BlockProof = blockProof
			clientRelayMessage.height = relayMessage.BlockProof.BlockWitness.Height
		}

		size += len(receiptProof.Proof)
		trp := &ReceiptProof{
			Index:       receiptProof.Index,
			Proof:       receiptProof.Proof,
			EventProofs: make([]*EventProof, 0),
		}

		for j, eventProof := range receiptProof.EventProofs {
			if s.client.IsTransactionOverLimit(len(eventProof.Proof)) {
				return nil, fmt.Errorf("invalid EventProof.Proof size")
			}

			size += len(eventProof.Proof)
			if s.client.IsTransactionOverLimit(size) {
				if j == 0 && len(clientRelayMessage.BlockUpdates) == 0 {
					return nil, fmt.Errorf("BlockProof + ReceiptProof + EventProof > limit")
				}

				segment := &Segment{
					Height:              clientRelayMessage.height,
					NumberOfBlockUpdate: clientRelayMessage.numberOfBlockUpdate,
					EventSequence:       clientRelayMessage.eventSequence,
					NumberOfEvent:       clientRelayMessage.numberOfEvent,
				}

				if segment.TransactionParam, err = s.client.BMCRelayMethodTransactionParam(s.wallet, s.destination, s.source, relayMessage.From.String(), clientRelayMessage, s.options.StepLimit); err != nil {
					return nil, err
				}
				segments = append(segments, segment)

				clientRelayMessage = &RelayMessageClient{
					BlockUpdates:  make([][]byte, 0),
					ReceiptProofs: make([][]byte, 0),
					BlockProof:    blockProof,
				}

				size = len(eventProof.Proof)
				size += len(receiptProof.Proof)
				size += len(blockProof)

				trp = &ReceiptProof{
					Index:       receiptProof.Index,
					Proof:       receiptProof.Proof,
					EventProofs: make([]*EventProof, 0),
				}
			}

			trp.EventProofs = append(trp.EventProofs, eventProof)
			clientRelayMessage.eventSequence = receiptProof.Events[j].Sequence
			clientRelayMessage.numberOfEvent += 1
		}

		if byteData, err = codec.RLP.MarshalToBytes(trp); err != nil {
			return nil, err
		}

		clientRelayMessage.ReceiptProofs = append(clientRelayMessage.ReceiptProofs, byteData)
	}

	segment := &Segment{
		Height:              clientRelayMessage.height,
		NumberOfBlockUpdate: clientRelayMessage.numberOfBlockUpdate,
		EventSequence:       clientRelayMessage.eventSequence,
		NumberOfEvent:       clientRelayMessage.numberOfEvent,
	}

	if segment.TransactionParam, err = s.client.BMCRelayMethodTransactionParam(s.wallet, s.destination, s.source, relayMessage.From.String(), clientRelayMessage, s.options.StepLimit); err != nil {
		return nil, err
	}

	segments = append(segments, segment)
	return segments, nil
}

func (s *sender) StopMonitorLoop() {
	s.client.CloseAllMonitor()
}

func (s *sender) FinalizeLatency() int {
	return 1
}
