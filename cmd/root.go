package cmd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"text_analyze/internal/model"
	"text_analyze/internal/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pg storage.Storage

var rootCmd = &cobra.Command{
	Use:   "readxml",
	Short: "Reads XML and processes it",
	Run: func(cmd *cobra.Command, args []string) {
		sensesFilePath := viper.GetString("s")
		serviceWordFilePath := viper.GetString("sw")

		senseItemSet, err := readXML(sensesFilePath)
		if err != nil {
			fmt.Println("Ошибка", err.Error())
			os.Exit(1)
		}

		senseMap := make(map[string]string)

		for _, item := range senseItemSet.Items {
			senseMap[item.Name] = item.Lemma
		}

		serviceWordItemSet, err := readXML(serviceWordFilePath)
		if err != nil {
			fmt.Println("Ошибка", err.Error())
			os.Exit(1)
		}

		serviceWordMap := make(map[string]struct{})

		for _, item := range serviceWordItemSet.Items {
			serviceWordMap[item.Text] = struct{}{}
		}

		input := strings.TrimSpace(strings.ToLower(viper.GetString("input")))

		inputFields := strings.Fields(input)

		normalizedInput := make([]string, 0, len(inputFields))

		for _, field := range inputFields {
			if _, ok := serviceWordMap[field]; !ok {
				normalizedInput = append(normalizedInput, field)
			}
		}

		twoWordsInput := make([]string, 0)

		if len(normalizedInput) <= 1 {
			twoWordsInput = normalizedInput
		} else {
			for i := 1; i < len(normalizedInput); i++ {
				twoWordsInput = append(twoWordsInput, fmt.Sprintf("%s %s", normalizedInput[i-1], normalizedInput[i]))
			}
		}

		lemmanizedInput := make(map[string]struct{}, 0)
		for _, exp := range twoWordsInput {
			if lemma, ok := senseMap[strings.ToUpper(exp)]; ok {
				lemmanizedInput[lemma] = struct{}{}
			}
		}

		if len(lemmanizedInput) == 0 {
			for _, exp := range normalizedInput {
				if lemma, ok := senseMap[strings.ToUpper(exp)]; ok {
					lemmanizedInput[lemma] = struct{}{}
				}
			}
		}
		if len(lemmanizedInput) == 0 {
			fmt.Println("Ничего не найдено")
			os.Exit(1)
		}
		for lemma := range lemmanizedInput {
			output, err := pg.GetOutput(lemma)
			if err != nil {
				fmt.Println("Ошибка", err.Error())
			}
			if output != "" {
				fmt.Println(output)
			}
		}
	},
}

func readXML(filePath string) (*model.ItemSet, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %v", err)
	}

	var itemSet model.ItemSet
	if err := xml.Unmarshal(byteValue, &itemSet); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге файла: %v", err)
	}

	return &itemSet, nil
}

func Execute() {
	pg = storage.NewStorage()
	pg.SetUp()

	rootCmd.PersistentFlags().String("input", "", "Input string")

	rootCmd.PersistentFlags().String("s", "", "Path to the senses XML file")
	viper.BindPFlag("s", rootCmd.PersistentFlags().Lookup("s"))

	rootCmd.PersistentFlags().String("sw", "", "Path to the service words XML file")
	viper.BindPFlag("sw", rootCmd.PersistentFlags().Lookup("sw"))

	var input string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите запрос: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	input = strings.TrimSpace(input)

	viper.Set("input", input)

	cobra.CheckErr(rootCmd.Execute())
}
