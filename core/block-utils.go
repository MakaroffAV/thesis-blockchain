//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the operations, which use block structure

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import "encoding/json"

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// unMarshal is function for
// converting bytes array to block structure
func unMarshal(b []byte) (*block, error) {

	var block block

	if blockUnMarshalErr := json.Unmarshal(b, &block); blockUnMarshalErr == nil {
		return &block, nil
	} else {
		return nil, blockUnMarshalErr
	}

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// marshal is method of block structure
// for converting block data structure to bytes array
func (b block) marshal() ([]byte, error) {
	return json.Marshal(b)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
