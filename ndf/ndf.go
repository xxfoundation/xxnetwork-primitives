////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Package ndf contains the structure for the network definition, which matches
// the JSON encoded network definition file (NDF). This file is passed, in some
// form, to all members of the network to relay connection and network
// information.
// The NDF also contains a base64 encoded signature of the network definition
// information.

package ndf

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// NO_NDF is a string that the permissioning server responds with when a member
// of the network requests an NDF from it but the NDF is not yet available.
const NO_NDF = "Permissioning server does not have an ndf to give"

// NetworkDefinition structure hold connection and network information. It
// matches the JSON structure generated in Terraform.
type NetworkDefinition struct {
	Timestamp    time.Time
	Gateways     []Gateway
	Nodes        []Node
	Registration Registration
	Notification Notification
	UDB          UDB   `json:"Udb"`
	E2E          Group `json:"E2e"`
	CMIX         Group `json:"Cmix"`
}

// Gateway contains the connection and identity information of a gateway on the
// network.
type Gateway struct {
	ID             []byte `json:"Id"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Node contains the connection and identity information of a node on the
// network.
type Node struct {
	ID             []byte `json:"Id"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Registration contains the connection information for the permissioning
// server.
type Registration struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Notification contains the connection information for the notification bot.
type Notification struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// UDB contains the ID and public key in PEM form for user discovery.
type UDB struct {
	ID        []byte `json:"Id"`
	PubKeyPem string `json:"Public_key_PEM"`
}

// Group contains the information used to reconstruct a cyclic group.
type Group struct {
	Prime      string
	SmallPrime string `json:"Small_prime"`
	Generator  string
}

// StripNdf returns a stripped down copy of the NetworkDefinition to be used by
// Clients.
func (ndf *NetworkDefinition) StripNdf() *NetworkDefinition {
	// Remove address and TLS cert for every node.
	var strippedNodes []Node
	for _, node := range ndf.Nodes {
		newNode := Node{
			ID: node.ID,
		}
		strippedNodes = append(strippedNodes, newNode)
	}

	// Create a new NetworkDefinition with the stripped information
	return &NetworkDefinition{
		Timestamp:    ndf.Timestamp,
		Gateways:     ndf.Gateways,
		Nodes:        strippedNodes,
		Registration: ndf.Registration,
		Notification: ndf.Notification,
		UDB:          ndf.UDB,
		E2E:          ndf.E2E,
		CMIX:         ndf.CMIX,
	}
}

// GetNodeId unmarshalls the node's ID bytes into an id.ID and returns it.
func (n *Node) GetNodeId() (*id.ID, error) {
	return id.Unmarshal(n.ID)
}

// GetGatewayId unmarshalls the gateway's ID bytes into an id.ID and returns it.
func (g *Gateway) GetGatewayId() (*id.ID, error) {
	return id.Unmarshal(g.ID)
}

// GetGatewayId unmarshalls the UDB ID bytes into an id.ID and returns it.
func (u *UDB) GetUdbId() (*id.ID, error) {
	return id.Unmarshal(u.ID)
}

// NewTestNDF generates a sample NDF used for testing.
func NewTestNDF(i interface{}) *NetworkDefinition {
	switch i.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic(fmt.Sprintf("NewTestNDF is restricted to testing only. Got %T", i))
	}

	ndfEncoded := "eyJUaW1lc3RhbXAiOiIyMDIwLTEwLTI2VDE2OjA5OjU0Ljk5ODQwOTE5WiIsIkdhdGV3YXlzIjpbeyJJZCI6IkVwWUdDclRESjFwV3U1d0c4VmF4cElFUW1jWFNaWnhqcWgxbW1GdjArOGNCIiwiQWRkcmVzcyI6IjEyNy4wLjAuMToyMjg0MCIsIlRsc19jZXJ0aWZpY2F0ZSI6Ii0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLVxuTUlJRnpUQ0NBN1dnQXdJQkFnSVVEL2lCdDhIK2NCckx5UkZoRzFyVm9wT2cvZUl3RFFZSktvWklodmNOQVFFTFxuQlFBd2daSXhDekFKQmdOVkJBWVRBazVNTVJJd0VBWURWUVFJREFsR2JHVjJiMnhoYm1ReEVUQVBCZ05WQkFjTVxuQ0V4bGJIbHpkR0ZrTVJJd0VBWURWUVFLREFsNGVHNWxkSGR2Y21zeEVqQVFCZ05WQkFzTUNYUmxjM1JPYjJSbFxuY3pFVE1CRUdBMVVFQXd3S2VIZ3VibVYwZDI5eWF6RWZNQjBHQ1NxR1NJYjNEUUVKQVJZUVlXUnRhVzVBZUhndVxuYm1WMGQyOXlhekFlRncweU1EQTVNamt4TURBMk5ERmFGdzB5TWpBNU1qa3hNREEyTkRGYU1JR1NNUXN3Q1FZRFxuVlFRR0V3Sk9UREVTTUJBR0ExVUVDQXdKUm14bGRtOXNZVzVrTVJFd0R3WURWUVFIREFoTVpXeDVjM1JoWkRFU1xuTUJBR0ExVUVDZ3dKZUhodVpYUjNiM0pyTVJJd0VBWURWUVFMREFsMFpYTjBUbTlrWlhNeEV6QVJCZ05WQkFNTVxuQ25oNExtNWxkSGR2Y21zeEh6QWRCZ2txaGtpRzl3MEJDUUVXRUdGa2JXbHVRSGg0TG01bGRIZHZjbXN3Z2dJaVxuTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElDRHdBd2dnSUtBb0lDQVFDMjlQU1JUZ05ISkFNU3dZREluQTk1QTZyTVxuQzJIZzRKOWxBWjJNNExMQzRWOXZGYTVvR0M4ODhkaFM2UmlUUUNXNTNWbGZuZHhwaEdLSkJveVJJK1VTenRsZlxuT25zR3BQd1ZCOWE0Uy9rdmNrL2hSOU01eEZHcElEQkxDektTMk5ocGVQUk5yclFHZXR2cWJ5U1JBQzkzd0k5YVxub3dteTdCcjArZW9NcTh6enhLMWpFUlJicHJDdnJGNHV2ait1Z0t2V1V3UWtsV002QWk1RDJzanVLam1hUE0yNlxuUEZpZ01VQXQ5UUtyOWFKYVZtdE50cFY4T29hSzZkQWlJekRuUFZWZExjRDlSeUVHVHVGL0lLZCt4ZXUzNW0yYlxuTTlla25rbGo4MWwrOHVsK3FuakhYS3lJRW9NWVp3OXUzZFQ2WXE3VlFWK3VkMXQ3T3NxVkZsbmhtUVhjT2VUaVxudzI5YnVHdjE0clVEem5DZ1Jvd0c5Qzg5L0NLRmhVTlNTQkhzclhORjkzOHFqZmw5R1I4NHZ6QnRpbXdnZ24ySlxuRStQRW5FRDFTM2FhUkdLU1hsQ1hSWFE3MUFRNVVJa2Z0dkdZSWN4OHdBUDFQSE1jNHQvYWcvNEwyWEV2V3BmelxuK1gzZFozVzk4dm1WOUV4aDN1dWQ2dEU4NGEyc3Z2SldTN1A4Y0ZiYmppbDdxaWk4Y0tZY2kwb0dqVUtBcTlIY1xuQ25uT1BkQ25qcnF4eXFnOElkRmVyb2hiTVRQZ1hENTZMZSsvYWJ0eEUvdGZwdENqRlNUNHN4K0twWGxSUGRwbFxuQXdUZUlkaXcrU1hIQjdvNENJL3hjZmJqT3FFb2FZTDk2WE0yVjVtN3grTHdzTGptbHkvcWpCZHp3VVV0czFsV1xuNlBHN1JNY3AyR0hVbk5TLzZ3SURBUUFCb3hrd0Z6QVZCZ05WSFJFRURqQU1nZ3A0ZUM1dVpYUjNiM0pyTUEwR1xuQ1NxR1NJYjNEUUVCQ3dVQUE0SUNBUUJTLzRtMEhtU2dyS3lqWGhMV3p2RHZnVlRnMm5RSGxsNXpZRWZVbTlTbFxuRU8ydVNGdmlBMmlTcE0zVWphck9JSG1FOXEweWNKSFZCZDlrdEhFdHdwTHhkVjg0ZFZlN25GdzcwZEpwekhTclxuQnlTNHl1RklTaHErTUlpU0JYOFI1MVA1VHgwOGxSVTJtalphQ1pqYy9OQ2tLNHY3ZU9zUDdzQzNsejQ1WldVMFxuTDB2NUU1Q1E3WDBwbFMzZ3VReTlaYzQ0TUVBTGp0Y0lQcjdXWWtzL1hDUUtld2l6Y1J5elA1NXFOeEhVNjRIK1xuaWpqY0s0cWVIdmhlcUg1SlZ4Zk92K212cXU3emNyQXhGRmkrU0UvWDV3UnBuR2VHV2kvcnhZWXVJblovdTRXV1xub3hzcmNVRDFvbVJTYTVoNlRxRkpTck4yYy9vZXdGK3RrK01FbjNlWnQ1cDB5MFpMeFAxaVdJd3A0bGVpM1hHVVxud3dZNXlNamF6OU1PdmZ0Wm82UkQ3dng2RzdtY3lPbDVlR2lTSEE4eloyN2ZKNytPZWJJRXIwOEo3VGdkMk5vUlxuM0xJQVUxdjI1ZE8zZnU4b29RU3hmR3VZenZKSERaWTl4bVRBTk5HVWNnVExSd2JOVzZYbzV4MjB1OTBxQ2ZBRVxubmpoMmpYWjRYbFFwb3hsTVVkZmFoNFhxK0NNT0J5UzQ4b0FwQ0QrU2NyUEVRcC9MMlhlSVM2NE96b3BtN2pONFxuQThkU29JUUg4a0YwK0taSXUvSEdlSDVETWxWMUJBM1RCNmhOU3hZY3BXNHB3NXNJM0hlbWROQ2xxdW8wN3NXUlxud3pleUpvWEhaL3lXSkt5bm5aZWRUZExuSngyeDJpcXlQa1lNTFdmZnZFck50NThBem0rQUtSa0xRUTd0a2g2NlxudFE9PVxuLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLVxuIn0seyJJZCI6Ikh3eEo4eVEyUjJscW1BOHVZMzNsK0c5U0RvbSszbUhGNVdIR05GUmMwd3dCIiwiQWRkcmVzcyI6IjEyNy4wLjAuMToyMjg0MCIsIlRsc19jZXJ0aWZpY2F0ZSI6Ii0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLVxuTUlJRnV6Q0NBNk9nQXdJQkFnSVVSeDBsL3k5T1NwNXVMRDZ4MURnZ0tZMVhhSDR3RFFZSktvWklodmNOQVFFTFxuQlFBd2dZa3hDekFKQmdOVkJBWVRBbE5MTVFzd0NRWURWUVFJREFKVFN6RVRNQkVHQTFVRUJ3d0tRbkpoZEdselxuYkdGMllURVNNQkFHQTFVRUNnd0plSGh1WlhSM2IzSnJNUTR3REFZRFZRUUxEQVZ1YjJSbGN6RVRNQkVHQTFVRVxuQXd3S2VIZ3VibVYwZDI5eWF6RWZNQjBHQ1NxR1NJYjNEUUVKQVJZUVlXUnRhVzVBZUhndWJtVjBkMjl5YXpBZVxuRncweU1EQTJNamN4TmpBM05EQmFGdzB5TWpBMk1qY3hOakEzTkRCYU1JR0pNUXN3Q1FZRFZRUUdFd0pUU3pFTFxuTUFrR0ExVUVDQXdDVTBzeEV6QVJCZ05WQkFjTUNrSnlZWFJwYzJ4aGRtRXhFakFRQmdOVkJBb01DWGg0Ym1WMFxuZDI5eWF6RU9NQXdHQTFVRUN3d0ZibTlrWlhNeEV6QVJCZ05WQkFNTUNuaDRMbTVsZEhkdmNtc3hIekFkQmdrcVxuaGtpRzl3MEJDUUVXRUdGa2JXbHVRSGg0TG01bGRIZHZjbXN3Z2dJaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQ1xuRHdBd2dnSUtBb0lDQVFEWmhtUGVtdnFvZXNnQWI5TGpGVGdoVGdOTmx0RnRldDFBM2c3ejc4VVFIbjNOYUdjRVxuWUZDQXZJK0ZxOEowWElSMFp1UmJrdnNJVmlXOUsxa3NUaFdnckFIQzJTcHA1Yk9GN01qTGUreGY2MS9wbmlUVVxuemdGZlBOU0w5aVRRN2libkw1TFl6OFNwVHpmOFpkK0RzYmJXc2RvTTJmd2paYmovNlVUM3ZReFFxL2gxNWhORlxuQktucEtmdE43bVFFR1dlMm1yVWhrVGF5aVZCU1VwVTduRTVabFdLK2RKQTM4RUVpcVpOeXVPdlZ2U1ZDQzRnTVxuTE11NytlRUx1V3BBMjdpMVpGRFEybzY5WW05QXZvemNzQU5qb3o2eSs5c2k5NnlrWjlQNGRKb2Y5bmVLOVl0ZVxuYkVXM1lRL0ZNeGgrMHloWXBFQ2hmL0ppK2t4Z0s5Tkk5OGJMREEzNzFrZGtLTXdCeHJ2OEZsUTBwSnhLSHAwT1xuRVgwSm1aVTR3VDZQVmwwUENTeGJFenFHSUhMbDM0TVlnM0s4SUJKT2tSU2pTY1gyZlJkbDQyVm9DVURmMHVaMlxuS1U1Vk1TeXpLZVdKZGhVWEphQmgzMUJSb3FlYkQ5Y1pIdVNtcTkvZG1JaXJHU21DaWdodkZRcTRQQ1E1blJJL1xuQ0tKTXBYWms2L3JjMnk2UFRVb0FNWVdha2NhNGxTcnZYTGpJUEZDTThXTFhoMGFlMWU2S3dMQXk2a0ZCcitVSlxuQ2Y1S3hwT2VPTEVsNTF5cEd2YnlkdXhweVJwd280MnU1OEUyUDQzd2NwaE5tZWcxaHNkVWtwK0FQOTYxLzZBM1xudTR4a1BvVk16ZmFiL05MTlg0Sm9KL3BRV1RLdGFFRE1RVkpTMXVmTDNnQ3JFN1Zwdkc5VWFYZCtFd0lEQVFBQlxub3hrd0Z6QVZCZ05WSFJFRURqQU1nZ3A0ZUM1dVpYUjNiM0pyTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFBSlxuZ0todFp5Q1VtL0NScW5nS0NaRGZpR21yQVl4a0ZaMmlhY1VIZ2ZtUnc0blZrTGQ3K093OWorNU9uSGRKaVRYM1xuSEpWNmVaNDZCcGNKN1NXWWd4dk9HdytBN2lZcjVkaWRaSzkwRGtmRWZlTHNBcXgzMEM3b25pNGVMNXBWdi9BdFxuaFd1MU9hdW1GdG5EWUdLU2ZhUXN1RGNYNlhjYjRXRG1JRUN6Q3RJTGFDeXZySnRxRElTOHdERWhIVTZaeVJtVVxuS0pEbnh5bkoyZEZWZVhxZnUwdmdudHFiYkUzR0VObEJHclAvWVJSeFZ1ajFPb0tsVHpZVFBpMXJLS0poaGNmMVxuaXp6RE40dmVueUxRckh4QXJhcDNndVFXVThFNFdlZUtsWW56Y1Y2TEtZQUNCd3ZuYUlsZ3VVTjRVVG41MW9PNVxud0RMK2ZjZm12dHBaejNGbGMzNlE1RVlwS2lZMi85dENtTmNIWENRZnhmaGR4STJ4UmFTUlhmdVQ5UUViZlAxVlxuU2FQNVN5UFFOWGZTdGMwQjI3bnFLaFhENDJkeFpXdlBSTDNtQ3p6dHhUZWljNzgvZWZKM05zejFvNkFXVEoyb1xubkkyYWptZTZKNjdSeTVPRTg1dDRnZlNNeDg1R3hEb1dWSFRPZDF6YklJdURkQU90Vi96TnphZ2tia3N4SkcraVxuaEhhbUloOTBBbUVMUjBqVjZiVlhBYVdreElwTjlFUmM2Qmc4d1FpVCtMYTczNkdBTUVUUnRlOERXdWVVL3V2Z1xuR05Yd0dQcGkzaHppN3pISEo2d0J1Q2hyMkZGeG5WY3NkK3VjTnN3OEZaTll1V3N1SC9GYm5raXpwSTF1ZjZxdlxuUGFITnZ3dU90M3BxdlNQWnpIQTM1eVpSTHorVVJRbzhqOC9oakxNUXZnPT1cbi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS1cbiJ9LHsiSWQiOiJ0UElEcWQ3bWJqeTQ2TmFZcjdQK1NpdlpkTjNpczk2eGZFUXBTQ2xSTUxvQiIsIkFkZHJlc3MiOiIxNzguMTI0LjIxMS4yMDE6MjI4NDAiLCJUbHNfY2VydGlmaWNhdGUiOiItLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS1cbk1JSUYxVENDQTcyZ0F3SUJBZ0lVVGZIa0xVWm5vTjl5NzJHSUlKM0xtNW03QTZBd0RRWUpLb1pJaHZjTkFRRUxcbkJRQXdnWll4Q3pBSkJnTlZCQVlUQWtKWk1SZ3dGZ1lEVlFRSURBOU5hVzV6YTJGNVlTQnZZbXhoYzNReERqQU1cbkJnTlZCQWNNQlUxcGJuTnJNUkl3RUFZRFZRUUtEQWw0ZUc1bGRIZHZjbXN4RWpBUUJnTlZCQXNNQ1hSbGMzUk9cbmIyUmxjekVUTUJFR0ExVUVBd3dLZUhndWJtVjBkMjl5YXpFZ01CNEdDU3FHU0liM0RRRUpBUllSTXpRMU56VXlcbk5rQm5iV0ZwYkM1amIyMHdIaGNOTWpBd056STNNVFV3TmpRNVdoY05Nakl3TnpJM01UVXdOalE1V2pDQmxqRUxcbk1Ba0dBMVVFQmhNQ1Fsa3hHREFXQmdOVkJBZ01EMDFwYm5OcllYbGhJRzlpYkdGemRERU9NQXdHQTFVRUJ3d0ZcblRXbHVjMnN4RWpBUUJnTlZCQW9NQ1hoNGJtVjBkMjl5YXpFU01CQUdBMVVFQ3d3SmRHVnpkRTV2WkdWek1STXdcbkVRWURWUVFEREFwNGVDNXVaWFIzYjNKck1TQXdIZ1lKS29aSWh2Y05BUWtCRmhFek5EVTNOVEkyUUdkdFlXbHNcbkxtTnZiVENDQWlJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dJUEFEQ0NBZ29DZ2dJQkFKdkJyZUxBTWtBcXEwbEZcbm1KUDJpdGtnQXl2aWJzeDFpYm1vTlh4MGE0ZFNFS3d0d1ZVRnU0dEVDUGxVSmJ0emtQNDdRUzlkdXF0M1ZsKzFcbkF5cklQYVdhOEZpd1RjZkZRY0pWeE9zUVVzL3lzREFseVAzK0tKUmRta2JKc0FDRm1mVjA2TGJWMDZzdFNhZlhcbmEyb0dlNndrRGNjRkRsR2x4RHRrREd4SVh2TlJadUY1cGdsWDNuOUFVbG8yeXU3TmxlKytFbkJJOFVuNnBFWktcbmZZbDM5Ym5WTFdNNDlSZy93dkdGQ3ZGNjRSWDN4b1QrTU1rWm5Eb3FxY25jUTRXRnA2RHRXUVdWZE5GMzVPVnhcbjFlWkRQMkRaTEdCWWlva1ZxakxvOEhUb01XcDJHQ1VreUtoSmZ3R2J2amNBbDFPaHJncFVFdndhSkxoS1NVeEtcbjd1Q2dZOEpxZ0E5V2RPNTJZK05kdVRpcVJ0WkIwbjVlSkVlMWJWaWRDL1JkRTlFZGVmMTNDVUFlZVNERGkyWHBcbjVnRGI2T0VNRnROWlRFbythOVM1ODA2ZytyU2s2WlpHbW5VL2c1anlmMkVIUzBrd3RzU0NRL1JNUWpKUnZJMUZcbjFvWVYxYkFVNy9MRDMyRXFnMExQbVlvTllsOUJBUVlzOHlmVFJxVjhCYzI4a1prVm9BTm1yOTFQUDRrQlFBWjlcbmpmeE5QdWhPekM5S0xzczIvcmNPNWM1NmtqMFd5N0p3ZjlyTmVTRHFzS2VXWDFwV1BScVZlNGZSelQ1a3lEcmZcbnJ3YnVEdmtUN3NDUXhqSnF4aXFCSDhXSWUyc1AxeUFWU1l0U2lqSFd1M3kxNmlBVlB4VmFJcjdBWklPT05YQU1cbmVjWlIrZ3JCK1VGb1RZME5sZk9NVExLQ1ovVXBBZ01CQUFHakdUQVhNQlVHQTFVZEVRUU9NQXlDQ25oNExtNWxcbmRIZHZjbXN3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUd5ZU96N0ZRR2hEdnNNc3lPREZGYmgzSXc5V0xVWW5cbm93YUp2azhScmxkR1JIUDVoRS9wcWRqZmJqTGU2YjJPdzN0eGEzaW1jSy9HSGpFU1paTnE1aUYxenRXWGVvcEJcbnhWYVAvNk1oRzlFeWh5Q05FbUNWaldVYXdFT1duS0tOeWUvWGNpbEVTY1pDT0g2b3RILzh1SCswVU5MaWFmQ0VcblhjK0tyTmRENHRGREZtQlRXNytqanliTG1jY2tQeEZzT3BvQ0sybTg4OGpKMW5HUmI2dm9TV25pTU5LcXlESDBcblhpb3hGdDZsY2dZWHpUbmplVHkrRzRYR29OM3Jkb29mL05nRk1Tbmc4TlN1SGxYWHNYOEQvTEVHcWU3TkdQc3RcbmVkR2YrTTJacHovRDJNN00zc0t6aWQ1V3RBS0VpSVg3YnIwN3RoRTM0dnJkbzZrM3hGd2R4b0RSRTI1VTFsaS9cbnN6ZVVPaUtjYUlPNm9uTHVKU0VwNWFoOEpVbXFHYmJwUUdNOU0yLzBVRzZEcU52T0c1WWtLbFVrc0ovbTRsS0VcblQxZjBtK0J1RzlUU0VzY1MvYWorYnZFck02VVJML2JPajFnejRYcFh5ZVVqYXdVUUxqZGgxUUJoSk9Gdi9rZmFcblkyL2N5U25JZGRhSktuQVJMUWpjRzV1N0tqcTE2T1JYblJHRnEzL0xRa3pjdktxeXZYWDhmWjNpOTZrTVVzTnlcbmNIallMLytHZnM1NnV5RTlaK0Y2TG9EVHE4WFFFZGZSR0tSb3AzMlc0a3lZR1ZGYlJvUFdKamFMazloL0lhejhcbkVyTUoxcVhSL1liMXpuU2MxVUJTK1pmR1ZyazlSY2NFR0MyZFVySzVZdndKTEdDQ2RrTVg0SmNENzZDTDI2NE5cbk0wU3o5RTllVHI3aFxuLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLVxuIn1dLCJOb2RlcyI6W3siSWQiOiJFcFlHQ3JUREoxcFd1NXdHOFZheHBJRVFtY1hTWlp4anFoMW1tRnYwKzhjQyIsIkFkZHJlc3MiOiIzNy4xNDguMTk2LjkzOjExNDIwIiwiVGxzX2NlcnRpZmljYXRlIjoiLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tXG5NSUlGelRDQ0E3V2dBd0lCQWdJVVA5bERJVXpPVnZZVDRrRjdzQ3Vha0d2bUdCOHdEUVlKS29aSWh2Y05BUUVMXG5CUUF3Z1pJeEN6QUpCZ05WQkFZVEFrNU1NUkl3RUFZRFZRUUlEQWxHYkdWMmIyeGhibVF4RVRBUEJnTlZCQWNNXG5DRXhsYkhsemRHRmtNUkl3RUFZRFZRUUtEQWw0ZUc1bGRIZHZjbXN4RWpBUUJnTlZCQXNNQ1hSbGMzUk9iMlJsXG5jekVUTUJFR0ExVUVBd3dLZUhndWJtVjBkMjl5YXpFZk1CMEdDU3FHU0liM0RRRUpBUllRWVdSdGFXNUFlSGd1XG5ibVYwZDI5eWF6QWVGdzB5TURBNU1qa3hNREEyTkRGYUZ3MHlNakE1TWpreE1EQTJOREZhTUlHU01Rc3dDUVlEXG5WUVFHRXdKT1RERVNNQkFHQTFVRUNBd0pSbXhsZG05c1lXNWtNUkV3RHdZRFZRUUhEQWhNWld4NWMzUmhaREVTXG5NQkFHQTFVRUNnd0plSGh1WlhSM2IzSnJNUkl3RUFZRFZRUUxEQWwwWlhOMFRtOWtaWE14RXpBUkJnTlZCQU1NXG5Dbmg0TG01bGRIZHZjbXN4SHpBZEJna3Foa2lHOXcwQkNRRVdFR0ZrYldsdVFIaDRMbTVsZEhkdmNtc3dnZ0lpXG5NQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUNEd0F3Z2dJS0FvSUNBUURBcEJOMVI0Yis1ZkE5MXd2alJSc21oUDJWXG4vSGQvYkdiRVBGdVh4NHBzbExrQytMRkJncDExL1h1N1padlk4SzNINnN0bGhBenZWSU1Qd1lyOC9xOG9jQjloXG5tdlNHdTBaNGo5UzNjNU5CT2NEQmgvWDZBM1BBSE9OZTNqUXJTOXlTMUh6SlhvcGwrd2RuSG9LT3hlajlJNEx6XG5PRkdwMkFCQ3lpSG5Wa2xCbmQ0MjRzMVJ3dHN2cEpURkVOcy9RMnBtMWVlSlM1RVJHZEFMcUM4WVZEWkZubTByXG5wbFQ4eWFYaktPM3ZMWUhnVjNQYmtFaXIxVXZGeEllT29BdUJyei9BQjlJampxOWprR3dEdGdaM1BJK2hrZVQ4XG4wK3VEb2JxZjU0T1lYa01MZFdrcEpKeEl1SzFOREtVd3lPQXROZ0wrWG01dXpSTTA4c0p1ejJJODFZc2JlQ2NLXG5ndGJJZkZRdGV0SDBad3pqb3ZIeitmaWQwOFNoVHBIRVRGeU5oZXhXVGcxODVJMzRQOWorSis4MzhNZmlyc0w4XG56NmxzcWlYWWdKYjRUTEFVMWkrenlxOU5jSWE1alRvMG9ObFlSTkVOeGo3d2VxMWRlaS9KaHN0VGRhMjVOaCtxXG5mZm5NdEo0M2haQXAxN2htVnFxbEdBZ2xEYkU2ZEpodlJsMTJNZ1N1M3kvZ2NsVERodUhzaEZMOEt2T2x1TmllXG5SZ2FCZ1UzNkVQb0dlNXJRd2lSUzZHZGRDc04xYm9YOFJ4Tng0ZTVwQWlzVWp6OHduZThLZXN1ZmpuZ0xPSitEXG5JSnJhVEZyYVQxbUlGbzdkcXNaUUMrUXNRNWsyN0dHcHZGWGNkWWdoUkx4eGI3NFJvdFZ0SVMwSTNnMHlBUEQ5XG5JM2RkZFM3UWl4RUNWRDl0R3dJREFRQUJveGt3RnpBVkJnTlZIUkVFRGpBTWdncDRlQzV1WlhSM2IzSnJNQTBHXG5DU3FHU0liM0RRRUJDd1VBQTRJQ0FRQjNpVjQ2N3JXeWNPZCtNZTJzTnh2WTlwbXZyOE9DS0Jlb0hmdDNGOGxiXG5jVW4wT1duSGNSTnFuUXRHcmpVdFZybHdRMFdQRjM4NmhXUkRWV2VIRkpwM3dsbU1HWUgvbFhLUFluczV6YU5DXG5rc0dVUUFXS3BxMkd1TmpFbkZjRXloMXFSZnAzaTltT0ZrVEFFQmZGRzRRTStUcCtobm9zUTFVTnRPMGwwK0Q3XG5aRmQ4Y3JDcmZHVE4vVEpXNXdUK1JLRWRPRWU0R3liZE5ZN05teDA2a1BqbEF4dDFEZW4zOUgrZGNCZGVRNlVFXG5VY2MwYU1SUlpwNnZ3VUlFWG94Z1FXa1FpODVjVXc1cGlGQ3BzN3p1LzJWV0h1TEpUYVd6MjI5MmJjWjNQWnJDXG5NY2plUUNJNnYzQ2dRNWFpSzJtTUVLaW1qc3ovTGNMajM3eGN1Y2l1K3JtQm5HQUhZOUFWS2FIZFl0VVQrZ0FxXG4zQ2xyODdnUk13QW9LTFBTSEkycmpCVVJDNGphR1ZpbitlNG5hYXZueDlDNEk3eFpIVEVJcGp1Vmx2emF5SEJqXG5Rdm9iSDVIWmNlWjI2WnVPVU50UXE5eHZveWFvQTJQczNWYjU1ZEhrUEtUK3VTeVhnVTZNTWpYQXVOcUw1S2RmXG5DcUNqZWNndzdDSGh3M3M4RXNHc3oyNTFwaU9JRi8rYlU0UXBMalNQZlRCL1h0emRMUXlmOVJKbnViUmpyQWhqXG5WSmhqcVoweEIrZ3NLTnBWM2U1eXBOYUxzRzRGeDFKaDNmOVBocURGSHA1bGpUZU14YXRQcnBmSTJQeDk4SGhwXG5pVFdvTFJHNjFBMGZYT2JKSXBtbGtTWnRQdFMxeWgxVDBSa1M1blZ6K2xHQzdSbVVZaDRlZWpIVk83NUNjZGZiXG5LUT09XG4tLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tXG4ifSx7IklkIjoiSHd4Sjh5UTJSMmxxbUE4dVkzM2wrRzlTRG9tKzNtSEY1V0hHTkZSYzB3d0MiLCJBZGRyZXNzIjoiOTAuNjQuMTU3LjExNjoxMTQyMCIsIlRsc19jZXJ0aWZpY2F0ZSI6Ii0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLVxuTUlJRnV6Q0NBNk9nQXdJQkFnSVVHTUpuNDVwckFQajcyRVVYNW11RHpNbjFwVE13RFFZSktvWklodmNOQVFFTFxuQlFBd2dZa3hDekFKQmdOVkJBWVRBbE5MTVFzd0NRWURWUVFJREFKVFN6RVRNQkVHQTFVRUJ3d0tRbkpoZEdselxuYkdGMllURVNNQkFHQTFVRUNnd0plSGh1WlhSM2IzSnJNUTR3REFZRFZRUUxEQVZ1YjJSbGN6RVRNQkVHQTFVRVxuQXd3S2VIZ3VibVYwZDI5eWF6RWZNQjBHQ1NxR1NJYjNEUUVKQVJZUVlXUnRhVzVBZUhndWJtVjBkMjl5YXpBZVxuRncweU1EQTJNamN4TmpBM016bGFGdzB5TWpBMk1qY3hOakEzTXpsYU1JR0pNUXN3Q1FZRFZRUUdFd0pUU3pFTFxuTUFrR0ExVUVDQXdDVTBzeEV6QVJCZ05WQkFjTUNrSnlZWFJwYzJ4aGRtRXhFakFRQmdOVkJBb01DWGg0Ym1WMFxuZDI5eWF6RU9NQXdHQTFVRUN3d0ZibTlrWlhNeEV6QVJCZ05WQkFNTUNuaDRMbTVsZEhkdmNtc3hIekFkQmdrcVxuaGtpRzl3MEJDUUVXRUdGa2JXbHVRSGg0TG01bGRIZHZjbXN3Z2dJaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQ1xuRHdBd2dnSUtBb0lDQVFEQnluQlkyRGdlMmNGbkhudkMzdHNRUHdIWnJkbllnZi85ZFVPMFBFRGxRUjNLR2twZVxuYXhwYTc4SEdLbnpZajNDYVhtVlFXYzd2NWVvNkRpbnovdEh4ZVBDNmlYaDZQSldPQ3FFUisxVFFRamhiRVJtY1xuUW02ZmtPck1oRm1HZ3lLQWEyTTZoZHJTQ2ZKU3l1WDk0Q1NnQXVlRTdKQzJLYTVFWUpWY2M5dkJKTVNwQmY0clxuY21saE9QL2tyS0RvbnJxYWM5Qjg3eDYrVEk3S1BwVVplSENZMWlZVFdPeHpLdjVHNjduekhrZ1dSbWtLQkZKalxueWdOckk4Tnk2V2twcmlNZ3d5MVlGa1BWa0NoeElOSEpnTThFYXR5VUlxWFRveGdKYjRMcy9WMkI1eGFOUXFHVVxuYUNQa3dNRUFQZ1gvMGgvNFdBRDlXcklXTGt3akNzcFFZcHBLVmFIckJ4VGx0d1pkbVZrMURkajhZVm5pWS9FNVxuZWtpSEluR3BqSlFqNS9TbG01UkljWVFxdjdvNWZPT0VSeEl5cWFmWk5tWkFYcy9KbHdXVjlHN01OSXVwOXhPZFxuTTdBOFcxTXk5bWlsbVFET1JWL1JuT3l3bkZGOFFGV0o2ZWtZV0RkNmJsYWhoUmdOd1ZUQi9lYlRFOGxFemprZlxuYzF1KzB1Y1NSSS9Qa1VFR29ST3Z1Q0lXK2JoYnBjbHE2K3lFL01yTGxEY0Y2dlQ5WmZVcXN5bE1Na3cyclV1VVxuMlFhcmpNMnQ2NzlUaVJVOWRURXJSWXFZWXJFS1d4WHZDbnBQOTVxckdtK1I2NmsxMy9zSHFOQTVjMWJuTnpmNlxuS0VUNjlFVUpmNGc3MUowVW0vV2xMdEg2aXNrNHNjZXBmbXJHTEpucUl5QXRrekY2RytlNmdwNUJTUUlEQVFBQlxub3hrd0Z6QVZCZ05WSFJFRURqQU1nZ3A0ZUM1dVpYUjNiM0pyTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFCV1xuU2dMSFFmd2c3NDBsRzNQdms4ZEtGbmJYeERkWUVMa3lKTDdza3BVRitPdFFkYnk4eXZxbHJ4blZ6bzlCdXpPYlxubFdiT3JkZ240c0tTRmJpcXAzWG5XYXF2a2YxajBSSE44bEQ1Y3kzenNrWHYwTm8ydmp6K1cxVVJIYlNFT09nOVxudGNHV0s1QWNGT3NRcFZjbDlTNHZQRWQ2TlAwa3J3ZysvOUJaNGpNUmNHbC9SZldTNWRWbkhVREhUYjlCampaS1xuamVobmRSY2o3VzNKTXN1T1JTRFVvVERVWG9oaUt5Y2dXVnkvR2QxemNxK1JGZlJValN0R1l3d2JkMXF4QmdMeVxuOFRTcXpQblAvVlhHSE0wYkhzR1VVa0xBQUhoSFdwc21PRkRycDFQT2Y4UE96LysyTEllWGdVWkJpR3JjMXFCVFxuSUI3c0xKNmVzalpvUm5aZlBSWVFqc09razU3bzZqVlA4ajBITDR6UTF3eXdwVzN3eWdtc1pGZ1VYSDhJWjlNL1xuaFVPTzJpam11NnBBYkRsbzR2bVA5S2d5YWdmY1hYY3R6OVBsOWNkeS9MVE02WlNHVUVvamMrNWJKOVVzc3NXMVxudU4wOVlYV3pwWTliVnBSRzJDb3kyUHJ5eFVBbGYwS2Q4cGVUNi9vQUJpd1p1Z1pGN09mWTYvSWpub1RQSVRha1xuZHJkM0Z0d3oxT0xneFgzQjFhU1BaVVF0RG82Q1dWcGpIQTRnRmlEZWRDQXdlM0RKUzN1L1psNzFWV1VsM2FSSlxuV0dlenIyM0g3UTRiMm9GeEo2RHZjb283OERlYmRKcUxnT2E2a2NuYy9CeWlhWkVJanpobVJnZEEwODJXYTIzNFxuano4Qmx5UmFrU2d0UUJtTU5MbEpWM0ZyaVdYZzJNVlRtOHNhSEZLTzRnPT1cbi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS1cbiJ9LHsiSWQiOiJ0UElEcWQ3bWJqeTQ2TmFZcjdQK1NpdlpkTjNpczk2eGZFUXBTQ2xSTUxvQyIsIkFkZHJlc3MiOiIxNzguMTI0LjIxMS4yMDE6MTE0MjAiLCJUbHNfY2VydGlmaWNhdGUiOiItLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS1cbk1JSUYxVENDQTcyZ0F3SUJBZ0lVYVJZeENpWDI4NHR3NXRxK3grOU80S2N1SkJVd0RRWUpLb1pJaHZjTkFRRUxcbkJRQXdnWll4Q3pBSkJnTlZCQVlUQWtKWk1SZ3dGZ1lEVlFRSURBOU5hVzV6YTJGNVlTQnZZbXhoYzNReERqQU1cbkJnTlZCQWNNQlUxcGJuTnJNUkl3RUFZRFZRUUtEQWw0ZUc1bGRIZHZjbXN4RWpBUUJnTlZCQXNNQ1hSbGMzUk9cbmIyUmxjekVUTUJFR0ExVUVBd3dLZUhndWJtVjBkMjl5YXpFZ01CNEdDU3FHU0liM0RRRUpBUllSTXpRMU56VXlcbk5rQm5iV0ZwYkM1amIyMHdIaGNOTWpBd056STNNVFV3TmpRNVdoY05Nakl3TnpJM01UVXdOalE1V2pDQmxqRUxcbk1Ba0dBMVVFQmhNQ1Fsa3hHREFXQmdOVkJBZ01EMDFwYm5OcllYbGhJRzlpYkdGemRERU9NQXdHQTFVRUJ3d0ZcblRXbHVjMnN4RWpBUUJnTlZCQW9NQ1hoNGJtVjBkMjl5YXpFU01CQUdBMVVFQ3d3SmRHVnpkRTV2WkdWek1STXdcbkVRWURWUVFEREFwNGVDNXVaWFIzYjNKck1TQXdIZ1lKS29aSWh2Y05BUWtCRmhFek5EVTNOVEkyUUdkdFlXbHNcbkxtTnZiVENDQWlJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dJUEFEQ0NBZ29DZ2dJQkFPdVhRNXFlUGdCYWJlckhcbjJJaEVuckxsUUNRd2U3UmphVnpNeXp1OHZwc1IzMTZBOHJOQzBSaEx4VlRCVmcyR2UyWi8vUUFCMzQ5dzhCeXVcbnFzSUczSjB4VGFzMHYveEJMVnB5L1RQU29pUGlPbTFOa1VFT1JLZEI5NjlOR3R6clJnQnR6OC9mRWQwWC9uS3pcbnhlb2ExZkZhTzhNcWR1eHpFMHhHbzBiTWlBOGo0SGpsZU5PVjVFYk14Rkovc0l4QjlpNkRnNXd5V2NtOVV2andcblBxWkppaGYyek9BV3RpQ2M0TlY3R0oyUmduaVVFblUvemVJbUxHNjZaRXdPbFowUVhCT3QvMWJxako1UnhUM0dcbkVpUXhXMWRqaVgzM2d5UVYxUlNuYklKRDFoRG1PNy9jSkptU2R3bEVmb3RMNzZXSmVMS0xpNXBQdlI3dnpwd0pcbmZmVTF2VC94aXgvem1UUjhFd0YyMXV2eWlHOUJWbmptUjVWTE9wSFdzOXVOR2h3RFZ4ZFF6ZnR5Y0pGNG9hMTNcbm5pMjRZSHk3KzF6eDExb0JJdVNDdVNoRm1pZ05kQzlCRFBIZitObHh3WDJFZUlWUkhJSUZOVXFRZFh4UHRDMVhcbnhqRG9NSUtOYmlnQ1ZRS1BpYmcvR3RRRDZHSUZ6VDJFeFVLc2tlVjFZdnNLSzc4TUdLZlBKdWtKbStzbGUyMXZcbmZITHMrYjNjNWVINlZrU1lNeTROaittdHg0bDcweGxSYi9UMkRITzdhSXpjTWkyOG1udmg0SENEWGFjQmJZTmtcbk9aYVZ1MUZONTNKMi9GVGpFTFhlb2NsY3FGcGRSb2J4Mk1yUS8zbTMrc3RLa3VCdVFpWXoyeXUyK29lOHZlRkxcbndHQ00rRHNFZGFtaHRNalhJRzM1UlVDc0VwL0hBZ01CQUFHakdUQVhNQlVHQTFVZEVRUU9NQXlDQ25oNExtNWxcbmRIZHZjbXN3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUhVUSthLytYckFIbFhvaCsxcE16Zmt4WVpIUW00NHdcbjhwb0xtMmtNTWJCYXFUM29TUkNmQkFkeW90NW5sM3NZV1JqZDZKL1VzWnBINnMwYVhhZWJXNzRjWnVyWTBHSnlcbkU0TTA3djZYYk9iQytIZExhZU1hQ1Y4MFRqTDBXSE1sczBmUzdFb3RPazVmbFgwMFYwcWxFMzZUcS84a2RFRHRcbjBWei9UNHJ0dGFTbXlHYXJzMFVyMmJFUkdXWUJxZ3psMllOVFRiTDNsNEpXREtVaUxwcVFDV2hseDRuVTd4dTdcbmJodnlYY1RvMzhGT2lMdnk1cG1sNURSOE9pd2l4Rzh0cUhlM0d1QjQzZ1d6OEZEcGJjWlZXT2IzUEVNN0htMHZcbmJoc1d1VG5oVUk1Sm9Way9JQ24weUtyZmFFa0hOUmFDUVl6dEtTd2Q5WjIveUFoYWh4cUZUd0xHT1hBRHY4MStcbmc3OWhsemF1Y1h0VURERGozK1h0bHhROExxa0FId0tkQ3NXUWJkZzZvWVRiQ1VLRGlXYW5jQzh2Y3l6UU96L3hcbkVUUHQvdm9MQkVWdFFwdldLcHBhU2ROa0didzJkRy84ZFMrUUlLWURRZTRzcitmczQ4cWZnblRBelFVS3cvc3dcbi9BbnY4RnJTdTRNME1mWDhqSU52WjZvN3JCbFhXK0ZEWlJ1WGh3dW1lTkovOWVlTFNodVNQOWk2N3RpWEwxYmZcbld3dUdmZVExZ1AxalRWUFpCTjhGeDQvdFhZblIwTXA1NVFCSUlWaHEzUkl4eGpvMFFHRHlsbSthTXM0czJ6bGhcbmIrWXBPV01XbGQ2TjdzRzE4ckNMbmZkYXVDZmVETyt4MVNZSVA3M2lPNXloVkFWdmlNTHdQQmYwa0JyVkszdjFcbjlPUGd4U2NlQTRDRVxuLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLVxuIn1dLCJSZWdpc3RyYXRpb24iOnsiQWRkcmVzcyI6InBlcm1pc3Npb25pbmcucHJvZC5jbWl4LnJpcDoxMTQyMCIsIlRsc19jZXJ0aWZpY2F0ZSI6Ii0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLVxuTUlJRnFUQ0NBNUdnQXdJQkFnSVVHNTlqWFJUdjJIQ3lCWldzb2hHdWY2RUlyWHN3RFFZSktvWklodmNOQVFFTFxuQlFBd2dZQXhDekFKQmdOVkJBWVRBa3RaTVJRd0VnWURWUVFIREF0SFpXOXlaMlVnVkc5M2JqRVRNQkVHQTFVRVxuQ2d3S2VIZ2dibVYwZDI5eWF6RVBNQTBHQTFVRUN3d0dSR1YyVDNCek1STXdFUVlEVlFRRERBcDRlQzV1WlhSM1xuYjNKck1TQXdIZ1lKS29aSWh2Y05BUWtCRmhGaFpHMXBibk5BZUhndWJtVjBkMjl5YXpBZUZ3MHlNREEyTVRneVxuTURVek1qbGFGdzB5TWpBMk1UZ3lNRFV6TWpsYU1JR0FNUXN3Q1FZRFZRUUdFd0pMV1RFVU1CSUdBMVVFQnd3TFxuUjJWdmNtZGxJRlJ2ZDI0eEV6QVJCZ05WQkFvTUNuaDRJRzVsZEhkdmNtc3hEekFOQmdOVkJBc01Ca1JsZGs5d1xuY3pFVE1CRUdBMVVFQXd3S2VIZ3VibVYwZDI5eWF6RWdNQjRHQ1NxR1NJYjNEUUVKQVJZUllXUnRhVzV6UUhoNFxuTG01bGRIZHZjbXN3Z2dJaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRREFONEFWTHp4TVxuQnpDTXo2VVJPYXo2cGZSL3BjdFFnaEFhdnRoVzNTM3FzSSs3cGhrYlFvK2gxbHVJMVdtU3pwcnBYcTRDb0d0MFxuQXkrd2Z5SCtXMWxmRDNGUU9ETzh1SmxyNUo3VDdhRDhqRU5yVHZISy95b0RDRitQcHZpSk1HYkdKTzBHVGgwS1xuUWNWVWJkeFJIZTRlNDhnTXh4UW1jNTUrZGxaWmVKOGFHRHkvaHFZSEtmUVhQaHdQTENXRXFBVEpGa1lYOHI1R1xucWpPYXBvbkpQNGFPendmbFJFbjZ3MjVTSDcwWHBIbXFXKzZzclNqNWhjc1kxbGE5QTdYaEFYM1dvQ2NDSmtkaFxuRWJTS1U4dWFOSWUrM3JqN1Z0akxmbjBLZEZlOS9VYjBKWTdEdFI4Q1V4cEQ5K0RDc2pHRGw2TnVrQWZWNEo0blxuMk1zcnBuUVcyWDgwMCs1R2UwcC90ekliTXYwbzZyVmVGaE4zMFI4WjYrVEVPY2pkSzhNcGFTYUpXR1ljYUVnWVxuT3F5WFRPK3hIYjZnNG1WaXZYbWFVbUhSQnRYbkRtZ3JZanRSVlVGcU9YeFM2bDA1SHlKK29TTUF2MXNXcE4wcFxuT1lPSnU4RzNYMmZKcElRMWJUSjVRU1VhQncrTG1LWURtK1IzbWJaaCtOSTNBN2YyV1lTc0ZJQTFWcnlqcmRnS1xuanQ1ZHhxZ3Bta3FPV3A1Ym1lbmQ5MUtoTXY1Q2U5MFQ3ZERkQkt6b0Y2b3NLS2pKb1NnTWhnYnRjZGFTRktkeVxuamNhNjEyK2xmcFhCU1lQRzAza2JrSFh3b1BKWEZYd2RrUWt0UDZSdlo5NDNpa0FCSXU0cDJyMWZxV05zdnlqSVxuTW5zR3RKSTRFUkQrM21VeHFrbHprWVFnNjdMVHZJcWZMUUlEQVFBQm94a3dGekFWQmdOVkhSRUVEakFNZ2dwNFxuZUM1dVpYUjNiM0pyTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFDUU41TmFtdEpqZmcvZUUzVWxnYmFOSkhTTFxuaTVkZVdWYlJHNXJ0QVRHZXAxeTFyd21JUHRmZTN1QmgzMVh6NUhvMGtMNVhRbGdtVkFDQmhyMDZQNXJPSUM3SlxuRndKeGZnSC9qc3k2THJLN3IxeFBTMW9mcHR2ajRsditTRVdjYmE5eWdqZEQ4SUVVQkdnZ2ZYRUlLT0ZBVjVNVlxuV0ZBNDg4SlZITXVXcEE5dVFOdHRWWERVWW01OTdtc3d5dEF6Z3ozc05QU09iSUFxMmdUa041dmdQbFoySXhSa1xuSURERk5nblVLSm1rT0VSTm1oMkxEVC9mZk1oMGMrcnZuWDk2NnMzcmovWk1IZHh5bHlLZzlEM2I5dTBuTld5SFxuRlNkREkzMGhjTWVPalJneVJjbDdndllvbUNjRkdXYXJPd0w3eGFvQ0NkZkZ4MWkzdVlmaWhYVU1rRzNlaWhmSVxucWkwMFgxMjh5YThMZU5Da1ZzbHRkb2dxdzBzQXY2VUNDR0U2cDlOdFFKVTdIQjZCYk1Yd2Z5d0hkT0EwREUzd1xuTEM4bGthUUNkdTJLbWtnbWtxc1o0N2ZEY2tLTlp2NFJYNTlJUTY3ZFZzNUsrcjZNVlRxNURBb2RCd1hYaGxhV1xuMmJkUVUyNFpkVnpobXZKZVBSL0xvWTBUWHB5QWd5dlBHdE5MWUNUVWVEN0tvbnZkNGFOdm9EVlNpb2RRckNMVFxuMzBWK3MrMU9pTWNOaEVZdGdSZnlFeDNaUFBhdDc0QnU4VEhzSkp4TDVja3NWUVN4OXlPM05wcmlXWWtRTnhJRlxuYkg0NWdHV2U5QncvV2FaNiszRkZkK2JLSWFwVmJtV0FvMEluUnFQbUIxNUViWi8veTZNbW5lSFIyRjM3MjByNVxuTFMzS0hmaFE2M0Vuc0MrUHJBPT1cbi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS1cbiJ9LCJOb3RpZmljYXRpb24iOnsiQWRkcmVzcyI6IjMuMTI3LjIxNS4xMzM6MTE0MjAiLCJUbHNfY2VydGlmaWNhdGUiOiItLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS1cbk1JSUZ0akNDQTU2Z0F3SUJBZ0lKQU9hSjdoN09KUTY2TUEwR0NTcUdTSWIzRFFFQkN3VUFNSUdNTVFzd0NRWURcblZRUUdFd0pWVXpFTE1Ba0dBMVVFQ0F3Q1EwRXhFakFRQmdOVkJBY01DVU5zWVhKbGJXOXVkREVRTUE0R0ExVUVcbkNnd0hSV3hwZUhocGNqRVVNQklHQTFVRUN3d0xSR1YyWld4dmNHMWxiblF4RXpBUkJnTlZCQU1NQ21Wc2FYaDRcbmFYSXVhVzh4SHpBZEJna3Foa2lHOXcwQkNRRVdFR0ZrYldsdVFHVnNhWGg0YVhJdWFXOHdIaGNOTWpBd01USTBcbk1qTTFPRFExV2hjTk1qRXdNVEl6TWpNMU9EUTFXakNCakRFTE1Ba0dBMVVFQmhNQ1ZWTXhDekFKQmdOVkJBZ01cbkFrTkJNUkl3RUFZRFZRUUhEQWxEYkdGeVpXMXZiblF4RURBT0JnTlZCQW9NQjBWc2FYaDRhWEl4RkRBU0JnTlZcbkJBc01DMFJsZG1Wc2IzQnRaVzUwTVJNd0VRWURWUVFEREFwbGJHbDRlR2x5TG1sdk1SOHdIUVlKS29aSWh2Y05cbkFRa0JGaEJoWkcxcGJrQmxiR2w0ZUdseUxtbHZNSUlDSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQWc4QU1JSUNcbkNnS0NBZ0VBeWFZa0FnNHp1UjExa3ZTR2tuS2puSFhzYWNOaUdJaUFMRFlUOUVXL2t0VEVEcTNkZG1kMm0wNDJcblNOMHdydzhwc1lDT1p5TDk3V2V3RUZNWnlKbDZnNlFWdVRla3Q0T3c1QVpXSDZ2ZUI2Zm0xYWorYUZPTnpyWlJcbkM3akxMMDRBSlFwQy96UkI5OEtEaTgwamhUQzFoZHU2NWMvaTQwenF6bFlVN3JFVGFVUENrTFpFNkNrckV1RHVcbkRFMWZJbDdwVkIxZ3BNNkxvN2VVclJ6N1piNWlIak9rSlFDQjFYNnZnN2NuL2QvamJPdVppNUU0OGZwbnoyRDRcblJnUGhGRlJ5V3AvOUc3TnJ3eTVTNG4wZExGbVhGQW5IYzVjQ1hza0ptUHFjZnBvUUhZQkpPWUdKWkhYY2VBak9cbmNTaG41MnhRUURSK2V2ckRsNmtQMGNIejV6WkFOZTgyOCtwNll3SUY1ZFFlblFLcHBPK1BrZlhLYnhrczRQVWtcbnFBKy9kVlErajNxRG1HSE1mTVphWUd5Qi9vdGlta2NZUFVxTlJINDYySkE1MFh3ck1QYUkzV1VwMEMzeEVLd2lcbkxYdHdUNE5RTUh0Sy9pdkU0bUpod1VqbEFqNEZJc050NzNLaUxsKytiSkJ4QkhqeE9IWWpGWnV2Z1hYWEFCcU5cbm5pTm84OGFjL0QyOFNEeFBzVDBSMGhsaDlwRFV5NWNoVGs2ckRGTnVkdGFhZ2I3UC9vWDVpQ0tWa0FaOHBQZ3JcbjI1WXUxVVcwM2t6WGZiYmV3QWFJYjIyeDFSVmdLdHRpa1FQWnNpTG5MZ1VBZVViNE8xdTVoV3lBWXJ4WDFqek9cbnREaFVKcHpLZGx5THBHQW1xWENITFgyclBnbWdYUk9PTEEwUDdESlJ2cXNvcGw5OU9GOENBd0VBQWFNWk1CY3dcbkZRWURWUjBSQkE0d0RJSUtaV3hwZUhocGNpNXBiekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBZ0VBU0M1NjRhaU9cbitoTDM5a2RnRjh1Q21mTW1DamhoREpjaVoxWmdkZUQyK2ozY2YwdS8ybVJZMTBBZUlBaXI0SnR1dG9HN2dHaEdcbjJDcXBDRmRoc2licEpoSDBvWFdwSzRjK2Nya1JUeVNFUG9tOGJnT3ZhQmV6bkRSVXdJcTJGbHcxaW1EQWxDRGpcbk9iTk9Jdks3KzcwbGg0UkRBemFJV09xbklxbkY1em1TS3hOUnc3ZUNuQXZDTDdnTkJ4YW9Rb3JBY1hTQUV4cmVcblo3MEdRb1BNK2JiOWRLTmdnMlQvVnRxWE1jZVkrWHM5UHFYeUVBaGZMR1lrWUZBL1V2WGo4a0k3SHQxWDFTZEtcblhkcG0xQytCWEFud2oyU2ZQbUdCeVVIakpTVjNqVm0yN2NnM09qMVBKdWFXVTFuN05MdW41SlBWMm9yRnpZVVJcbk9aWTVIZE9KaHkxU2dQbFZLTytYOHFCeEkza3M4Q0FtQnFnN21XK3VBQUVzNTZiMldKNU83NVF6dkNzWFUxZzdcbmJTMHdTU2xLSFhiV3dZemVDZk16eDNXTzBDOWY5L1N0Uk5KS3Q0RXhwRUVESTRpOThSVVRNQUMyVlJoUzR5SzNcbkhWU0Y0ZGtJOWo4dnl2dmt5THNKaUtGb2dCRVc2RWxrazE3M1NXZVJ1NTMrU2sza0l1WEV5WTdlZ3ZhdElBZkdcblR0c1NBSVFxUndLLzhiRWdzNmZQaDl6ZXM4UkdWbTFGcmRiL1h1eHVyODZMMmJpa0UrS0RKY1hYTm9BVnREU3BcbnpZVlROOGE0MDY1dkhPU2hJTUhtZ2lISFdFMVdBa0gvWWRadG9NWVVMTFN6TEpKeGN5SnZZSGRyaEFlMWlvQXlcbkhxMUh5VDVuZWdDRGN0N1htdGx4ZkZoWXFwcTMzWVlwQ2hJPVxuLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLVxuIn0sIlVkYiI6eyJJZCI6IkFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQU09IiwiUHVibGljX2tleV9QRU0iOiItLS0tLUJFR0lOIFBVQkxJQyBLRVktLS0tLVxuTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBS0g4dndDSyt2YXptNk9qU0ZjdEtBZ0hCdEZhcU15dlxuZzBsclBXcDZ3K2N0bm5QNU1FbmhRa3ZHaDZCOUtlNUFXRHFGNk9VZG5SZ3dxMnNpcmhnRW9GVUNBd0VBQVE9PVxuLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tIn0sIkUyZSI6eyJQcmltZSI6IkUyRUU5ODNEMDMxREMxREI2RjFBN0E2N0RGMEU5QThFNTU2MURCOEU4RDQ5NDEzMzk0QzA0OUI3QThBQ0NFREMyOTg3MDhGMTIxOTUxRDlDRjkyMEVDNUQxNDY3MjdBQTRBRTUzNUIwOTIyQzY4OEI1NUIzREQyQUVERjZDMDFDOTQ3NjREQUI5Mzc5MzVBQTgzQkUzNkU2Nzc2MDcxM0FCNDRBNjMzN0MyMEU3ODYxNTc1RTc0NUQzMUY4QjlFOUFEODQxMjExOEM2MkEzRTJFMjlERjQ2QjA4NjREMEM5NTFDMzk0QTVDQkJEQzZBREM3MThERDJBM0UwNDEwMjNEQkI1QUIyM0VCQjQ3NDJERTlDMTY4N0I1QjM0RkE0OEMzNTIxNjMyQzRBNTMwRThGRkIxQkM1MURBRERGNDUzQjBCMjcxN0MyQkM2NjY5RUQ3NkI0QkRENUM5RkY1NThFODhGMjZFNTc4NTMwMkJFREJDQTIzRUFDNUFDRTkyMDk2RUU4QTYwNjQyRkI2MUU4RjNEMjQ5OTBCOENCMTJFRTQ0OEVFRjc4RTE4NEM3MjQyREQxNjFDNzczOEYzMkJGMjlBODQxNjk4OTc4ODI1QjQxMTFCNEJDM0UxRTE5ODQ1NTA5NTk1ODMzM0Q3NzZEOEIyQkVFRUQzQTFBMUEyMjFBNkUzN0U2NjRBNjRCODM5ODFDNDZGRkREQzFBNDVFM0Q1MjExQUFGOEJGQkMwNzI3NjhDNEY1MEQ3RDc4MDNEMkQ0RjI3OERFODAxNEE0NzMyMzYzMUQ3RTA2NERFODFDMEM2QkZBNDNFRjBFNjk5ODg2MEYxMzkwQjVEM0ZFQUNBRjE2OTYwMTVDQjc5QzNGOUMyRDkzRDk2MTEyMENEMEU1RjEyQ0JCNjg3RUFCMDQ1MjQxRjk2Nzg5QzM4RTg5RDc5NjEzOEU2MzE5QkU2MkUzNUQ4N0IxMDQ4Q0EyOEJFMzg5QjU3NUU5OTREQ0E3NTU0NzE1ODRBMDlFQzcyMzc0MkRDMzU4NzM4NDdBRUY0OUY2NkU0Mzg3MyIsIlNtYWxsX3ByaW1lIjoiIiwiR2VuZXJhdG9yIjoiMiJ9LCJDbWl4Ijp7IlByaW1lIjoiRkZGRkZGRkZGRkZGRkZGRkM5MEZEQUEyMjE2OEMyMzRDNEM2NjI4QjgwREMxQ0QxMjkwMjRFMDg4QTY3Q0M3NDAyMEJCRUE2M0IxMzlCMjI1MTRBMDg3OThFMzQwNERERUY5NTE5QjNDRDNBNDMxQjMwMkIwQTZERjI1RjE0Mzc0RkUxMzU2RDZENTFDMjQ1RTQ4NUI1NzY2MjVFN0VDNkY0NEM0MkU5QTYzN0VENkIwQkZGNUNCNkY0MDZCN0VERUUzODZCRkI1QTg5OUZBNUFFOUYyNDExN0M0QjFGRTY0OTI4NjY1MUVDRTQ1QjNEQzIwMDdDQjhBMTYzQkYwNTk4REE0ODM2MUM1NUQzOUE2OTE2M0ZBOEZEMjRDRjVGODM2NTVEMjNEQ0EzQUQ5NjFDNjJGMzU2MjA4NTUyQkI5RUQ1MjkwNzcwOTY5NjZENjcwQzM1NEU0QUJDOTgwNEYxNzQ2QzA4Q0ExODIxN0MzMjkwNUU0NjJFMzZDRTNCRTM5RTc3MkMxODBFODYwMzlCMjc4M0EyRUMwN0EyOEZCNUM1NURGMDZGNEM1MkM5REUyQkNCRjY5NTU4MTcxODM5OTU0OTdDRUE5NTZBRTUxNUQyMjYxODk4RkEwNTEwMTU3MjhFNUE4QUFBQzQyREFEMzMxNzBEMDQ1MDdBMzNBODU1MjFBQkRGMUNCQTY0RUNGQjg1MDQ1OERCRUYwQThBRUE3MTU3NUQwNjBDN0RCMzk3MEY4NUE2RTFFNEM3QUJGNUFFOENEQjA5MzNENzFFOEM5NEUwNEEyNTYxOURDRUUzRDIyNjFBRDJFRTZCRjEyRkZBMDZEOThBMDg2NEQ4NzYwMjczM0VDODZBNjQ1MjFGMkIxODE3N0IyMDBDQkJFMTE3NTc3QTYxNUQ2Qzc3MDk4OEMwQkFEOTQ2RTIwOEUyNEZBMDc0RTVBQjMxNDNEQjVCRkNFMEZEMTA4RTRCODJEMTIwQTkyMTA4MDExQTcyM0MxMkE3ODdFNkQ3ODg3MTlBMTBCREJBNUIyNjk5QzMyNzE4NkFGNEUyM0MxQTk0NjgzNEI2MTUwQkRBMjU4M0U5Q0EyQUQ0NENFOERCQkJDMkRCMDRERThFRjkyRThFRkMxNDFGQkVDQUE2Mjg3QzU5NDc0RTZCQzA1RDk5QjI5NjRGQTA5MEMzQTIyMzNCQTE4NjUxNUJFN0VEMUY2MTI5NzBDRUUyRDdBRkI4MUJERDc2MjE3MDQ4MUNEMDA2OTEyN0Q1QjA1QUE5OTNCNEVBOTg4RDhGRERDMTg2RkZCN0RDOTBBNkMwOEY0REY0MzVDOTM0MDYzMTk5RkZGRkZGRkZGRkZGRkZGRiIsIlNtYWxsX3ByaW1lIjoiIiwiR2VuZXJhdG9yIjoiMiJ9fQ=="

	ndfData, err := base64.StdEncoding.DecodeString(ndfEncoded)
	if err != nil {
		panic(err)
	}
	def := &NetworkDefinition{}
	err = json.Unmarshal(ndfData, def)
	if err != nil {
		panic(err)
	}

	return def
}
