/*
Copyright 2019 Cortex Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package json

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/lib/files"
)

func Marshal(obj interface{}) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, errStrMarshalJSON)
	}
	return jsonBytes, nil
}

func Unmarshal(data []byte, dst interface{}) error {
	if err := json.Unmarshal(data, dst); err != nil {
		return errors.Wrap(err, errStrUnmarshalJSON)
	}
	return nil
}

func DecodeWithNumber(jsonBytes []byte, dst interface{}) error {
	d := json.NewDecoder(bytes.NewReader(jsonBytes))
	d.UseNumber()
	if err := d.Decode(&dst); err != nil {
		return errors.Wrap(err, errStrUnmarshalJSON)
	}

	return nil
}

func MarshalJSONStr(obj interface{}) (string, error) {
	jsonBytes, err := Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func WriteJSON(obj interface{}, outPath string) error {
	jsonBytes, err := Marshal(obj)
	if err != nil {
		return err
	}
	if err := files.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return err
	}

	if err := files.WriteFile(outPath, jsonBytes, 0644); err != nil {
		return err
	}
	return nil
}

func Pretty(obj interface{}) (string, error) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
