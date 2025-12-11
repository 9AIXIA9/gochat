//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

type ProfileIDGenerator interface {
	Generate() ProfileID
}
