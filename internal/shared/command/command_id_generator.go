//go:generate mockgen -source=command_id_generator.go -destination=./mocks/mock_command_id_generator.go -package=mocks
package command

type IDGenerator interface {
	Generate() ID
}
