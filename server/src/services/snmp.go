package services

import (
	"fmt"
	"net"
	"time"

	"github.com/gosnmp/gosnmp"
)

type SNMPService struct {
	timeout time.Duration
}

type SNMPCredentials struct {
	Community     string `json:"community"`
	Version       string `json:"version"` // "v1", "v2c", "v3"
	Username      string `json:"username,omitempty"`
	AuthProtocol  string `json:"auth_protocol,omitempty"`
	AuthPassword  string `json:"auth_password,omitempty"`
	PrivProtocol  string `json:"priv_protocol,omitempty"`
	PrivPassword  string `json:"priv_password,omitempty"`
	SecurityLevel string `json:"security_level,omitempty"`
}

type SNMPRequest struct {
	Target      string          `json:"target"`
	Port        int             `json:"port"`
	OIDs        []string        `json:"oids"`
	Credentials SNMPCredentials `json:"credentials"`
}

type SNMPResponse struct {
	Target    string           `json:"target"`
	Success   bool             `json:"success"`
	Variables []gosnmp.SnmpPDU `json:"variables"`
	Error     string           `json:"error,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

type SNMPWalkResponse struct {
	Target    string                   `json:"target"`
	Success   bool                     `json:"success"`
	Results   []map[string]interface{} `json:"results"`
	Error     string                   `json:"error,omitempty"`
	Timestamp time.Time                `json:"timestamp"`
}

func NewSNMPService() *SNMPService {
	return &SNMPService{
		timeout: 10 * time.Second,
	}
}

func (s *SNMPService) GetSNMPData(req SNMPRequest) (*SNMPResponse, error) {
	snmp := &gosnmp.GoSNMP{
		Target:    req.Target,
		Port:      uint16(req.Port),
		Community: req.Credentials.Community,
		Timeout:   s.timeout,
	}

	switch req.Credentials.Version {
	case "v1":
		snmp.Version = gosnmp.Version1
	case "v2c":
		snmp.Version = gosnmp.Version2c
	case "v3":
		snmp.Version = gosnmp.Version3
		snmp.SecurityModel = gosnmp.UserSecurityModel

		switch req.Credentials.SecurityLevel {
		case "noAuthNoPriv":
			snmp.MsgFlags = gosnmp.NoAuthNoPriv
		case "authNoPriv":
			snmp.MsgFlags = gosnmp.AuthNoPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
		case "authPriv":
			snmp.MsgFlags = gosnmp.AuthPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
			snmp.PrivProtocol = s.getPrivProtocol(req.Credentials.PrivProtocol)
			snmp.PrivPassword = req.Credentials.PrivPassword
		}

		snmp.SecurityName = req.Credentials.Username
	default:
		snmp.Version = gosnmp.Version2c
	}

	err := snmp.Connect()
	if err != nil {
		return &SNMPResponse{
			Target:    req.Target,
			Success:   false,
			Error:     fmt.Sprintf("Connection failed: %v", err),
			Timestamp: time.Now(),
		}, nil
	}
	defer snmp.Conn.Close()

	result, err := snmp.Get(req.OIDs)
	if err != nil {
		return &SNMPResponse{
			Target:    req.Target,
			Success:   false,
			Error:     fmt.Sprintf("SNMP Get failed: %v", err),
			Timestamp: time.Now(),
		}, nil
	}

	if result.Error != gosnmp.NoError {
		return &SNMPResponse{
			Target:    req.Target,
			Success:   false,
			Error:     fmt.Sprintf("SNMP error: %v", result.Error),
			Timestamp: time.Now(),
		}, nil
	}

	return &SNMPResponse{
		Target:    req.Target,
		Success:   true,
		Variables: result.Variables,
		Timestamp: time.Now(),
	}, nil
}

func (s *SNMPService) WalkSNMP(req SNMPRequest) (*SNMPWalkResponse, error) {
	snmp := &gosnmp.GoSNMP{
		Target:    req.Target,
		Port:      uint16(req.Port),
		Community: req.Credentials.Community,
		Timeout:   s.timeout,
	}

	switch req.Credentials.Version {
	case "v1":
		snmp.Version = gosnmp.Version1
	case "v2c":
		snmp.Version = gosnmp.Version2c
	case "v3":
		snmp.Version = gosnmp.Version3
		snmp.SecurityModel = gosnmp.UserSecurityModel

		switch req.Credentials.SecurityLevel {
		case "noAuthNoPriv":
			snmp.MsgFlags = gosnmp.NoAuthNoPriv
		case "authNoPriv":
			snmp.MsgFlags = gosnmp.AuthNoPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
		case "authPriv":
			snmp.MsgFlags = gosnmp.AuthPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
			snmp.PrivProtocol = s.getPrivProtocol(req.Credentials.PrivProtocol)
			snmp.PrivPassword = req.Credentials.PrivPassword
		}

		snmp.SecurityName = req.Credentials.Username
	default:
		snmp.Version = gosnmp.Version2c
	}

	err := snmp.Connect()
	if err != nil {
		return &SNMPWalkResponse{
			Target:    req.Target,
			Success:   false,
			Error:     fmt.Sprintf("Connection failed: %v", err),
			Timestamp: time.Now(),
		}, nil
	}
	defer snmp.Conn.Close()

	if len(req.OIDs) == 0 {
		return &SNMPWalkResponse{
			Target:    req.Target,
			Success:   false,
			Error:     "At least one OID is required for walk operation",
			Timestamp: time.Now(),
		}, nil
	}

	result, err := snmp.WalkAll(req.OIDs[0])
	if err != nil {
		return &SNMPWalkResponse{
			Target:    req.Target,
			Success:   false,
			Error:     fmt.Sprintf("SNMP Walk failed: %v", err),
			Timestamp: time.Now(),
		}, nil
	}

	var results []map[string]interface{}
	for _, variable := range result {
		result := map[string]interface{}{
			"oid":       variable.Name,
			"type":      variable.Type.String(),
			"value":     s.convertValue(variable),
			"timestamp": time.Now(),
		}
		results = append(results, result)
	}

	return &SNMPWalkResponse{
		Target:    req.Target,
		Success:   true,
		Results:   results,
		Timestamp: time.Now(),
	}, nil
}

func (s *SNMPService) TestConnection(req SNMPRequest) error {
	snmp := &gosnmp.GoSNMP{
		Target:    req.Target,
		Port:      uint16(req.Port),
		Community: req.Credentials.Community,
		Timeout:   s.timeout,
	}

	switch req.Credentials.Version {
	case "v1":
		snmp.Version = gosnmp.Version1
	case "v2c":
		snmp.Version = gosnmp.Version2c
	case "v3":
		snmp.Version = gosnmp.Version3
		snmp.SecurityModel = gosnmp.UserSecurityModel

		switch req.Credentials.SecurityLevel {
		case "noAuthNoPriv":
			snmp.MsgFlags = gosnmp.NoAuthNoPriv
		case "authNoPriv":
			snmp.MsgFlags = gosnmp.AuthNoPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
		case "authPriv":
			snmp.MsgFlags = gosnmp.AuthPriv
			snmp.AuthProtocol = s.getAuthProtocol(req.Credentials.AuthProtocol)
			snmp.AuthPassword = req.Credentials.AuthPassword
			snmp.PrivProtocol = s.getPrivProtocol(req.Credentials.PrivProtocol)
			snmp.PrivPassword = req.Credentials.PrivPassword
		}

		snmp.SecurityName = req.Credentials.Username
	default:
		snmp.Version = gosnmp.Version2c
	}

	err := snmp.Connect()
	if err != nil {
		return err
	}
	defer snmp.Conn.Close()

	_, err = snmp.Get([]string{"1.3.6.1.2.1.1.1.0"})
	return err
}

func (s *SNMPService) getAuthProtocol(protocol string) gosnmp.SnmpV3AuthProtocol {
	switch protocol {
	case "MD5":
		return gosnmp.MD5
	case "SHA":
		return gosnmp.SHA
	default:
		return gosnmp.MD5
	}
}

func (s *SNMPService) getPrivProtocol(protocol string) gosnmp.SnmpV3PrivProtocol {
	switch protocol {
	case "DES":
		return gosnmp.DES
	case "AES":
		return gosnmp.AES
	case "AES192":
		return gosnmp.AES192
	case "AES256":
		return gosnmp.AES256
	default:
		return gosnmp.AES
	}
}

func (s *SNMPService) convertValue(variable gosnmp.SnmpPDU) interface{} {
	switch variable.Type {
	case gosnmp.OctetString:
		return string(variable.Value.([]byte))
	case gosnmp.Integer:
		return variable.Value.(int)
	case gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Uinteger32:
		return variable.Value.(uint)
	case gosnmp.Counter64:
		return variable.Value.(uint64)
	case gosnmp.ObjectIdentifier:
		return variable.Value.(string)
	case gosnmp.IPAddress:
		return net.IP(variable.Value.([]byte)).String()
	default:
		return variable.Value
	}
}
