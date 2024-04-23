package src

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/setavenger/gobip352"
	"io"
	"io/ioutil"
	"net/http"
)

/*
Most of this will probably be removed in favour of binary encodings (proto buffs)
*/

type Client struct {
	BaseUrl string
}

type Filter struct {
	FilterType  uint8  `json:"filter_type,omitempty"`
	BlockHeight uint64 `json:"block_height,omitempty"`
	BlockHash   []byte `json:"block_hash,omitempty"`
	Data        []byte `json:"data,omitempty"`
}

type UTXO struct {
	Txid         [32]byte `json:"txid"`
	Vout         uint32   `json:"vout"`
	Value        uint64   `json:"value"`
	ScriptPubKey [34]byte `json:"scriptpubkey"`
	BlockHeight  uint64   `json:"block_height"`
	BlockHash    [32]byte `json:"block_hash"`
	Timestamp    uint64   `json:"timestamp"`
}

func (c Client) GetTweaks(blockHeight, dustLimit uint64) ([][33]byte, error) {
	url := fmt.Sprintf("%s/tweaks/%d", c.BaseUrl, blockHeight)
	if dustLimit > 0 {
		url = fmt.Sprintf("%s?dustLimit=%d", url, dustLimit)
	}

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into a []string
	var data []string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Convert []string to [][33]byte
	var bytesData [][33]byte
	for _, hexStr := range data {
		// Each string should be exactly 66 characters long (33 bytes)
		if len(hexStr) != 66 {
			return nil, errors.New(fmt.Sprintf("Invalid hex string length: %d", len(hexStr)))
		}
		// Decode hex string to byte slice
		byteSlice, err := hex.DecodeString(hexStr)
		if err != nil {
			return nil, err
		}
		// Convert byte slice to [33]byte
		var byteArray [33]byte
		copy(byteArray[:], byteSlice[:33]) // Ensure only the first 33 bytes are copied
		bytesData = append(bytesData, byteArray)
	}

	return bytesData, nil
}

func (c Client) GetFilter(blockHeight uint64) (*Filter, error) {
	url := fmt.Sprintf("%s/filter/%d", c.BaseUrl, blockHeight)

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		FilterType  uint8  `json:"filter_type"`
		BlockHeight uint64 `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		Data        string `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	blockHash, err := hex.DecodeString(data.BlockHash)
	if err != nil {
		return nil, err
	}
	filterData, err := hex.DecodeString(data.Data)
	if err != nil {
		return nil, err
	}

	filter := &Filter{
		FilterType:  data.FilterType,
		BlockHeight: data.BlockHeight,
		BlockHash:   blockHash,
		Data:        filterData,
	}

	return filter, err
}

func (c Client) GetUTXOs(blockHeight uint64) ([]*UTXO, error) {
	url := fmt.Sprintf("%s/utxos/%d", c.BaseUrl, blockHeight)

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dataSlice []struct {
		Txid         string `json:"txid"`
		Vout         uint32 `json:"vout"`
		Value        uint64 `json:"value"`
		ScriptPubKey string `json:"scriptpubkey"`
		BlockHeight  uint64 `json:"block_height"`
		BlockHash    string `json:"block_hash"`
		Timestamp    uint64 `json:"timestamp"`
	}

	err = json.Unmarshal(body, &dataSlice)
	if err != nil {
		return nil, err
	}

	var utxos []*UTXO
	for _, data := range dataSlice {
		blockHashBytes, err := hex.DecodeString(data.BlockHash)
		if err != nil {
			return nil, err
		}
		scriptPubKeyBytes, err := hex.DecodeString(data.ScriptPubKey)
		if err != nil {
			return nil, err
		}
		txidBytes, err := hex.DecodeString(data.Txid)
		if err != nil {
			return nil, err
		}

		utxo := &UTXO{
			Txid:         gobip352.ConvertToFixedLength32(txidBytes),
			Vout:         data.Vout,
			Value:        data.Value,
			BlockHeight:  data.BlockHeight,
			BlockHash:    gobip352.ConvertToFixedLength32(blockHashBytes),
			ScriptPubKey: ConvertToFixedLength34(scriptPubKeyBytes),
			Timestamp:    data.Timestamp,
		}

		utxos = append(utxos, utxo)
	}

	return utxos, err
}
