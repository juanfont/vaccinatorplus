package main

import (
	"fmt"
	"log"
	"os"

	"github.com/juanfont/vaccinatorplus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaccinatorCmd = &cobra.Command{
	Use:   "vaccinatorplus",
	Short: "vaccinatorplus - a Telegram Bot to get notifications when you can get vaccinated in NL",
	Long: `
Juan Font Alonso <juanfontalonso@gmail.com> - 2021
https://github.com/juanfont/vaccinatorplus`,
}

var runCmd = &cobra.Command{
	Use:   "run YEAR",
	Short: "Launches the VaccinatorPlus",
	Run: func(cmd *cobra.Command, args []string) {
		year, err := cmd.Flags().GetInt("year")
		fmt.Println(year)
		token := viper.GetString("TELEGRAM_TOKEN")
		v, err := vaccinatorplus.NewVaccinator(token, "db.sqlite", year)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("💉 Launching VaccinatorPlus")
		v.Run()
	},
}

func main() {
	viper.AutomaticEnv()
	vaccinatorCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().IntP("year", "y", 1983, "Year")
	runCmd.MarkPersistentFlagRequired("year")

	if err := vaccinatorCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
