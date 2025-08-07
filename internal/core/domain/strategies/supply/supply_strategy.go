package supply

type SupplyStrategy interface {
	ExecuteSupply() error
}