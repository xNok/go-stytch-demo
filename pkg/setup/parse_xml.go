package setup

import (
	"encoding/xml"
)

// Define XML structures corresponding to the XML document structure
type EntityDescriptor struct {
	XMLName          xml.Name         `xml:"EntityDescriptor"`
	EntityID         string           `xml:"entityID,attr"`
	IDPSSODescriptor IDPSSODescriptor `xml:"IDPSSODescriptor"`
}

type IDPSSODescriptor struct {
	WantAuthnRequestsSigned    string                `xml:"WantAuthnRequestsSigned,attr"`
	ProtocolSupportEnumeration string                `xml:"protocolSupportEnumeration,attr"`
	KeyDescriptors             []KeyDescriptor       `xml:"KeyDescriptor"`
	NameIDFormat               string                `xml:"NameIDFormat"`
	SingleSignOnServices       []SingleSignOnService `xml:"SingleSignOnService"`
}

type KeyDescriptor struct {
	Use     string  `xml:"use,attr"`
	KeyInfo KeyInfo `xml:"KeyInfo"`
}

type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

type SingleSignOnService struct {
	Binding  string `xml:"Binding,attr"`
	Location string `xml:"Location,attr"`
}

func parseXML(xmlData string) (*EntityDescriptor, error) {
	var descriptor EntityDescriptor

	err := xml.Unmarshal([]byte(xmlData), &descriptor)

	return &descriptor, err
}
