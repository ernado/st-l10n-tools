package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/st-10n/martian/resource"
)

var genCmd = &cobra.Command{
	Use: "generate",
	Aliases: []string{
		"gen", "g",
	},
	Short: "Generate po files",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			outDir, inDir string
			postfixes     []string
			limit         []string
			err           error
			languages     Languages

			english Language
		)
		if inDir, err = f.GetString("input"); err != nil {
			return err
		}
		if outDir, err = f.GetString("output"); err != nil {
			return err
		}
		if len(outDir) == 0 {
			return errors.New("blank output dir")
		}
		if err = viper.UnmarshalKey("languages", &languages); err != nil {
			return err
		}
		if limit, err = f.GetStringSlice("limit"); err != nil {
			return err
		}
		for _, lang := range languages {
			if lang.Code == "EN" {
				english = lang
				break
			}
		}
		if english.Code == "" {
			return errors.New("no english language configured (code=EN)")
		}
		if err = filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
			if path != inDir && info.IsDir() {
				return filepath.SkipDir
			}
			base := filepath.Base(path)
			if strings.HasPrefix(base, "english") {
				postfixes = append(postfixes, strings.TrimPrefix(base, "english"))
			}
			return nil
		}); err != nil {
			return err
		}
		if len(postfixes) == 0 {
			return errors.New("no english files found in input folder")
		}
		fmt.Println("postfixes:", postfixes)
	Loop:
		for _, lang := range languages {
			if len(limit) > 0 {
				isSelected := false
				for _, limitLang := range limit {
					if strings.ToLower(limitLang) == strings.ToLower(lang.Name) {
						isSelected = true
					}
					if strings.ToLower(limitLang) == strings.ToLower(lang.Code) {
						isSelected = true
					}
				}
				if !isSelected {
					continue Loop
				}
			}
			fmt.Println("Language:", lang.Name)
			if lang.Prefix == "" {
				lang.Prefix = strings.ToLower(
					strings.Replace(lang.Name, " ", "_", -1),
				)
			}
			if lang.Locale == "" {
				lang.Locale = strings.ToLower(lang.Code)
			}
			fmt.Printf("  prefix: %s\n", lang.Prefix)
			fmt.Printf("  code: %s\n", lang.Code)
			fmt.Printf("  locale: %s\n", lang.Locale)
			var entries resource.Entries
			for _, p := range postfixes {
				original, err := readFile(filepath.Join(
					inDir, "english"+p,
				))
				if err != nil {
					return fmt.Errorf("failed to read english translation file: %v", err)
				}
				translated, err := readFile(filepath.Join(
					inDir, lang.Prefix+p,
				))
				if err != nil {
					if !(os.IsNotExist(err) && p != ".xml") {
						return fmt.Errorf("failed to find translated file for %s", lang.Code)
					}
				}
				gotEntries, err := resource.Gen(resource.GenOptions{
					Original:   original,
					Translated: translated,
					Simplified: viper.GetStringSlice("simplified"),
				})
				if err != nil {
					return err
				}
				entries = append(entries, gotEntries...)
			}
			fmt.Printf("  entries: %d\n", entries.TranslatedCount())
			outDirStat, err := os.Stat(outDir)
			if err != nil {
				return err
			}
			targetDir := filepath.Join(outDir, lang.Locale)
			if err = os.MkdirAll(targetDir, outDirStat.Mode()); err != nil {
				return err
			}
			for _, name := range entries.Files() {
				fileName := fmt.Sprintf("%s.po", name)
				outFile, createErr := os.Create(path.Join(targetDir, fileName))
				if createErr != nil {
					return createErr
				}
				if err = entries.WriteFile(name, outFile); err != nil {
					return err
				}
				if err = outFile.Close(); err != nil {
					return err
				}

				fileName = fmt.Sprintf("%s.pot", name)
				outFile, createErr = os.Create(path.Join(targetDir, fileName))
				if createErr != nil {
					return createErr
				}
				if err = entries.WriteTemplateFile(name, outFile); err != nil {
					return err
				}
				if err = outFile.Close(); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	{
		f := genCmd.Flags()
		f.StringP("output", "o", ".", "output directory")
		f.StringP("input", "i", ".", "input directory")
		f.StringSlice("limit", nil, "limit languages")
	}
	rootCmd.AddCommand(
		genCmd,
	)
}
