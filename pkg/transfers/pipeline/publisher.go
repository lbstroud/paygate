// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/moov-io/paygate/pkg/config"
)

// XferPublisher is an interface for pushing Transfers (and their ACH files) to be
// uploaded to an ODFI. These implementations can be to push Transfers onto streams
// (e.g. kafka, rabbitmq) or inmem (the default in our OSS PayGate).
type XferPublisher interface {
	Upload(xfer Xfer) error
	Cancel(msg CanceledTransfer) error
	Shutdown(ctx context.Context)
}

func NewPublisher(cfg config.Pipeline) (XferPublisher, error) {
	if cfg.Stream != nil {
		return createStreamPublisher(cfg.Stream)
	}
	return nil, errors.New("unknown Pipeline config")
}

func createMetadata(xf Xfer) map[string]string {
	out := make(map[string]string)
	out["transferID"] = xf.Transfer.TransferID
	return out
}

func createBody(xf Xfer) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(xf); err != nil {
		return nil, fmt.Errorf("trasferID=%s json encode: %v", xf.Transfer.TransferID, err)
	}
	return buf.Bytes(), nil
}
