package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func buildRouterIdentity() RouterIdentity {
	router_ident_data := make([]byte, 128+256)
	router_ident_data = append(router_ident_data, []byte{0x05, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00}...)
	return RouterIdentity(router_ident_data)
}

func buildDate() []byte {
	date_data := []byte{0x00, 0x00, 0x00, 0x00, 0x05, 0x26, 0x5c, 0x00}
	return date_data
}

func mappingData() Mapping {
	mapping, _ := GoMapToMapping(map[string]string{"host": "127.0.0.1", "port": "4567"})
	return mapping
}

func buildRouterAddress() RouterAddress {
	router_address_bytes := []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	str, _ := ToI2PString("foo")
	router_address_bytes = append(router_address_bytes, []byte(str)...)
	router_address_bytes = append(router_address_bytes, mappingData()...)
	return RouterAddress(router_address_bytes)
}

func buildFullRouterInfo() RouterInfo {
	router_info_data := make([]byte, 0)
	router_info_data = append(router_info_data, buildRouterIdentity()...)
	router_info_data = append(router_info_data, buildDate()...)
	router_info_data = append(router_info_data, 0x01)
	router_info_data = append(router_info_data, buildRouterAddress()...)
	router_info_data = append(router_info_data, 0x00)
	router_info_data = append(router_info_data, mappingData()...)
	router_info_data = append(router_info_data, make([]byte, 40)...)
	return RouterInfo(router_info_data)
}

func TestPublishedReturnsCorrectDate(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	date, err := router_info.Published()
	assert.Nil(err)
	assert.Equal(int64(86400), date.Time().Unix(), "RouterInfo.Published() did not return correct date")
}

func TestPublishedReturnsCorrectErrorWithPartialDate(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	router_info = router_info[:387+4]
	_, err := router_info.Published()
	if assert.NotNil(err) {
		assert.Equal("error parsing date: not enough data", err.Error())
	}
}

func TestPublishedReturnsCorrectErrorWithInvalidData(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	router_info = router_info[:56]
	_, err := router_info.Published()
	if assert.NotNil(err) {
		assert.Equal("error parsing KeysAndCert: data is smaller than minimum valid size", err.Error())
	}
}

func TestRouterAddressCountReturnsCorrectCount(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	count, err := router_info.RouterAddressCount()
	assert.Nil(err)
	assert.Equal(1, count, "RouterInfo.RouterAddressCount() did not return correct count")
}

func TestRouterAddressCountReturnsCorrectErrorWithInvalidData(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	router_info = router_info[:387+8]
	count, err := router_info.RouterAddressCount()
	if assert.NotNil(err) {
		assert.Equal("error parsing router addresses: not enough data", err.Error())
	}
	assert.Equal(0, count)
}

func TestRouterAddressesReturnsAddresses(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	router_addresses, err := router_info.RouterAddresses()
	assert.Nil(err)
	if assert.Equal(1, len(router_addresses)) {
		assert.Equal(
			0,
			bytes.Compare(
				[]byte(buildRouterAddress()),
				[]byte(router_addresses[0]),
			),
		)
	}
}

func TestRouterAddressesReturnsAddressesWithMultiple(t *testing.T) {}

func TestPeerSizeIsZero(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	size := router_info.PeerSize()
	assert.Equal(0, size, "RouterInfo.PeerSize() did not return 0")
}

func TestOptionsAreCorrect(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	options := router_info.Options()
	assert.Equal(
		0,
		bytes.Compare(
			[]byte(mappingData()),
			[]byte(options),
		),
	)
}

func TestSignatureIsCorrectSize(t *testing.T) {
	assert := assert.New(t)

	router_info := buildFullRouterInfo()
	signature := router_info.Signature()
	assert.Equal(40, len(signature))
}

func TestRouterIdentityIsCorrect(t *testing.T) {}
