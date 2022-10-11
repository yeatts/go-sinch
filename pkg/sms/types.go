package sms

import "encoding/json"

type DeliveryReport int

const (
	None DeliveryReport = iota
	Summary
	Full
	PerRecipient
)

func (dr DeliveryReport) String() string {
	return dr.toString()
}

func (dr DeliveryReport) toString() string {
	switch dr {
	case None:
		return "none"
	case Summary:
		return "summary"
	case Full:
		return "full"
	case PerRecipient:
		return "per_recipient"
	}
	return "none"
}

func toDeliveryReport(s string) DeliveryReport {
	switch s {
	case "none":
		return None
	case "summary":
		return Summary
	case "full":
		return Full
	case "per_recipient":
		return PerRecipient
	}
	return None
}

func (dr DeliveryReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(dr.String())
}

func (dr *DeliveryReport) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*dr = toDeliveryReport(s)
	return nil
}

type Type int

const (
	Text Type = iota
	Binary
)

func (t Type) String() string {
	return t.toString()
}

func (t Type) toString() string {
	switch t {
	case Text:
		return "mt_text"
	case Binary:
		return "mt_binary"
	}
	return "mt_text"
}

func toType(s string) Type {
	switch s {
	case "mt_text":
		return Text
	case "mt_binary":
		return Binary
	}
	return Text
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*t = toType(s)
	return nil
}
