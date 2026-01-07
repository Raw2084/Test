package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

type Product struct {
	Name    string
	Version string
}

func readInstalledProgramsFromRegistry(rootKey registry.Key, path string) ([]Product, error) {
	var products []Product

	key, err := registry.OpenKey(rootKey, path, registry.READ)
	if err != nil {
		return products, err
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return products, err
	}

	for _, subkey := range subkeys {
		progKey, err := registry.OpenKey(rootKey, path+`\`+subkey, registry.READ)
		if err != nil {
			continue
		}

		name, _, err := progKey.GetStringValue("DisplayName")
		if err != nil || name == "" {
			progKey.Close()
			continue
		}

		version, _, _ := progKey.GetStringValue("DisplayVersion")
		products = append(products, Product{Name: name, Version: version})

		progKey.Close()
	}

	return products, nil
}

func main() {
	paths := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	var allProducts []Product
	for _, p := range paths {
		products, err := readInstalledProgramsFromRegistry(p.root, p.path)
		if err != nil {
			log.Printf("Fehler beim Auslesen von %s: %v\n", p.path, err)
			continue
		}
		allProducts = append(allProducts, products...)
	}

	fmt.Println("Installierte Programme und deren Versionen (Registry):")
	for _, product := range allProducts {
		fmt.Printf("Name: %s, Version: %s\n", product.Name, product.Version)
	}

	fmt.Println("\nDrücke Enter, um das Programm zu beenden...")
	fmt.Scanln()
}
