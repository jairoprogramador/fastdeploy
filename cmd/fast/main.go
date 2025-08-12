package main

func main() {
	rootCmd := NewRootCmd()
	rootCmd.AddCommand(NewTestCmd())
	rootCmd.AddCommand(NewSupplyCmd())
	rootCmd.AddCommand(NewPackageCmd())
	rootCmd.AddCommand(NewDeployCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.Execute()
}
