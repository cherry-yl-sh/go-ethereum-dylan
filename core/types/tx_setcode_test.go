// Copyright 2024 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestParseDelegation tests a few possible delegation designator values and
// ensures they are parsed correctly.
func TestParseDelegation(t *testing.T) {
	addr := common.Address{0x42}
	for _, tt := range []struct {
		val  []byte
		want *common.Address
	}{
		{ // simple correct delegation
			val:  append(DelegationPrefix, addr.Bytes()...),
			want: &addr,
		},
		{ // wrong address size
			val: append(DelegationPrefix, addr.Bytes()[0:19]...),
		},
		{ // short address
			val: append(DelegationPrefix, 0x42),
		},
		{ // long address
			val: append(append(DelegationPrefix, addr.Bytes()...), 0x42),
		},
		{ // wrong prefix size
			val: append(DelegationPrefix[:2], addr.Bytes()...),
		},
		{ // wrong prefix
			val: append([]byte{0xef, 0x01, 0x01}, addr.Bytes()...),
		},
		{ // wrong prefix
			val: append([]byte{0xef, 0x00, 0x00}, addr.Bytes()...),
		},
		{ // no prefix
			val: addr.Bytes(),
		},
		{ // no address
			val: DelegationPrefix,
		},
	} {
		got, ok := ParseDelegation(tt.val)
		if ok && tt.want == nil {
			t.Fatalf("expected fail, got %s", got.Hex())
		}
		if !ok && tt.want != nil {
			t.Fatalf("failed to parse, want %s", tt.want.Hex())
		}
	}

}

// https://odyssey-explorer.ithaca.xyz/tx/0xe3a7660b86626560968e63bd9faa428dd6555cc6c09054c144295d2b6c7e2da9?tab=index
func TestAuthority(t *testing.T) {
	jsonData := `{
        "address": "0xadeebe459e44222ed40fa615be9a929d2fa77893",
        "chainId": "0xde9fb",
        "nonce": "0x1",
        "r": "0x380c4db8e1b82461e4b6b235c775649995f8b06e0b999fa69b5633f539f1352c",
        "s": "0x50a11fa0d2f2e66bb8a95948db232ce24a22913973dcdb50e3b1cdc697e032d1",
        "yParity": "0x1"
    }`

	// 创建SetCodeAuthorization实例
	var auth SetCodeAuthorization

	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		panic(err)
	}
	fmt.Printf("SetCodeAuthorization:\n")
	fmt.Printf("  Address: %s\n", auth.Address.Hex())
	fmt.Printf("  ChainID: %s\n", auth.ChainID.String())
	fmt.Printf("  Nonce: %d\n", auth.Nonce)
	fmt.Printf("  V: %d\n", auth.V)
	fmt.Printf("  R: %s\n", auth.R.Hex())
	fmt.Printf("  S: %s\n", auth.S.Hex())
	addr, err := auth.Authority()
	if err != nil {
		fmt.Printf("Authority error: %v\n", err)
	}
	fmt.Printf("Authority: %s\n", addr.Hex())
	expectedAddr := common.HexToAddress("0x7c0ea167b05a85c6e6ce7e919983af3f3cea379e")
	assert.Equal(t, expectedAddr, addr)

}
